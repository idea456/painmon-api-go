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
	"github.com/idea456/painmon-api-go/pkg/utils"
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

func setDataCategory[T types.Entry](obj map[string]map[string]interface{}, kind string, entryPath string, id string) {
	entry := opener.OpenJSON[T](entryPath)

	if _, ok := obj[kind]; !ok {
		obj[kind] = make(map[string]interface{})
	}

	v, _ := any(entry).(types.Material)
	if kind == "talentmaterialtypes" {
		v.Type = utils.TALENT_MATERIAL_TYPE
	} else if kind == "weaponmaterialtypes" {
		v.Type = utils.WEAPON_MATERIAL_TYPE
	}
	obj[kind][strings.Split(id, ".")[0]] = v
}

func (db *Database) InsertAll() {
	for _, directory := range opener.GetDirectoryFiles("db/data/src/data/English") {
		files := opener.OpenDataDirectory(directory.Name())
		obj := make(map[string]map[string]interface{})

		for _, file := range files.Files {
			fileName := file.Name()
			if !file.IsDir() {
				entryPath := filepath.Join(files.Path, fileName)
				switch directory.Name() {
				case "artifacts":
					setDataCategory[types.Artifact](obj, "artifacts", entryPath, fileName)
				case "talentmaterialtypes":
					setDataCategory[types.Material](obj, "talentmaterialtypes", entryPath, fileName)
				case "talents":
					setDataCategory[types.Talent](obj, "talent", entryPath, fileName)
				case "weaponmaterialtypes":
					setDataCategory[types.Material](obj, "weaponmaterialtypes", entryPath, fileName)
				case "weapons":
					setDataCategory[types.Weapon](obj, "weapons", entryPath, fileName)
				case "domains":
					setDataCategory[types.Domain](obj, "domains", entryPath, fileName)
				}
			}

			for key := range obj {
				Set(key, obj[key])
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

func GetCategory[T types.Entry](category string) map[string]T {
	db := InitializeDatabase()
	res, err := db.Handler.JSONGet(category, ".")
	if err != nil {
		log.Fatalf("Failed to GET category %s: %v", category, err)
	}
	log.Printf("[Database] Get category %s\n", category)

	var buffer map[string]T
	json.Unmarshal(res.([]uint8), &buffer)

	return buffer
}

// golang does not support type parameters on methods as of now: https://github.com/golang/go/issues/49085
// therefore initializing methods such as Get, Query and GetCategory as functions instead of Database's methods
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

// TODO: Improve daily queries with Redisearch indexes
// func Query(day string) {
// 	db := InitializeDatabase()

// 	db.Connection.Do("FT.CREATE", weaponMaterialsIdx, )
// }
