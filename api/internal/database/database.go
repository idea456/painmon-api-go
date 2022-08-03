package database

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	redis "github.com/gomodule/redigo/redis"
	opener "github.com/idea456/painmon-api-go/internal/opener"
	"github.com/idea456/painmon-api-go/internal/types"
	rejson "github.com/nitishm/go-rejson/v4"
)

type Config struct {
	Port int
}

type Database struct {
	Handler    *rejson.Handler
	Connection redis.Conn
}

func InitializeDatabase() *Database {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Unable to establish connection to Redis: %s", err)
	}

	// defer func() {
	// 	_, err := conn.Do("FLUSHALL")
	// 	log.Printf("Closing Redis connection...\n")
	// 	err = conn.Close()
	// 	if err != nil {
	// 		log.Fatalf("Failed to communicate with Redis: %v", err)
	// 	}
	// }()

	log.Println("Successfully established connection to Redis!")

	rh := rejson.NewReJSONHandler()
	rh.SetRedigoClient(conn)

	return &Database{
		Handler: rh,
	}
}

var ctx = context.Background()

func (db *Database) InsertAll() {
	directories := []string{"artifacts", "talentmaterialtypes", "talents", "weaponmaterialtypes", "weapons"}

	for _, directory := range directories {
		files := opener.OpenDirectory(directory)

		for _, file := range files.Files {
			fileName := file.Name()
			var obj interface{}
			if !file.IsDir() {
				switch directory {
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
				db.Set(key, obj)
			}
		}
	}
}

func (db *Database) Set(key string, obj interface{}) {
	res, err := db.Handler.JSONSet(key, ".", obj)
	if err != nil {
		log.Fatalf("Failed to SET object: %v\n", err)
	}
	log.Printf("Set %s: %+v\n", key, res)

}

func (db *Database) Get(key string) interface{} {
	res, err := db.Handler.JSONGet(key, ".")
	if err != nil {
		log.Fatalf("Failed to GET object: %v", err)
	}
	log.Printf("Get %s: %+v\n", key, res)
	return res
}
