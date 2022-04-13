package main

import (
	"flag"
	"fmt"
	"model/model"
)

func main() {

	fs, err := model.Init()

	var configPath string
	flag.StringVar(&configPath, "config", "", "please specify an input file using `-config=myconfig.json`")
	flag.Parse()

	if configPath == "" {
		fmt.Println("given a blank path...")
		fmt.Println("please specify an input file using `-config=myconfig.json`")
		return
	}

	payload := "/data/payload.yaml"

	payloadInstructions, err := model.LoadPayloadFromS3(payload, fs)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// verify this is the right plugin
	if payloadInstructions.Plugin != "hydrograph_scaler" {
		fmt.Println("error", "expecting", "hydrograph_scaler", "got", payloadInstructions.Plugin)
		return
	}

	for _, location := range payloadInstructions.DischargeModels {
		hsml, err := model.NewHydrographScalerLocationFromS3(location.Model.Input, fs)
		if err != nil {
			fmt.Println("error:", err)
			return
		} else {
			fmt.Println(hsml)
			hsml.Compute(&payloadInstructions)

		}
	}

}
