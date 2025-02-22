package generator

import (
	"fmt"
	"math/rand"
	"time"

	models "github.com/kanna-karuppasamy/smart-grid-monitoring-producer/internal/models"
)

var (
	locations = []string{
		"New York Downtown",
		"Brooklyn Industrial",
		"Queens Residential",
		"Manhattan Commercial",
		"Bronx North",
		"Staten Island South",
		"Long Island City",
		"Williamsburg",
		"Astoria",
		"Forest Hills",
	}

	statuses        = []string{"online", "offline"}
	powerUsageRange = [2]float64{0.5, 10.0}
)

func GenerateMeterID(location string, number int) string {
	return fmt.Sprintf("%s_METER_%03d", location[:3], number)
}

func GenerateRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GenerateTimestamp() string {
	now := time.Now()
	twoHoursAgo := now.Add(-2 * time.Hour)
	randomTime := twoHoursAgo.Add(time.Duration(rand.Int63n(int64(2 * time.Hour))))
	return randomTime.UTC().Format(time.RFC3339)
}

func GenerateTransaction() models.Transaction {
	location := locations[rand.Intn(len(locations))]
	meterNumber := rand.Intn(20) + 1

	return models.Transaction{
		MeterID:    GenerateMeterID(location, meterNumber),
		Timestamp:  GenerateTimestamp(),
		Location:   location,
		PowerUsage: float64(int(GenerateRandomFloat(powerUsageRange[0], powerUsageRange[1])*100)) / 100,
		Status:     statuses[rand.Intn(len(statuses))],
	}
}
