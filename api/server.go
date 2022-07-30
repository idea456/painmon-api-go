package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/idea456/painmon-api-go/graph"
	"github.com/idea456/painmon-api-go/graph/generated"
	database "github.com/idea456/painmon-api-go/internal/database"
	types "github.com/idea456/painmon-api-go/internal/types"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	// fmt.Println(os.Getwd())

	// obj := OpenJSON[Talent]("db/data/src/data/English/talents/albedo.json")
	// fmt.Printf("%v", obj.Costs.Cost[1]["name"])
	rh := database.InitializeDBConnection()
	database.InsertAll[types.Talent](rh, "talents")
	// Set(rh, "ballad", obj)
	// res := Get(rh, "ballad")
	// var objGet Talent
	// err := json.Unmarshal(res.([]byte), &objGet)
	// if err != nil {
	// 	log.Fatalf("Unable to unmarshall object: %+v\n", err)
	// }
	// log.Printf("Got object: %+v\n", objGet.Costs["lvl10"])

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}
