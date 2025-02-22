package models

type Transaction struct {
	MeterID    string  `json:"meter_id"`
	Timestamp  string  `json:"timestamp"`
	Location   string  `json:"location"`
	PowerUsage float64 `json:"power_usage"`
	Status     string  `json:"status"`
}
