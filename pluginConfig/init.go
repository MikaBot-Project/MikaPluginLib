package pluginConfig

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
)

var configPath string
var readJsonMap map[string]any
var readAllJsonMap map[string]any

func init() {
	readJsonMap = make(map[string]any)
	readAllJsonMap = make(map[string]any)
	configPath = os.Args[1]
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	pluginIO.OperatorMap["config"] = operatorHandler
}

func ReadJson(fileName string, config any) {
	readJsonMap[fileName] = config
	jsonFile, err := os.OpenFile(fmt.Sprintf("%s%s", configPath, fileName),
		os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)
	fileInfo, _ := jsonFile.Stat()
	var bytes []byte
	if fileInfo.Size() < 3 {
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

func ReadAllJson[T any](path string, config *map[string]T) {
	readAllJsonMap[path] = *config
	*config = make(map[string]T)
	DirPath := filepath.Join(configPath, path)
	err := os.MkdirAll(DirPath, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsonFolder, _ := os.ReadDir(DirPath)
	for _, jsonFile := range jsonFolder {
		if jsonFile.IsDir() {
			continue
		}
		nameArr := strings.Split(jsonFile.Name(), ".")
		if nameArr[len(nameArr)-1] != "json" {
			continue
		}
		var msgMap T
		ReadJson(filepath.Join(path, jsonFile.Name()), &msgMap)
		(*config)[strings.Join(nameArr[:len(nameArr)-1], ".")] = msgMap
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
