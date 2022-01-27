package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalyze(t *testing.T) {
	spec := []struct {
		Name     string
		Expected Reading
		Input    map[ReadingType]map[DeviceType][]float64
	}{
		{
			Name:     "temp-1",
			Expected: map[DeviceType]string{"temp-1": "precise"},
			Input:    map[ReadingType]map[DeviceType][]float64{"thermometer": {"temp-1": {72.4, 76, 79.1, 75.6, 71.2, 71.4, 69.2, 65.2, 62.8, 61.4, 64, 67.5, 69.4}}},
		},
		{
			Name:     "temp-2",
			Expected: map[DeviceType]string{"temp-2": "ultra precise"},
			Input:    map[ReadingType]map[DeviceType][]float64{"thermometer": {"temp-2": {69.5, 70.1, 71.3, 71.5, 69.8}}},
		},
		{
			Name:     "hum-1",
			Expected: map[DeviceType]string{"hum-1": "keep"},
			Input:    map[ReadingType]map[DeviceType][]float64{"humidity": {"hum-1": {45.2, 45.3, 45.1}}},
		},
		{
			Name:     "hum-2",
			Expected: map[DeviceType]string{"hum-2": "discard"},
			Input:    map[ReadingType]map[DeviceType][]float64{"humidity": {"hum-2": {44.4, 43.9, 44.9, 43.8, 42.1}}},
		},
	}
	reference := Reference{Thermometer: 70.0, Humidity: 45.0}

	for _, s := range spec {
		res := analyze(reference, s.Input)
		assert.Equal(t, s.Expected, res, "")
	}
}
