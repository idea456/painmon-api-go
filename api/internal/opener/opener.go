package opener

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	types "github.com/idea456/painmon-api-go/internal/types"
)

type Files struct {
	Files []fs.FileInfo
	Path  string
}

func ChangeHomeDirectory() string {
	home := os.Getenv("BASE_DIR")
	if home == "" {
		home = "/Users/idea456/Documents/painmon-api-go"
	}
	os.Chdir(home)
	return home
}

func OpenJSON[T types.Entry](path string) T {
	home := ChangeHomeDirectory()
	jsonFile, err := os.Open(filepath.Join(home, path))
	if err != nil {
		log.Fatal("Unable to open JSON file!")
	}
	defer jsonFile.Close()

	// convert the json to byte array
	bytes, _ := ioutil.ReadAll(jsonFile)

	var entry T
	json.Unmarshal(bytes, &entry)
	return entry
}

func OpenDirectory(directory string) *Files {
	ChangeHomeDirectory()
	directoryPath := filepath.Join("db/data/src/data/English", directory)
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatalf("Failed to open directory: %v\n", err)
	}

	return &Files{
		Files: files,
		Path:  directoryPath,
	}
}
