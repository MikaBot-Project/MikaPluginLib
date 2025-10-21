package pluginFile

import (
	"encoding/gob"
	"log"
	"os"
)

func ReadBinary[T any](name string, data *T) {
	file := openFile(name)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	decoder := gob.NewDecoder(file)
	err := decoder.Decode(data)
	if err != nil {
		log.Println("Binary Decode error:", err)
	}
}

func SaveBinary(name string, data any) {
	file := openFile(name)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(data)
	if err != nil {
		log.Println("binary.Write failed:", err)
	}
}
