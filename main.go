package main

import (
	"fmt"
	"os"
	"time"

	"github.com/henrygeorgist/hydrographscalar/model"
)

func main() {

	fs, err := model.Init()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// var configPath string
	// flag.StringVar(&configPath, "config", "", "please specify an input file using `-config=myconfig.json`")
	// flag.Parse()

	// if configPath == "" {
	// 	fmt.Println("given a blank path...")
	// 	fmt.Println("please specify an input file using `-config=myconfig.json`")
	// 	return
	// }
	//fmt.Printf("sleeping for 20 seconds, current unix time: %v\n", time.Now().Unix())

	//time.Sleep(20 * time.Second)

	payload := "/payload.yaml"
	payloadInstructions := model.Payload{}
	success := false
	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
	fs.Walk("", func(path string, file os.FileInfo) error {
		fmt.Println(path)
		if path == payload {
			payloadInstructions, err = model.LoadPayloadFromS3(path, fs)
			if err != nil {
				fmt.Println("error:", err)
				return err
			} else {
				success = true
			}
		}
		return nil
	})

	if !success {
		fmt.Println("not successful")
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
			fmt.Println(hsm)
			hsm.Compute(&payloadInstructions)

		}
	}
	fmt.Println("Made it to the end.....")
}
