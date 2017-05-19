package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/ONSDigital/go-ns/log"
	"github.com/ONSdigital/dp-published-collection-api/publishedcollection"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func createDB(dbAccess string) *sql.DB {
	db, dbErr := sql.Open("postgres", dbAccess)
	if dbErr != nil {
		log.ErrorC("Failed to connect to database", dbErr, log.Data{})
		panic(dbErr)
	}
	return db
}

func getEnvironmentVariable(name string, defaultValue string) string {
	environmentValue := os.Getenv(name)
	if environmentValue != "" {
		return environmentValue
	}
	return defaultValue
}

func main() {
	dbAccess := getEnvironmentVariable("DB_ACCESS", "user=dp dbname=dp sslmode=disable")
	port := getEnvironmentVariable("PORT", "9090")
	log.Namespace = "dp-publish-report"
	log.Debug("Starting published collection API", log.Data{"port": port})
	db := createDB(dbAccess)
	defer db.Close()
	api, err := publishedcollection.NewAPI(db)
	if err != nil {
		log.ErrorC("Failed to setup publish collection API", err, log.Data{"db": db})
		panic(err)
	}
	defer api.Close()
	router := mux.NewRouter()
	router.HandleFunc("/publishedcollection", api.GetList)
	router.HandleFunc("/publishedcollection/{colllectionId}", api.GetCollection)
	router.HandleFunc("/health", api.Health)
	errServer := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Error(errServer, log.Data{})
		panic(err)
	}
}
