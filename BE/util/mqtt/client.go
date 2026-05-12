package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MessageCallback func(topic string, payload []byte)

type Client struct {
	client    paho.Client
	cfg       Config
	mu        sync.RWMutex
	connected bool
	onMessage MessageCallback
}

type Config struct {
	BrokerURL       string
	ClientID        string
	Username        string
	Password        string
	QOS             byte
	SubscribeTopics []string
	UseTLS          *bool
}

func NewClient(cfg Config, onMessage MessageCallback) (*Client, error) {
	if cfg.QOS > 2 {
		cfg.QOS = 1
	}

	broker, err := normalizeBrokerURL(cfg.BrokerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid broker url: %w", err)
	}
	cfg.BrokerURL = broker

	c := &Client{
		cfg:       cfg,
		onMessage: onMessage,
	}

	opts := paho.NewClientOptions().
		AddBroker(cfg.BrokerURL).
		SetClientID(cfg.ClientID).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetMaxReconnectInterval(5 * time.Minute).
		SetConnectTimeout(10 * time.Second).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetOrderMatters(false).
		SetCleanSession(true).
		SetOnConnectHandler(c.onConnect).
		SetConnectionLostHandler(c.onConnectionLost)
	isSecureScheme(cfg.BrokerURL)

	// Determine TLS requirement
	needsTLS := true

	if needsTLS {
		tlsConfig, err := buildTLSConfig()
		if err != nil {
			return nil, fmt.Errorf("build tls: %w", err)
		}
		opts.SetTLSConfig(tlsConfig)
	}

	// Wire paho logs to stdout so you see the real network error
	paho.ERROR = log.New(log.Writer(), "[MQTT ERR] ", log.LstdFlags)
	paho.CRITICAL = log.New(log.Writer(), "[MQTT CRT] ", log.LstdFlags)
	paho.WARN = log.New(log.Writer(), "[MQTT WRN] ", log.LstdFlags)

	c.client = paho.NewClient(opts)
	return c, nil
}

func normalizeBrokerURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("broker URL is empty")
	}

	if !strings.Contains(raw, "://") {
		raw = "tcp://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// Map user-friendly schemes to paho-compatible ones
	switch u.Scheme {
	case "mqtt":
		u.Scheme = "tcp"
	case "mqtts":
		u.Scheme = "ssl"
	}

	// Validate host
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("broker URL missing host: %s", raw)
	}

	// Validate / default port
	portStr := u.Port()
	if portStr == "" {
		port := 1883
		if isSecureScheme(u.String()) {
			port = 8883
		}
		u.Host = net.JoinHostPort(host, strconv.Itoa(port))
	} else {
		if _, err := strconv.Atoi(portStr); err != nil {
			return "", fmt.Errorf("invalid broker port: %s", portStr)
		}
	}

	return u.String(), nil
}

func isSecureScheme(urlStr string) bool {
	return strings.HasPrefix(urlStr, "ssl://") ||
		strings.HasPrefix(urlStr, "tls://") ||
		strings.HasPrefix(urlStr, "mqtts://") ||
		strings.HasPrefix(urlStr, "wss://")
}

func buildTLSConfig() (*tls.Config, error) {
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM([]byte(isrgRootX1PEM)); !ok {
		return nil, fmt.Errorf("failed to parse root CA")
	}
	return &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
	}, nil
}

func (c *Client) Connect() error {
	token := c.client.Connect()
	if !token.WaitTimeout(30 * time.Second) {
		return fmt.Errorf("mqtt connection timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt connect: %w", err)
	}
	return nil
}

func (c *Client) Disconnect(quiesceMs uint) {
	c.client.Disconnect(quiesceMs)
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// onConnect runs on paho's internal goroutine — must not block it.
func (c *Client) onConnect(client paho.Client) {
	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	go func() {
		for _, topic := range c.cfg.SubscribeTopics {
			t := client.Subscribe(topic, c.cfg.QOS, func(_ paho.Client, msg paho.Message) {
				if c.onMessage != nil {
					go c.onMessage(msg.Topic(), msg.Payload())
				}
			})
			if ok := t.WaitTimeout(10 * time.Second); !ok {
				log.Printf("subscribe timeout on %s", topic)
				continue
			}
			if err := t.Error(); err != nil {
				log.Printf("subscribe error %s: %v", topic, err)
			} else {
				log.Printf("subscribed: %s", topic)
			}
		}
	}()
}

func (c *Client) onConnectionLost(_ paho.Client, err error) {
	c.mu.Lock()
	c.connected = false
	c.mu.Unlock()
	log.Printf("mqtt connection lost: %v", err)
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload []byte) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("mqtt client not connected")
	}
	token := c.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) PublishJSON(topic string, qos byte, retained bool, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Publish(topic, qos, retained, b)
}

const isrgRootX1PEM = `-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIRAIIQz7DSQONZRGPgu2OCiwAwDQYJKoZIhvcNAQELBQAw
TzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh
cmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMTUwNjA0MTEwNDM4
WhcNMzUwNjA0MTEwNDM4WjBPMQswCQYDVQQGEwJVUzEpMCcGA1UEChMgSW50ZXJu
ZXQgU2VjdXJpdHkgUmVzZWFyY2ggR3JvdXAxFTATBgNVBAMTDElTUkcgUm9vdCBY
MTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAK3oJHP0FDfzm54rVygc
h77ct984kIxuPOZXoHj3dcKi/vVqbvYATyjb3miGbESTtrFj/RQSa78f0uoxmyF+
0TM8ukj13Xnfs7j/EvEhmkvBioZxaUpmZmyPfjxwv60pIgbz5MDmgK7iS4+3mX6U
A5/TR5d8mUgjU+g4rk8Kb4Mu0UlXjIB0ttov0DiNewNwIRt18jA8+o+u3dpjq+sW
T8KOEUt+zwvo/7V3LvSye0rgTBIlDHCNAymg4VMk7BPZ7hm/ELNKjD+Jo2FR3qyH
B5T0Y3HsLuJvW5iB4YlcNHlsdu87kGJ55tukmi8mxdAQ4Q7e2RCOFvu396j3x+UC
B5iPNgiV5+I3lg02dZ77DnKxHZu8A/lJBdiB3QW0KtZB6awBdpUKD9jf1b0SHzUv
KBds0pjBqAlkd25HN7rOrFleaJ1/ctaJxQZBKT5ZPt0m9STJEadao0xAH0ahmbWn
OlFuhjuefXKnEgV4We0+UXgVCwOPjdAvBbI+e0ocS3MFEvzG6uBQE3xDk3SzynTn
jh8BCNAw1FtxNrQHusEwMFxIt4I7mKZ9YIqioymCzLq9gwQbooMDQaHWBfEbwrbw
qHyGO0aoSCqI3Haadr8faqU9GY/rOPNk3sgrDQoo//fb4hVC1CLQJ13hef4Y53CI
rU7m2Ys6xt0nUW7/vGT1M0NPAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNV
HRMBAf8EBTADAQH/MB0GA1UdDgQWBBR5tFnme7bl5AFzgAiIyBpY9umbbjANBgkq
hkiG9w0BAQsFAAOCAgEAVR9YqbyyqFDQDLHYGmkgJykIrGF1XIpu+ILlaS/V9lZL
ubhzEFnTIZd+50xx+7LSYK05qAvqFyFWhfFQDlnrzuBZ6brJFe+GnY+EgPbk6ZGQ
3BebYhtF8GaV0nxvwuo77x/Py9auJ/GpsMiu/X1+mvoiBOv/2X/qkSsisRcOj/KK
NFtY2PwByVS5uCbMiogziUwthDyC3+6WVwW6LLv3xLfHTjuCvjHIInNzktHCgKQ5
ORAzI4JMPJ+GslWYHb4phowim57iaztXOoJwTdwJx4nLCgdNbOhdjsnvzqvHu7Ur
TkXWStAmzOVyyghqpZXjFaH3pO3JLF+l+/+sKAIuvtd7u+Nxe5AW0wdeRlN8NwdC
jNPElpzVmbUq4JUagEiuTDkHzsxHpFKVK7q4+63SM1N95R1NbdWhscdCb+ZAJzVc
oyi3B43njTOQ5yOf+1CceWxG1bQVs5ZufpsMljq4Ui0/1lvh+wjChP4kqKOJ2qxq
4RgqsahDYVvTH9w7jXbyLeiNdd8XM2w9U/t7y0Ff/9yi0GE44Za4rF2LN9d11TPA
mRGunUHBcnWEvgJBQl9nJEiU0Zsnvgc/ubhPgXRR4Xq37Z0j4r7g1SgEEzwxA57d
emyPxgcYxn/eR44/KJ4EBs+lVDR3veyJm+kXQ99b21/+jh5Xos1AnX5iItreGCc=
-----END CERTIFICATE-----`
