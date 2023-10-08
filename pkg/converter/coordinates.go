package converter

import (
	"fmt"
	"strconv"
	"strings"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func CoordinatesToString(coordinates []float64) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%f", coordinates[0]))
	builder.WriteString(",")
	builder.WriteString(fmt.Sprintf("%f", coordinates[1]))

	return builder.String()
}

func StringToCoordinates(str string) (Coordinates, error) {
	coordinates := strings.Split(str, ",")

	latitude, err := strconv.ParseFloat(strings.TrimSpace(coordinates[0]), 64)
	if err != nil {
		return Coordinates{}, err
	}

	longitude, err := strconv.ParseFloat(strings.TrimSpace(coordinates[1]), 64)
	if err != nil {
		return Coordinates{}, err
	}

	return Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}
