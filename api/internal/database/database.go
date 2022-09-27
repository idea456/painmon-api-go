package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
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
	Config     *Config
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

func (db *Database) PreprocessDailies() {
	talents := GetCategory[types.Talent](utils.TALENTS)
	weapons := GetCategory[types.Weapon](utils.WEAPONS)
	characterMaterials := GetCategory[types.Material](utils.TALENT_MATERIAL)
	weaponMaterials := GetCategory[types.Material](utils.WEAPON_MATERIAL)

	daily := make(map[string][]types.Material, 0)

	for _, material := range characterMaterials {
		for characterKey, talent := range talents {
			if characterKey == "travelergeo" || characterKey == "traveleranemo" || characterKey == "travelerelectro" || characterKey == "travelerdendro" {
				characterKey = "lumine"
			}
			character := GetPath[types.Character](utils.CHARACTERS, fmt.Sprintf(".%s", characterKey))

			for _, cost := range talent.Costs.Cost {
				if cost.Name == material.FourStarName {
					material.Characters = append(material.Characters, character)
				}
			}

		}

		for _, day := range material.Day {
			daily[day] = append(daily[day], material)
		}
	}

	for materialKey, material := range weaponMaterials {
		for _, weapon := range weapons {
			key := strings.ToLower(strings.Replace(weapon.Material, " ", "", -1))
			rarity, _ := strconv.Atoi(weapon.Rarity)
			if len(weapon.Material) > 0 && materialKey == key && rarity >= 4 {
				material.Weapons = append(material.Weapons, weapon)
			}
		}

		for _, day := range material.Day {
			daily[day] = append(daily[day], material)
		}
	}

	Set("daily", daily)
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

	v, ok := any(entry).(types.Material)

	if ok {
		if kind == "talentmaterialtypes" {
			v.Type = utils.TALENT_MATERIAL_TYPE
		} else if kind == "weaponmaterialtypes" {
			v.Type = utils.WEAPON_MATERIAL_TYPE
		}
		obj[kind][strings.Split(id, ".")[0]] = v
	} else {
		obj[kind][strings.Split(id, ".")[0]] = entry
	}
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
				case "characters":
					setDataCategory[types.Character](obj, "characters", entryPath, fileName)
				case "domains":
					setDataCategory[types.Domain](obj, "domains", entryPath, fileName)
				case "talentmaterialtypes":
					setDataCategory[types.Material](obj, "talentmaterialtypes", entryPath, fileName)
				case "talents":
					setDataCategory[types.Talent](obj, "talents", entryPath, fileName)
				case "weaponmaterialtypes":
					setDataCategory[types.Material](obj, "weaponmaterialtypes", entryPath, fileName)
				case "weapons":
					setDataCategory[types.Weapon](obj, "weapons", entryPath, fileName)
				}
			}

			for key := range obj {
				err := Set(key, obj[key])
				if err != nil {
					continue
				}
			}

		}
	}
}

func Set(key string, obj interface{}) error {
	db := InitializeDatabase()
	res, err := db.Handler.JSONSet(key, ".", obj)
	if err != nil {
		log.Printf("Failed to SET object: %v\n", err)
		return err
	}
	log.Printf("[Database] Set %s: %+v\n", key, res)
	return nil
}

func SetPath(key string, obj interface{}, path string) error {
	db := InitializeDatabase()
	res, err := db.Handler.JSONSet(key, path, obj)
	if err != nil {
		log.Printf("Failed to set object %s at path %s : %v\n", key, path, err)
		return err
	}
	log.Printf("[Database] Set %s at path %s : %v\n", key, path, res)
	return nil
}

func GetCategory[T types.Entry | interface{}](category string) map[string]T {
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
func Get[T types.Entry | interface{}](key string) T {
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

func GetPath[T types.Entry](key string, path string) T {
	db := InitializeDatabase()
	res, err := db.Handler.JSONGet(key, path)
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
