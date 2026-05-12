package mqtt

import (
	entity "backend/entity/mqttpayload"
	"backend/util/mqtt"
)

type MQTTPublisher interface {
	SendCommand(cmd entity.ESP32Command) error
	SendConfig(cfg entity.ESP32Config) error
}

type mqttPublisher struct {
	client *mqtt.Client
}

func NewMQTTPublisher(client *mqtt.Client) MQTTPublisher {
	return &mqttPublisher{client: client}
}

func (p *mqttPublisher) SendCommand(cmd entity.ESP32Command) error {
	return p.client.PublishJSON("HyTerra/toESP32", 1, false, cmd)
}

func (p *mqttPublisher) SendConfig(cfg entity.ESP32Config) error {
	return p.client.PublishJSON("HyTerra/toESP32", 1, false, cfg)
}
