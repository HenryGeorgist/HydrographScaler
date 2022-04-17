package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/USACE/filestore"
	"github.com/usace/wat-api/wat"
	"gopkg.in/yaml.v2"
)

func Init() (filestore.FileStore, error) {
	mock := os.Getenv("S3_MOCK")
	disablessl := false
	s3fps := false
	mbool, err := strconv.ParseBool(mock)
	s3Conf := filestore.S3FSConfig{
		S3Id:     os.Getenv("AWS_ACCESS_KEY_ID"),
		S3Key:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Region: os.Getenv("AWS_DEFAULT_REGION"),
		S3Bucket: os.Getenv("S3_BUCKET"),
	}
	if err != nil {
		return nil, err
	}
	if mbool {
		dsslstring := os.Getenv("S3_DISABLE_SSL")
		disablessl, err = strconv.ParseBool(dsslstring)
		if err != nil {
			return nil, err
		}
		s3fpsstring := os.Getenv("S3_FORCE_PATH_STYLE")
		s3fps, err = strconv.ParseBool(s3fpsstring)
		if err != nil {
			return nil, err
		}
		s3Conf.Mock = mbool
		s3Conf.S3DisableSSL = disablessl
		s3Conf.S3ForcePathStyle = s3fps
		s3Conf.S3Endpoint = os.Getenv("S3_ENDPOINT")
	}
	fmt.Println(s3Conf)

	fs, err := filestore.NewFileStore(s3Conf)

	if err != nil {
		log.Fatal(err)
	}

	return fs, nil
}

// LoadPayload
func LoadPayloadFromS3(payloadFile string, fs filestore.FileStore) (wat.ModelPayload, error) {
	var p wat.ModelPayload
	fmt.Println("reading payload:", payloadFile)
	data, err := fs.GetObject(payloadFile)
	if err != nil {
		return p, err
	}

	body, err := ioutil.ReadAll(data)
	if err != nil {
		return p, err
	}
	//fmt.Println(string(body))
	err = yaml.Unmarshal(body, &p)
	if err != nil {
		return p, err
	}
	//fmt.Println(p)
	return p, nil
}

// LoadPayload
func UpLoadToS3(newS3Path string, fileBytes []byte, fs filestore.FileStore) (filestore.FileOperationOutput, error) {
	var repsonse *filestore.FileOperationOutput
	var err error
	repsonse, err = fs.PutObject(newS3Path, fileBytes)
	if err != nil {
		return *repsonse, err
	}

	return *repsonse, err
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
