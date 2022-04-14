package main

import (
	"flag"
	"fmt"

	"github.com/henrygeorgist/hydrographscalar/model"
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

	payload := "/workspaces/hydrographscaler/manifest/payload.yaml"

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

	for _, m := range payloadInstructions.DischargeModels {
		hsm, err := model.NewHydrographScalerModelFromS3(m.Model.Input, fs)
		fmt.Println(m.Model.Input)
		if err != nil {
			fmt.Println("error:", err)
			return
		} else {
			fmt.Println("computing model")
			fmt.Println(hsm)
			hsm.Compute(&payloadInstructions)

		}
	}

}
