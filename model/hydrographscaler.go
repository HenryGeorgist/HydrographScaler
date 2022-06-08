package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/filestore"
	wm "github.com/usace/wat-api/model"
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

//model implementation
func (hsm HydrographScalerLocation) ModelName() string {
	return hsm.Name
}

func (hsm HydrographScalerLocation) Compute(eventSeed int64, realizationSeed int64, timewindow wm.TimeWindow, outputdestination string, fs filestore.FileStore) error {
	// bootstrap first (this is inefficient because it should only happen once per realization)
	bootStrap := hsm.Distribution.Bootstrap(realizationSeed)
	randomPeakValue := rand.New(rand.NewSource(eventSeed))
	value := bootStrap.InvCDF(randomPeakValue.Float64())
	//fmt.Println("value", value)
	//fmt.Println("bootStrap", bootStrap)
	currentTime := timewindow.StartTime
	timestepPercent := float64(1) / float64(len(hsm.Flows)-1)
	//create a writer
	output := strings.Builder{}
	fmt.Println("preparing to write output to:", outputdestination)
	//output.Write([]byte("Time,Flow\n"))
	for idx, flow := range hsm.Flows {
		if timewindow.EndTime.After(currentTime) {
			msg := fmt.Sprintf("%v,%v\n", float64(idx)*timestepPercent, flow*value)
			output.Write([]byte(msg))
			currentTime = currentTime.Add(hsm.TimeStep)
		} else {
			fmt.Println("encountered more flows than the time window.")
		}
	}
	fmt.Println(output.String())
	fso, err := fs.PutObject(outputdestination, []byte(output.String()))
	if err != nil {
		fmt.Println(fso)
		return err
	}
	return nil
}
func (hsm HydrographScalerModel) Compute(event *wm.EventConfiguration, fs filestore.FileStore, outputdest string) {
	//create random generator for realization and event
	erng := rand.NewSource(event.Event.Seed)
	rrng := rand.NewSource(event.Realization.Seed)
	for _, location := range hsm.Locations {
		err := location.Compute(erng.Int63(), rrng.Int63(), event.EventTimeWindow, outputdest, fs)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
	}
}
