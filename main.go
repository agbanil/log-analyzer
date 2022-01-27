package main

import (
	"encoding/json"
	"github.com/hpcloud/tail"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Reference struct {
	Thermometer float64
	Humidity    float64
}

type Constraint struct {
	UpperLimit float64
	LowerLimit float64
	Handle     func([]float64, float64, float64) bool
}

type Constraints struct {
	Values []Constraint
	Result string
}

type ChannelResponse struct {
	Reading    Reading
	File       string
	DeviceType DeviceType
}

type ReadingType string
type DeviceType string
type Reading map[DeviceType]string
type ReadingValues map[ReadingType]map[DeviceType][]float64

func main() {
	files := strings.Split(os.Getenv("FILE_PATHS"), ",")
	if len(files) < 1 {
		panic("You need to specify comma delimited file paths using the FILE_PATHS env variable")
	}
	m := make(map[string]Reading)
	c := make(chan ChannelResponse)
	for _, file := range files {
		go processFile(file, c)
	}

	for ch := range c {
		logger := log.New(os.Stdout, "", 0)
		if _, ok := m[ch.File]; !ok {
			m[ch.File] = Reading{}
		}
		m[ch.File][ch.DeviceType] = ch.Reading[ch.DeviceType]
		jsonString, _ := json.Marshal(m)
		logger.Println(string(jsonString))
	}
}

func processFile(file string, c chan ChannelResponse) {
	readingType := ""
	deviceType := ""
	reference := Reference{}
	m := ReadingValues{}

	t, err := tail.TailFile(file, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		text := strings.ToLower(line.Text)
		textSplit := strings.Split(text, " ")
		if strings.Contains(text, "reference") { // Get reference values
			m[ReadingType("thermometer")] = map[DeviceType][]float64{}
			m[ReadingType("humidity")] = map[DeviceType][]float64{}
			therm, err := strconv.ParseFloat(textSplit[1], 64)
			if err != nil {
				panic("Not a valid float")
			}
			hum, err := strconv.ParseFloat(textSplit[2], 64)
			if err != nil {
				panic("Not a valid float")
			}
			reference = Reference{Thermometer: therm, Humidity: hum}
		} else if strings.Contains(text, "thermometer") || strings.Contains(text, "humidity") {
			a := analyze(reference, m)
			if len(a) > 0 {
				c <- ChannelResponse{File: t.Filename, Reading: a, DeviceType: DeviceType(deviceType)}
				// Clear map after channel send
				delete(m[ReadingType(readingType)], DeviceType(deviceType))
			}
			readingType = textSplit[0]
			deviceType = textSplit[1]
		} else if len(textSplit) == 2 { // Get readings
			if _, ok := m[ReadingType(readingType)]; !ok {
				m[ReadingType(readingType)] = map[DeviceType][]float64{}
			}

			if _, ok := m[ReadingType(readingType)][DeviceType(deviceType)]; !ok {
				m[ReadingType(readingType)][DeviceType(deviceType)] = []float64{}
			}
			th, err := strconv.ParseFloat(textSplit[1], 64)
			if err != nil {
				panic("Not a valid float")
			}
			// Store reading by device type (hum-1, temp-1, etc)
			m[ReadingType(readingType)][DeviceType(deviceType)] = append(m[ReadingType(readingType)][DeviceType(deviceType)], th)
		}
	}
}

func getConstraints(reference Reference) map[ReadingType][]Constraints {
	therm := []Constraints{
		{
			Values: []Constraint{
				{
					UpperLimit: reference.Thermometer + 0.5,
					LowerLimit: reference.Thermometer,
					Handle:     meanCalc,
				},
				{
					UpperLimit: 3,
					Handle:     standardDeviation,
				},
			},
			Result: "ultra precise",
		},
		{
			Values: []Constraint{
				{
					UpperLimit: reference.Thermometer + 0.5,
					LowerLimit: reference.Thermometer,
					Handle:     meanCalc,
				},
				{
					UpperLimit: 5,
					LowerLimit: 3,
					Handle:     standardDeviation,
				},
			},
			Result: "very precise",
		},
	}

	hum := []Constraints{
		{
			Values: []Constraint{
				{
					UpperLimit: reference.Humidity + (reference.Humidity / 100),
					LowerLimit: reference.Humidity,
					Handle:     percentageHum,
				},
			},
			Result: "keep",
		},
	}

	return map[ReadingType][]Constraints{
		"thermometer": therm,
		"humidity":    hum,
	}
}

func analyze(reference Reference, m map[ReadingType]map[DeviceType][]float64) Reading {
	defaultResult := map[ReadingType]string{
		"thermometer": "precise",
		"humidity":    "discard",
	}

	constraints := getConstraints(reference)
	results := Reading{}
	for rd, v := range m {
		for dt, rds := range v {
			for _, cons := range constraints[rd] {
				successCounter := 0
				for _, val := range cons.Values {
					if val.Handle(rds, val.UpperLimit, val.LowerLimit) {
						successCounter++
					}
				}
				if successCounter == len(cons.Values) {
					results[dt] = cons.Result
				}
				if _, ok := results[dt]; !ok {
					results[dt] = defaultResult[rd]
				}
			}
		}
	}

	return results
}

func percentageHum(fl []float64, upperLimit float64, lowerLimit float64) bool {
	for _, v := range fl {
		if lowerLimit > v || v > upperLimit {
			return false
		}
	}
	return true
}

func standardDeviation(fl []float64, upperLimit float64, lowerLimit float64) bool {
	mn := mean(fl)
	sqrarr := []float64{}
	for _, v := range fl {
		mns := v - mn
		sqrarr = append(sqrarr, mns*mns)
	}

	scmn := mean(sqrarr)
	sqrt := math.Sqrt(scmn)
	return lowerLimit < sqrt && sqrt < upperLimit
}

func meanCalc(fl []float64, upperLimit float64, lowerLimit float64) bool {
	mn := mean(fl)
	return lowerLimit < mn && mn < upperLimit
}

func mean(fl []float64) float64 {
	var res float64
	for _, v := range fl {
		res += v
	}
	return res / float64(len(fl))
}
