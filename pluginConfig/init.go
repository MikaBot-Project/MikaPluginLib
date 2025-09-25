package pluginConfig

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var configPath string

func init() {
	configPath = os.Args[1]
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func ReadJson(fileName string, config any) {
	jsonFile, err := os.OpenFile(fmt.Sprintf("%s%s", configPath, fileName),
		os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	fileInfo, _ := jsonFile.Stat()
	var bytes []byte
	if fileInfo.Size() < 2 {
		bytes, err = json.Marshal(config)
		if err != nil {
			return
		}
		_, err = jsonFile.Write(bytes)
		_, err = jsonFile.Write([]byte("\n"))
		if err != nil {
			return
		}
	} else {
		bytes = make([]byte, fileInfo.Size())
		bytes, _ = io.ReadAll(jsonFile)
		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return
		}
	}
	if err != nil {
		log.Fatal("read config", err)
		return
	}
}

func SaveJson(fileName string, config any) {
	jsonFile, err := os.OpenFile(fmt.Sprintf("%s%s", configPath, fileName),
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	var bytes []byte
	bytes, err = json.Marshal(config)
	if err != nil {
		return
	}
	_, err = jsonFile.Write(bytes)
	_, err = jsonFile.Write([]byte("\n"))
	if err != nil {
		return
	}
}
