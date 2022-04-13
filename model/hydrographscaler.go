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
	Name     string        `json:"name"`
	Flows    []float64     `json:"flows"`
	TimeStep time.Duration `json:"timestep"`
	LP3      LP3Moments    `json:"lp3_moments"`
}

type LP3Moments struct {
	Mean              float64 `json:"mean"`
	StandardDeviation float64 `json:"standard_deviation"`
	Skew              float64 `json:"skew"`
	EYOR              int     `json:"equivalent_years_of_record"`
}

func NewHydrographScalerLocationFromS3(filepath string, fs filestore.FileStore) (HydrographScalerLocation, error) {
	// var hss HydrographScalerStruct
	hsm := HydrographScalerLocation{}

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

func (hsm HydrographScalerLocation) Compute(event *Payload) error {
	// bootstrap first (this is inefficient because it should only happen once per realization)
	lp3 := statistics.LogPearsonIIIDistribution{
		Mean:                    hsm.LP3.Mean,
		StandardDeviation:       hsm.LP3.StandardDeviation,
		Skew:                    hsm.LP3.Skew,
		EquivalentYearsOfRecord: hsm.LP3.EYOR,
	}

	bootStrap := lp3.Bootstrap(event.Config.Realization.Seed)
	randomPeakValue := rand.New(rand.NewSource(event.Config.Event.Seed))
	value := bootStrap.InvCDF(randomPeakValue.Float64())

	currentTime := event.Config.TimeWindow.StartTime

	for _, flow := range hsm.Flows {
		if event.Config.TimeWindow.EndTime.After(currentTime) {

			_ = fmt.Sprintf("%v,%v", currentTime, flow*value)
			// fmt.Println(msg)

			currentTime = currentTime.Add(hsm.TimeStep)
		} else {
			fmt.Println("encountered more flows than the time window.")
		}
	}
	return nil
}
