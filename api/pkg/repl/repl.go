package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	database "github.com/idea456/painmon-api-go/internal/database"
	types "github.com/idea456/painmon-api-go/internal/types"
	utils "github.com/idea456/painmon-api-go/pkg/utils"
)

var TERMINATED string = "exit"

func isTerminate(line string) bool {
	return strings.EqualFold(strings.Trim(line, "\n"), TERMINATED)
}

func execute(db *database.Database, line string) {
	args := strings.Split(line, " ")

	if strings.ToLower(args[0]) == "get" {
		obj := db.Get(args[1])
		var objGet types.Talent
		err := json.Unmarshal(obj.([]byte), &objGet)
		if err != nil {
			log.Fatalf("Unable to unmarshall object: %+v\n", err)
		}
		log.Printf("Got object: %+v\n", objGet.Costs.Cost)
	}
}

func InitREPL(db *database.Database) {
	fmt.Println("painmon-cli 1.0.0")
	fmt.Printf(">>> ")

	reader := bufio.NewReader(os.Stdin)

	for {
		if line, _ := reader.ReadString('\n'); !isTerminate(line) {
			line = strings.Trim(line, "\n")
			utils.Validate(line, 1)
			execute(db, line)
			fmt.Printf(">>> ")
		} else {
			break
		}

	}
}
