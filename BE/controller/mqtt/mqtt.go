package mqtt

import (
	service "backend/service/mqttlisten"
)

type Controller struct {
	mqttService service.Servicer
}

func NewController(mqttService service.Servicer) *Controller {
	return &Controller{mqttService: mqttService}
}

func (c *Controller) OnMessage(topic string, payload []byte) {
	switch topic {
	case "HyTerra/toBE":
		c.mqttService.ProcessMQTTMessage(topic, payload)
	}
}
