package main

import (
	"graphql/graph"
	"graphql/graph/generated"
	"graphql/internal/database"

	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	
	_ "github.com/mattn/go-sqlite3"
	"database/sql"

)

const defaultPort = "8080"

func main() {

	// Handling db con
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal("Failed to open database: %v", err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	courseDb := database.NewCourse(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		CategoryDB: categoryDb,
		CourseDB: courseDb,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
