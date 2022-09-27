package main

import (
	// "net/http"
	// "os"

	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/idea456/painmon-api-go/graph"
	"github.com/idea456/painmon-api-go/graph/generated"
	"github.com/idea456/painmon-api-go/internal/database"
	"github.com/idea456/painmon-api-go/pkg/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = utils.DEFAULT_PORT
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	db := database.InitializeDatabase()
	db.InsertAll()
	db.PreprocessDailies()
	// Set(rh, "ballad", obj)
	// res := Get(rh, "ballad")
	// var objGet Talent
	// err := json.Unmarshal(res.([]byte), &objGet)
	// if err != nil {
	// 	log.Fatalf("Unable to unmarshall object: %+v\n", err)
	// }
	// log.Printf("Got object: %+v\n", objGet.Costs["lvl10"])
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	// repl.InitREPL(db)

	// bw := browser.InitializeBrowser()

	defer func() {
		// bw.Close()
		db.Close()
	}()

}
