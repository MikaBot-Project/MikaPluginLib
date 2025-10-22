package pluginData

import (
	"log"
	"os"
	"path/filepath"
)

var dataPath string

func init() {
	dataPath = os.Args[2]
	err := os.MkdirAll(filepath.Dir(dataPath), 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func openFile(fileName string) *os.File {
	filePath := filepath.Join(dataPath, fileName)
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		log.Println("make dir ", filepath.Dir(filePath), " err:", err.Error())
		return nil
	}
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println("open file ", filePath, " err:", err.Error())
		return nil
	}
	return file
}

func readDir(path string) []os.DirEntry {
	dirPath := filepath.Join(dataPath, path)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Println("make dir ", dirPath, " err:", err.Error())
		return nil
	}
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println("open dir ", dirPath, " err:", err.Error())
		return nil
	}
	return dir
}
