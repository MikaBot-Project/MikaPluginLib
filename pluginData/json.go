package pluginData

import (
	"encoding/json"
	"log"
)

func ReadJson[T any](name string, data *T) {
	file := openFile(name)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(data)
	if err != nil {
		log.Println("Error decoding json:", err)
	}
}

func SaveJson(name string, data any) {
	file := openFile(name)
	defer file.Close()
	encoder := json.NewEncoder(file)
	err := encoder.Encode(data)
	if err != nil {
		log.Println("Error encoding json:", err)
	}
}
