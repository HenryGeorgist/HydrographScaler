package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/filestore"
)

type HydrographScalerModel struct {
	Locations []HydrographScalerLocation `json:"locations"`
}

type HydrographScalerLocation struct {
	Name         string                               `json:"name"`
	Flows        []float64                            `json:"flows"`
	TimeStep     time.Duration                        `json:"timestep"`
	Distribution statistics.LogPearsonIIIDistribution `json:"distribution"`
	//EquivalentYearsOfRecord int                                  `json:"equivalent_years_of_record"`
}

func NewHydrographScalerModelFromS3(filepath string, fs filestore.FileStore) (HydrographScalerModel, error) {
	// var hss HydrographScalerStruct
	hsm := HydrographScalerModel{}

	data, err := fs.GetObject(filepath)
	if err != nil {
		return hsm, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("read:", string(body))
	errjson := json.Unmarshal(body, &hsm)
	if errjson != nil {
		fmt.Println("Yep!")
		return hsm, errjson
	}

	return hsm, nil

}

//model implementation
func (hsm HydrographScalerLocation) ModelName() string {
	return hsm.Name
}

func (hsm HydrographScalerLocation) Compute(eventSeed int64, realizationSeed int64, event *Payload) error {
	// bootstrap first (this is inefficient because it should only happen once per realization)
	bootStrap := hsm.Distribution.Bootstrap(realizationSeed)
	randomPeakValue := rand.New(rand.NewSource(eventSeed))
	value := bootStrap.InvCDF(randomPeakValue.Float64())

	currentTime := event.Config.TimeWindow.StartTime

	for _, flow := range hsm.Flows {
		if event.Config.TimeWindow.EndTime.After(currentTime) {

			msg := fmt.Sprintf("%v,%v", currentTime, flow*value)
			fmt.Println(msg)

			currentTime = currentTime.Add(hsm.TimeStep)
		} else {
			fmt.Println("encountered more flows than the time window.")
		}
	}
	return nil
}
func (hsm HydrographScalerModel) Compute(event *Payload) {
	//create random generator for realization and event
	erng := rand.NewSource(event.Config.Event.Seed)
	rrng := rand.NewSource(event.Config.Realization.Seed)
	for _, location := range hsm.Locations {
		err := location.Compute(erng.Int63(), rrng.Int63(), event)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
	}
}
