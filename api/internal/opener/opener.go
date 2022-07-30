package opener

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	types "github.com/idea456/painmon-api-go/internal/types"
)

func OpenJSON[T types.Entry](path string) T {
	jsonFile, err := os.Open(path)

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
