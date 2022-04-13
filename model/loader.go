package model

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/USACE/filestore"
	"gopkg.in/yaml.v2"
)

func Init() (filestore.FileStore, error) {

	s3Conf := filestore.S3FSConfig{
		S3Id:     os.Getenv("AWS_ACCESS_KEY_ID"),
		S3Key:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Region: os.Getenv("AWS_DEFAULT_REGION"),
		S3Bucket: os.Getenv("S3_BUCKET")}

	fs, err := filestore.NewFileStore(s3Conf)

	if err != nil {
		log.Fatal(err)
	}

	return fs, nil
}

// LoadPayload
func LoadPayloadFromS3(payloadFile string, fs filestore.FileStore) (Payload, error) {
	var p Payload

	data, err := fs.GetObject(payloadFile)
	if err != nil {
		return p, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return p, err
	}

	err = yaml.Unmarshal(body, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}

// func LoadModelPayloadFromLocalJson(watPayload string) (wat.ModelPayload, error) {
// 	var wp wat.ModelPayload
// 	jsonFile, err := os.Open(watPayload)
// 	if err != nil {
// 		return wp, err
// 	}

// 	defer jsonFile.Close()

// 	jsonData, err := ioutil.ReadAll(jsonFile)
// 	if err != nil {
// 		return wp, err
// 	}

// 	errjson := json.Unmarshal(jsonData, &wp)
// 	if errjson != nil {
// 		return wp, errjson
// 	}
// 	return wp, nil

// }
