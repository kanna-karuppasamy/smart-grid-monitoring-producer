package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/config"
	"github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/models"
)

// Generator is responsible for generating mock energy consumption data
type Generator struct {
	config  config.GeneratorConfig
	meters  map[string]*meterData
	regions []config.Region
	rnd     *rand.Rand
}

// meterData holds the persistent data for a single meter
type meterData struct {
	MeterID         string
	Region          string
	Latitude        float64
	Longitude       float64
	BuildingType    models.BuildingType
	BaseConsumption float64
}

// NewGenerator creates a new data generator with the provided configuration
func NewGenerator(cfg config.GeneratorConfig) *Generator {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	g := &Generator{
		config:  cfg,
		meters:  make(map[string]*meterData, cfg.MeterCount),
		regions: cfg.Regions,
		rnd:     rng,
	}

	// Initialize meters
	g.initializeMeters()

	return g
}

// initializeMeters creates the initial set of meters
func (g *Generator) initializeMeters() {
	for i := 0; i < g.config.MeterCount; i++ {
		meterId := fmt.Sprintf("meter-%06d", i+1)

		// Select region based on distribution
		region := g.selectRegion()

		// Generate location within the region
		lat := g.rnd.Float64()*(region.MaxLat-region.MinLat) + region.MinLat
		long := g.rnd.Float64()*(region.MaxLong-region.MinLong) + region.MinLong

		// Determine building type with simple distribution
		buildingTypeDist := g.rnd.Float64()
		var buildingType models.BuildingType
		switch {
		case buildingTypeDist < 0.7:
			buildingType = models.BuildingResidential
		case buildingTypeDist < 0.9:
			buildingType = models.BuildingCommercial
		default:
			buildingType = models.BuildingIndustrial
		}

		// Simple base consumption varies by building type
		var baseConsumption float64
		switch buildingType {
		case models.BuildingResidential:
			baseConsumption = 0.5 + g.rnd.Float64()*1.5 // 0.5-2 kWh
		case models.BuildingCommercial:
			baseConsumption = 3 + g.rnd.Float64()*7 // 3-10 kWh
		case models.BuildingIndustrial:
			baseConsumption = 15 + g.rnd.Float64()*25 // 15-40 kWh
		}

		// Save meter data
		g.meters[meterId] = &meterData{
			MeterID:         meterId,
			Region:          region.Name,
			Latitude:        lat,
			Longitude:       long,
			BuildingType:    buildingType,
			BaseConsumption: baseConsumption,
		}
	}
}

// selectRegion chooses a region based on the configured distribution
func (g *Generator) selectRegion() config.Region {
	val := g.rnd.Float64()
	cumulativePerc := 0.0

	for _, region := range g.regions {
		cumulativePerc += region.MeterPerc
		if val <= cumulativePerc {
			return region
		}
	}

	// Fallback to last region
	return g.regions[len(g.regions)-1]
}

// GenerateTransaction creates a new mock transaction
func (g *Generator) GenerateTransaction() models.Transaction {
	// Pick a random meter
	meterIds := make([]string, 0, len(g.meters))
	for id := range g.meters {
		meterIds = append(meterIds, id)
	}

	meterId := meterIds[g.rnd.Intn(len(meterIds))]
	meter := g.meters[meterId]

	// Current timestamp
	now := time.Now()

	// Determine meter status (operational, fault, or offline)
	status := models.StatusOperational
	if g.rnd.Float64() < g.config.FaultProbability {
		status = models.StatusFault
	} else if g.rnd.Float64() < g.config.OfflineProbability {
		status = models.StatusOffline
	}

	// Calculate consumption with time-of-day variation
	consumption := meter.BaseConsumption

	// Simple time-of-day multiplier
	hour := now.Hour()
	if hour >= 8 && hour <= 20 {
		// Daytime - higher consumption
		consumption *= 1.0 + 0.5*g.rnd.Float64()
	} else {
		// Nighttime - lower consumption
		consumption *= 0.4 + 0.3*g.rnd.Float64()
	}

	// Apply status effects
	if status == models.StatusFault {
		// Fault can cause erratic readings
		if g.rnd.Float64() < 0.5 {
			consumption *= 2 + g.rnd.Float64() // Spike
		} else {
			consumption *= 0.3 // Drop
		}
	} else if status == models.StatusOffline {
		consumption = 0
	}

	// Create transaction
	tx := models.Transaction{
		ID:             uuid.New().String(),
		MeterID:        meter.MeterID,
		Timestamp:      now,
		ConsumptionKWh: consumption,
		Latitude:       meter.Latitude,
		Longitude:      meter.Longitude,
		Region:         meter.Region,
		Status:         status,
		BuildingType:   meter.BuildingType,
		PeakLoad:       hour >= 17 && hour <= 21, // Evening peak
	}

	return tx
}
