package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

type HydrographScalerModel struct {
	Name          string                                `json:"name"`
	Flows         []float64                             `json:"flows"`
	TimeStep      time.Duration                         `json:"timestep"`
	FlowFrequency statistics.BootstrappableDistribution `json:"flow_frequency"`
}
type HydrographScalerEvent struct {
	RealizationSeed   int64
	EventSeed         int64
	OutputDestination string
	StartTime         time.Time
	EndTime           time.Time
}

//model implementation
func (hsm HydrographScalerModel) ModelName() string {
	return hsm.Name
}

func (hsm HydrographScalerModel) Compute(event HydrographScalerEvent) error {
	//bootstrap first (this is inefficient because it should only happen once per realization)
	b := hsm.FlowFrequency.Bootstrap(event.RealizationSeed)
	//then sample event level peak value
	r := rand.New(rand.NewSource(event.EventSeed))
	value := b.InvCDF(r.Float64())
	outputdest := event.OutputDestination + hsm.ModelName() + ".csv"

	w, err := os.OpenFile(outputdest, os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		fmt.Println(err)
	}

	defer w.Close()

	currentTime := event.StartTime
	fmt.Fprintln(w, "Time,Flow")

	for _, flow := range hsm.Flows {
		if event.EndTime.After(currentTime) {
			fmt.Fprintln(w, fmt.Sprintf("%v,%v", currentTime, flow*value))

			currentTime = currentTime.Add(hsm.TimeStep)
		} else {
			fmt.Println("encountered more flows than the time window.")
		}
	}
	return nil
}
