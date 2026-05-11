#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

// ============= NET CONFIG =============
const char* WIFI_SSID = "PRODI INF";
const char* WIFI_PASS = "admin123";

const char* MQTT_BROKER = "mqtt.tenzly.codes";
const int MQTT_PORT = 13792;
const char* MQTT_USER = "Kitasan";
const char* MQTT_PASS = "Kitasan1234";
const char* CLIENT_ID = "ESP32_SoilMonitor_001";

// ============= SSL CERT =============
const char mqttCert[] PROGMEM = R"EOF(
-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----
)EOF";

// ============= TOPIC =============
const char* TOPIC_PUBLISH = "HyTerra/toBE";
const char* TOPIC_SUBSCRIBE = "HyTerra/toESP32";

// ============= SENSOR CONFIG =============
const int SENSOR_PIN = 34;          // GPIO34 (ADC1_CH6)
const int ADC_RESOLUTION = 4095;    // 12-bit

// Kalibrasi default — bisa di-update via MQTT
int adcDry = 3200;    // Udara (kering)
int adcWet = 1200;    // Air (basah)
unsigned long publishInterval = 10000;  // 10 detik default

// ============= GLOBAL =============
WiFiClientSecure net;
PubSubClient mqttClient(net);

unsigned long mqttReconnectTimer = 0;
const unsigned long MQTT_RECONNECT_INTERVAL = 5000;

unsigned long lastPublish = 0;

// ============= FORWARD DECLARATIONS =============
void connectWiFi();
bool reconnectMQTT();
void mqttCallback(char* topic, byte* payload, unsigned int length);
int readSoilMoisture(int& rawADC);
void publishSoilData();
void publishStatus(const char* message, bool success);

// ============= SETUP =============
void setup() {
  Serial.begin(115200);
  delay(1000);
  
  analogReadResolution(12);
  analogSetAttenuation(ADC_11db);
  
  net.setCACert(mqttCert);
  
  mqttClient.setServer(MQTT_BROKER, MQTT_PORT);
  mqttClient.setCallback(mqttCallback);
  
  connectWiFi();
}

// ============= WIFI =============
void connectWiFi() {
  Serial.print("Connecting to WiFi");
  WiFi.begin(WIFI_SSID, WIFI_PASS);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.printf("\n[WiFi] Connected. IP: %s\n", WiFi.localIP().toString().c_str());
}

// ============= MQTT FUNCTIONS =============
bool reconnectMQTT() {
  if (millis() - mqttReconnectTimer < MQTT_RECONNECT_INTERVAL) {
    return mqttClient.connected();
  }
  mqttReconnectTimer = millis();
  
  if (mqttClient.connected()) return true;
  
  Serial.print("[MQTT] Connecting...");
  
  if (mqttClient.connect(CLIENT_ID, MQTT_USER, MQTT_PASS)) {
    Serial.printf(" connected to %s:%d\n", MQTT_BROKER, MQTT_PORT);
    mqttClient.subscribe(TOPIC_SUBSCRIBE);
    Serial.printf("[MQTT] Subscribed to %s\n", TOPIC_SUBSCRIBE);
    publishStatus("ESP32 Soil Monitor started", true);
    return true;
  } else {
    Serial.printf(" failed, rc=%d. Retry in %ds...\n", 
                  mqttClient.state(), MQTT_RECONNECT_INTERVAL/1000);
    return false;
  }
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
  Serial.printf("[MQTT] Message received [%s]\n", topic);

  char msgBuffer[512];
  if (length >= sizeof(msgBuffer) - 1) {
    Serial.println("[ERROR] Payload too large");
    return;
  }

  memcpy(msgBuffer, payload, length);
  msgBuffer[length] = '\0';
  Serial.printf("[MQTT] Payload: %s\n", msgBuffer);

  StaticJsonDocument<512> doc;
  DeserializationError error = deserializeJson(doc, msgBuffer);
  if (error) {
    Serial.printf("[ERROR] JSON parse failed: %s\n", error.c_str());
    return;
  }

  if (!doc.containsKey("action")) {
    Serial.println("[ERROR] Missing 'action' field");
    return;
  }

  const char* action = doc["action"];

  if (strcmp(action, "calibrate_dry") == 0) {
    adcDry = analogRead(SENSOR_PIN);
    Serial.printf("[ACTION] Calibrate DRY: ADC=%d\n", adcDry);
    publishStatus("Calibrate dry updated", true);
  }
  else if (strcmp(action, "calibrate_wet") == 0) {
    adcWet = analogRead(SENSOR_PIN);
    Serial.printf("[ACTION] Calibrate WET: ADC=%d\n", adcWet);
    publishStatus("Calibrate wet updated", true);
  }
  else if (strcmp(action, "set_interval") == 0) {
    unsigned long newInterval = doc["value"] | 10000;
    publishInterval = newInterval;
    Serial.printf("[ACTION] Interval set to %lu ms\n", publishInterval);
    publishStatus("Interval updated", true);
  }
  else if (strcmp(action, "read_now") == 0) {
    Serial.println("[ACTION] Force read triggered");
    lastPublish = 0;
  }
  else {
    Serial.printf("[WARN] Unknown command: %s\n", action);
  }
}

// ============= SENSOR & PUBLISH =============
int readSoilMoisture(int& rawADC) {
  rawADC = analogRead(SENSOR_PIN);
  int percent = map(rawADC, adcDry, adcWet, 0, 100);
  return constrain(percent, 0, 100);
}

void publishSoilData() {
  if (!mqttClient.connected()) return;

  int rawADC;
  int moisture = readSoilMoisture(rawADC);

  StaticJsonDocument<512> doc;
  doc["device"] = CLIENT_ID;
  doc["moisture_percent"] = moisture;
  doc["raw_adc"] = rawADC;          // Tambahan: kirim raw juga
  doc["adc_dry"] = adcDry;
  doc["adc_wet"] = adcWet;
  doc["interval_ms"] = publishInterval;
  doc["uptime_sec"] = millis() / 1000;

  char jsonBuffer[512];
  size_t len = serializeJson(doc, jsonBuffer);

  Serial.printf("[PUB] %s\n", jsonBuffer);
  
  if (mqttClient.publish(TOPIC_PUBLISH, jsonBuffer)) {
    Serial.println("[PUB] Success");
  } else {
    Serial.println("[PUB] Failed!");
  }
}

void publishStatus(const char* message, bool success) {
  if (!mqttClient.connected()) return;

  StaticJsonDocument<256> doc;
  doc["device"] = CLIENT_ID;
  doc["type"] = "status";
  doc["message"] = message;
  doc["status"] = success ? "success" : "error";
  doc["uptime_sec"] = millis() / 1000;

  char jsonBuffer[256];
  serializeJson(doc, jsonBuffer);
  mqttClient.publish(TOPIC_PUBLISH, jsonBuffer);
}

// ============= LOOP =============
void loop() {
  if (WiFi.status() != WL_CONNECTED) {
    connectWiFi();
  }

  if (!mqttClient.connected()) {
    reconnectMQTT();
  } else {
    mqttClient.loop();
  }

  unsigned long now = millis();
  if (now - lastPublish >= publishInterval) {
    lastPublish = now;
    publishSoilData();
  }
}