package database

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	// "encoding/json"
	// "flag"
	// "fmt"
	rejson "github.com/nitishm/go-rejson/v4"

	"github.com/gomodule/redigo/redis"
	opener "github.com/idea456/painmon-api-go/internal/opener"
	types "github.com/idea456/painmon-api-go/internal/types"
)

type Config struct {
	Port int
}

var ctx = context.Background()

func InsertAll[T types.Entry](rh *rejson.Handler, directory string) {
	files := opener.OpenDirectory(directory)

	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() {
			obj := opener.OpenJSON[T](filepath.Join(directoryPath, fileName))
			key := strings.Split(fileName, ".")[0]
			Set(rh, key, obj)
		}
	}
}

func Set(rh *rejson.Handler, key string, obj interface{}) {
	res, err := rh.JSONSet(key, ".", obj)
	if err != nil {
		log.Fatalf("Failed to SET object: %v\n", err)
	}
	log.Printf("Set %s: %+v\n", key, res)

}

func Get(rh *rejson.Handler, key string) interface{} {
	res, err := rh.JSONGet(key, ".")
	if err != nil {
		log.Fatalf("Failed to GET object: %v", err)
	}
	log.Printf("Get %s: %+v\n", key, res)
	return res
}

func InitializeDBConnection() *rejson.Handler {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Unable to establish connection to Redis: %s", err)
	}

	// defer func() {
	// 	_, err := conn.Do("FLUSHALL");
	// 	err = conn.Close()
	// 	if err != nil {
	// 		log.Fatalf("Failed to communicate with Redis: %v", err)
	// 	}
	// }()

	log.Println("Successfully established connection to Redis!")

	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	return rh
}
