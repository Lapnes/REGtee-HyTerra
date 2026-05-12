package mqttpub

import (
	"strconv"
	"time"

	entity "backend/entity/mqttpayload"
	repository "backend/repository/mqtt"
)

type CommandService interface {
	UpdateClientID(deviceID, newID string) error
	SetPump(deviceID string, active bool) error
	SetInterval(deviceID string, ms int) error
}

type commandService struct {
	pub repository.MQTTPublisher
}

func NewCommandService(pub repository.MQTTPublisher) CommandService {
	return &commandService{pub: pub}
}

func (s *commandService) UpdateClientID(deviceID, newID string) error {
	return s.pub.SendCommand(entity.ESP32Command{
		Cmd:    "SET_CLIENT_ID",
		Value:  newID,
		Target: deviceID,
		SentAt: time.Now().Unix(),
	})
}

func (s *commandService) SetPump(deviceID string, active bool) error {
	cmd := "PUMP_OFF"
	if active {
		cmd = "PUMP_ON"
	}
	return s.pub.SendCommand(entity.ESP32Command{
		Cmd:    cmd,
		Value:  "",
		Target: deviceID,
		SentAt: time.Now().Unix(),
	})
}

func (s *commandService) SetInterval(deviceID string, ms int) error {
	return s.pub.SendCommand(entity.ESP32Command{
		Cmd:    "SET_INTERVAL",
		Value:  strconv.Itoa(ms),
		Target: deviceID,
		SentAt: time.Now().Unix(),
	})
}
