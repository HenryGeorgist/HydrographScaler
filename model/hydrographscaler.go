package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/filestore"
	"github.com/usace/wat-api/wat"
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
	fmt.Println("reading:", filepath)
	data, err := fs.GetObject(filepath)
	if err != nil {
		return hsm, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return hsm, err
	}

	// fmt.Println("read:", string(body))
	errjson := json.Unmarshal(body, &hsm)
	if errjson != nil {
		fmt.Println("error:", errjson)
		return hsm, errjson
	}

	return hsm, nil

}

//model implementation
func (hsm HydrographScalerLocation) ModelName() string {
	return hsm.Name
}

func (hsm HydrographScalerLocation) Compute(eventSeed int64, realizationSeed int64, timewindow wat.TimeWindow, outputdestination string, fs filestore.FileStore) error {
	// bootstrap first (this is inefficient because it should only happen once per realization)
	bootStrap := hsm.Distribution.Bootstrap(realizationSeed)
	randomPeakValue := rand.New(rand.NewSource(eventSeed))
	value := bootStrap.InvCDF(randomPeakValue.Float64())
	//fmt.Println("value", value)
	//fmt.Println("bootStrap", bootStrap)
	currentTime := timewindow.StartTime
	//create a writer
	output := strings.Builder{}
	fmt.Println("preparing to write output to:", outputdestination)
	output.Write([]byte("Time,Flow"))
	for _, flow := range hsm.Flows {
		if timewindow.EndTime.After(currentTime) {
			msg := fmt.Sprintf("%v,%v", currentTime, flow*value)
			output.Write([]byte(msg))
			currentTime = currentTime.Add(hsm.TimeStep)
		} else {
			fmt.Println("encountered more flows than the time window.")
		}
	}
	fso, err := UpLoadToS3(outputdestination, []byte(output.String()), fs)
	if err != nil {
		fmt.Println(fso)
		return err
	}
	return nil
}
func (hsm HydrographScalerModel) Compute(event *wat.ModelPayload, fs filestore.FileStore) {
	//create random generator for realization and event
	erng := rand.NewSource(event.EventSeed)
	rrng := rand.NewSource(event.RealizationSeed)
	for idx, location := range hsm.Locations {
		path := fmt.Sprintf("%v/%v/%v/%v", event.OutputDestination, event.RealizationNumber, event.EventNumber, event.NecessaryOutputs[idx].Name)
		err := location.Compute(erng.Int63(), rrng.Int63(), event.EventTimeWindow, path, fs)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
	}
}
