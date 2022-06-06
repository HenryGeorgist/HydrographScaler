package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/henrygeorgist/hydrographscalar/model"
	wm "github.com/usace/wat-api/model"
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

	payloadInstructions, err := utils.LoadModelPayloadFromS3(payload, fs)
	if err != nil {
		fmt.Println("not successful", err)
		return
	}
	modelpath := ""
	eventConfigPath := ""
	for _, ldd := range payloadInstructions.LinkedInputs {
		switch ldd.Name {
		case "Event Configuration":
			//event configuration
			eventConfigPath = ldd.Fragment
		default:
			//model file
			modelpath = ldd.Fragment
		}
	}
	//load event configuration into memory
	ec := wm.EventConfiguration{}
	err = utils.LoadJsonPluginModelFromS3(eventConfigPath, fs, &ec)
	//load the model data into memory.
	hsm := model.HydrographScalerModel{}
	err = utils.LoadJsonPluginModelFromS3(modelpath, fs, &hsm)

	if err != nil {
		fmt.Println("error:", err)
		return
	} else {
		fmt.Println("computing model")
		//fmt.Println(hsm)
		hsm.Compute(&ec, fs)

	}
	fmt.Println("Made it to the end.....")
}
