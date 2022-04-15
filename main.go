package main

import (
	"fmt"
	"time"

	"github.com/henrygeorgist/hydrographscalar/model"
)

func main() {

	fs, err := model.Init()

	// var configPath string
	// flag.StringVar(&configPath, "config", "", "please specify an input file using `-config=myconfig.json`")
	// flag.Parse()

	// if configPath == "" {
	// 	fmt.Println("given a blank path...")
	// 	fmt.Println("please specify an input file using `-config=myconfig.json`")
	// 	return
	// }
	fmt.Printf("sleeping for 20 seconds, current unix time: %v\n", time.Now().Unix())

	time.Sleep(20 * time.Second)

	fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
	payload := "/media/payload.yaml"

	/* PLACEHOLDER for rapid iteration during development
	Allows to push the desired payload, we are about to read:) */
	/*localPayloadFile := "/workspaces/manifest/payload.yaml"
	jsonFile, err := os.Open(localPayloadFile)
	if err != nil {
		fmt.Println("jsonFile error:", err)
		return
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("jsonData error:", err)
		return
	}

	quickFixUploadReponse, err := model.UpLoadToS3(payload, jsonData, fs)
	if err != nil {
		fmt.Println("quickFixUploadReponse error:", err)
		return
	}
	fmt.Println(quickFixUploadReponse)
	*/
	/* Resume regular program */
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
