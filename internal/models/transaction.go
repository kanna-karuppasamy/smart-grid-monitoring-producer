package models

import (
	"time"
)

// MeterStatus represents the status of a smart meter
type MeterStatus string

const (
	StatusOperational MeterStatus = "operational"
	StatusFault       MeterStatus = "fault"
	StatusOffline     MeterStatus = "offline"
)

// BuildingType represents the type of building where the meter is installed
type BuildingType string

const (
	BuildingResidential BuildingType = "residential"
	BuildingCommercial  BuildingType = "commercial"
	BuildingIndustrial  BuildingType = "industrial"
)

// Transaction represents an energy consumption transaction from a smart meter
type Transaction struct {
	ID             string       `json:"id"`             // Unique transaction ID
	MeterID        string       `json:"meterId"`        // ID of the smart meter
	Timestamp      time.Time    `json:"timestamp"`      // When the reading was taken
	ConsumptionKWh float64      `json:"consumptionKWh"` // Energy consumption in kilowatt-hours
	Latitude       float64      `json:"latitude"`       // Geographic latitude
	Longitude      float64      `json:"longitude"`      // Geographic longitude
	Region         string       `json:"region"`         // Region name (e.g., Urban, Suburban, Rural)
	Status         MeterStatus  `json:"status"`         // Current meter status
	BuildingType   BuildingType `json:"buildingType"`   // Type of building
	PeakLoad       bool         `json:"peakLoad"`       // Whether this reading is during peak load time
}

// MeterInfo contains information about a smart meter
// This can be used for listing meters with faults or generating metrics
type MeterInfo struct {
	ID           string       `json:"id"`
	Status       MeterStatus  `json:"status"`
	Region       string       `json:"region"`
	BuildingType BuildingType `json:"buildingType"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
}
