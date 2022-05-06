package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/henrygeorgist/hydrographscalar/model"
	"github.com/usace/wat-api/utils"
)

func main() {
	fmt.Println("hydrograph_scaler plugin intializing")
	var payload string
	flag.StringVar(&payload, "payload", "", "please specify an input file using `-payload=pathtopayload.yml`")
	flag.Parse()

	if payload == "" {
		fmt.Println("given a blank path...")
		fmt.Println("please specify an input file using `-payload=pathtopayload.yml`")
		return
	}
	fmt.Println("initializing filestore")
	loader, err := utils.InitLoader("")
	if err != nil {
		log.Fatal(err)
		return
	}
	fs, err := loader.InitStore()
	if err != nil {
		log.Fatal(err)
		return
	}
	cache, err := loader.InitRedis()
	if err != nil {
		log.Fatal(err)
		return
	}
	payloadInstructions, err := utils.LoadModelPayloadFromS3(payload, fs)
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
	hsm := model.HydrographScalerModel{}
	path := payloadInstructions.ModelConfigurationResources[0].Authority + payloadInstructions.ModelConfigurationResources[0].Fragment
	err = utils.LoadJsonPluginModelFromS3(path, fs, &hsm)

	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println("computing model")
		//fmt.Println(hsm)
		hsm.Compute(&payloadInstructions, fs)

	}
	//}
	key := payloadInstructions.PluginImageAndTag + "_" + payloadInstructions.Name + "_R" + fmt.Sprint(payloadInstructions.Realization.Index) + "_E" + fmt.Sprint(payloadInstructions.Event.Index)
	out := cache.Set(key, "complete", 0)
	fmt.Println(out)
	fmt.Println("Made it to the end.....")
}
