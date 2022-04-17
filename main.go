package main

import (
	"fmt"

	"github.com/henrygeorgist/hydrographscalar/model"
)

func main() {

	fs, err := model.Init()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	payload := "/data/hydrographscaler/modelPayload.yml"
	payloadInstructions, err := model.LoadPayloadFromS3(payload, fs)
	if err != nil {
		fmt.Println("not successful", err)
		return
	}
	// verify this is the right plugin
	if payloadInstructions.Plugin != "hydrograph_scaler" {
		fmt.Println("error", "expecting", "hydrograph_scaler", "got", payloadInstructions.Plugin)
		return
	}

	for _, m := range payloadInstructions.DischargeModels {
		if len(m.Model.ModelFiles) == 0 {
			fmt.Println("These aren't the droids you're looking for...")
			return
		}

		hsm, err := model.NewHydrographScalerModelFromS3(m.Model.ModelFiles[0], fs)

		if err != nil {
			fmt.Println("error:", err)
			return
		} else {
			fmt.Println("computing model")
			//fmt.Println(hsm)
			hsm.Compute(&payloadInstructions, fs)

		}
	}
	fmt.Println("Made it to the end.....")
}
