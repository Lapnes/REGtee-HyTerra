package entity

type MQTTSensorPayload struct {
	Device          string `json:"device"`
	MoisturePercent int    `json:"moisture_percent"`
	RawADC          int    `json:"raw_adc"`
	ADCDry          int    `json:"adc_dry"`
	ADCWet          int    `json:"adc_wet"`
	IntervalMs      int    `json:"interval_ms"`
	PumpActive      bool   `json:"pump_active"`
	Timestamp       string `json:"timestamp"`
}

type ESP32Command struct {
	Cmd    string `json:"cmd"`
	Value  string `json:"value"`
	Target string `json:"target"`
	SentAt int64  `json:"sent_at"`
}

type ESP32Config struct {
	Device     string `json:"device"`
	IntervalMs int    `json:"interval_ms"`
	ADCDry     int    `json:"adc_dry"`
	ADCWet     int    `json:"adc_wet"`
}
