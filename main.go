package main

import (
	"fmt"

	"github.com/henrygeorgist/hydrographscalar/model"
)

func main() {
	fmt.Println("hydrograph_scaler plugin intializing")
	fmt.Println("initializing filestore")
	fs, err := model.InitStore()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("initializing Redis")
	rc, err := model.InitRedis()
	if err != nil {
		fmt.Println(err)
		return
	}
	// we can call set with a `Key` and a `Value`.
	err = rc.Set("pluginName", "hydrograph_scaler", 0).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		fmt.Println(err)
	}
	val, err := rc.Get("pluginName").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	payload := "/data/hydrographscaler/watModelPayload.yml"
	payloadInstructions, err := model.LoadPayloadFromS3(payload, fs)
	if err != nil {
		fmt.Println("not successful", err)
		return
	}
	// verify this is the right plugin
	if payloadInstructions.TargetPlugin != "hydrograph_scaler" {
		fmt.Println("error", "expecting", "hydrograph_scaler", "got", payloadInstructions.TargetPlugin)
		return
	}
	//load the model data into memory.
	hsm, err := model.NewHydrographScalerModelFromS3(payloadInstructions.ModelConfigurationPaths[0], fs)

	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println("computing model")
		//fmt.Println(hsm)
		hsm.Compute(&payloadInstructions, fs)

	}
	//}
	fmt.Println("Made it to the end.....")
}
