package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/usace/wat-api/wat"
)

func LoadModelPayload(watPayload string) (wat.ModelPayload, error) {
	var ts wat.ModelPayload
	jsonFile, err := os.Open(watPayload)
	if err != nil {
		return ts, nil
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ts, err
	}

	json.Unmarshal(jsonData, &ts)
	return ts, nil

}
func LoadModel(modelresourcepath string) (HydrographScalerModel, error) {
	var ts HydrographScalerModel
	jsonFile, err := os.Open(modelresourcepath)
	if err != nil {
		return ts, nil
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ts, err
	}

	json.Unmarshal(jsonData, &ts)
	return ts, nil

}
func main() {
	fmt.Println("running hydrographscaler")

	var configPath string
	flag.StringVar(&configPath, "config", "", "please specify an input file using `-config=myconfig.json`")
	flag.Parse()

	if configPath == "" {
		fmt.Println("given a blank path...")
		fmt.Println("please specify an input file using `-config=myconfig.json`")
		return
	} else {
		if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
			fmt.Println("input file does not exist or is inaccessible")
			return
		}

	}
	// Load modelpayload data
	modelpayload, err := LoadModelPayload(configPath)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("recieved payload:", modelpayload)
	}

	// verify this plugin is the right plugin
	if modelpayload.TargetPlugin != "hydrographscaler" {
		fmt.Println("error", "expecting", "hydrographscaler", "got", modelpayload.TargetPlugin)
		return
	}
	hsm, err := LoadModel(modelpayload.ModelConfigurationPath)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	//load it from file.
	event := HydrographScalerEvent{
		RealizationSeed:   modelpayload.RealizationSeed,
		EventSeed:         modelpayload.EventSeed,
		OutputDestination: modelpayload.OutputDestination,
		StartTime:         modelpayload.EventTimeWindow.StartTime,
		EndTime:           modelpayload.EventTimeWindow.EndTime,
	}
	hsm.Compute(event)
}
