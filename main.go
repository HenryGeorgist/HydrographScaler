package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/henrygeorgist/hydrographscalar/model"
	"github.com/usace/wat-api/wat"
)

func LoadModelPayload(watPayload string) (wat.ModelPayload, error) {
	var ts wat.ModelPayload
	jsonFile, err := os.Open(watPayload)
	if err != nil {
		return ts, err
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ts, err
	}

	errjson := json.Unmarshal(jsonData, &ts)
	if errjson != nil {
		return ts, errjson
	}
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
		return
	} else {
		fmt.Println("recieved payload:", modelpayload)
	}

	// verify this plugin is the right plugin
	if modelpayload.TargetPlugin != "hydrographscaler" {
		fmt.Println("error", "expecting", "hydrographscaler", "got", modelpayload.TargetPlugin)
		return
	}
	//load model from file
	hsm, err := model.NewHydrographScalerModelFromFile(modelpayload.ModelConfigurationPath)
	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println(hsm)
	}
	event := model.HydrographScalerEvent{
		RealizationSeed:   modelpayload.RealizationSeed,
		EventSeed:         modelpayload.EventSeed,
		OutputDestination: modelpayload.OutputDestination,
		StartTime:         modelpayload.EventTimeWindow.StartTime,
		EndTime:           modelpayload.EventTimeWindow.EndTime,
	}
	hsm.Compute(event)
}
