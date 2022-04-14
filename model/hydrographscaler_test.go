package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestModelMarshal(t *testing.T) {
	file, err := os.Open("/workspaces/hydrographscaler/configs/hsm.json")
	if err != nil {
		t.Fail()
	}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fail()
	}
	hsm := HydrographScalerModel{}
	// fmt.Println("read:", string(body))
	errjson := json.Unmarshal(body, &hsm)
	if errjson != nil {
		fmt.Println("Yep!")
	}
}
