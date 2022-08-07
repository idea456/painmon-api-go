package database

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"
	"sync"

	redis "github.com/gomodule/redigo/redis"
	opener "github.com/idea456/painmon-api-go/internal/opener"
	"github.com/idea456/painmon-api-go/internal/types"
	rejson "github.com/nitishm/go-rejson/v4"
)

var lock = &sync.Mutex{}

type Config struct {
	Port int
}

type Database struct {
	Handler    *rejson.Handler
	Connection redis.Conn
}

var database *Database

func InitializeDatabase() *Database {
	if database == nil {
		lock.Lock()
		defer lock.Unlock()
		if database == nil {
			log.Println("Creating database connection...")
			conn, err := redis.Dial("tcp", ":6379")
			if err != nil {
				log.Fatalf("Unable to establish connection to Redis: %s", err)
			}
			log.Println("Successfully established connection to Redis!")
			rh := rejson.NewReJSONHandler()
			rh.SetRedigoClient(conn)
			database = &Database{
				Handler:    rh,
				Connection: conn,
			}
		}
	}

	return database
}

func (db *Database) Close() {
	_, err := db.Connection.Do("FLUSHALL")
	log.Printf("Closing Redis connection...\n")
	err = db.Connection.Close()
	if err != nil {
		log.Fatalf("Failed to communicate with Redis: %v", err)
	}
}

var ctx = context.Background()

func (db *Database) InsertAll() {
	for _, directory := range opener.GetDirectoryFiles("db/data/src/data/English") {
		files := opener.OpenDataDirectory(directory.Name())

		for _, file := range files.Files {
			fileName := file.Name()
			var obj interface{}
			if !file.IsDir() {
				switch directory.Name() {
				case "artifacts":
					obj = opener.OpenJSON[types.Artifact](filepath.Join(files.Path, fileName))
				case "talentmaterialtypes":
					obj = opener.OpenJSON[types.TalentMaterial](filepath.Join(files.Path, fileName))
				case "talents":
					obj = opener.OpenJSON[types.Talent](filepath.Join(files.Path, fileName))
				case "weaponmaterialtypes":
					obj = opener.OpenJSON[types.WeaponMaterial](filepath.Join(files.Path, fileName))
				case "weapons":
					obj = opener.OpenJSON[types.Weapon](filepath.Join(files.Path, fileName))
				}
				key := strings.Split(fileName, ".")[0]
				Set(key, obj)
			}
		}
	}
}

func Set(key string, obj interface{}) {
	db := InitializeDatabase()
	res, err := db.Handler.JSONSet(key, ".", obj)
	if err != nil {
		log.Fatalf("Failed to SET object: %v\n", err)
	}
	log.Printf("[Database] Set %s: %+v\n", key, res)

}

// golang does not support type parameters on methods as of now: https://github.com/golang/go/issues/49085
func Get[T types.Entry](key string) T {
	db := InitializeDatabase()
	res, err := db.Handler.JSONGet(key, ".")
	if err != nil {
		log.Fatalf("Failed to GET object: %v", err)
	}
	log.Printf("[Database] Get %s: OK\n", key)

	var obj T
	err = json.Unmarshal(res.([]byte), &obj)
	if err != nil {
		log.Printf("Error in unmarshalling object: %+v\n", err)
	}
	return obj
}
