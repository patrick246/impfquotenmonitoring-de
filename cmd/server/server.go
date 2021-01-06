package main

import (
	"github.com/patrick246/impfquotenmonitoring-de/pkg/persistence/mongodb"
	"github.com/patrick246/impfquotenmonitoring-de/pkg/server"
	"log"
	"os"
)

func main() {
	mongodbUri := os.Getenv("MONGODB_URI")
	if mongodbUri == "" {
		mongodbUri = "mongodb://localhost:27017/vaccinestats"
	}

	mongodbClient, err := mongodb.NewClient(mongodbUri)
	if err != nil {
		log.Fatalln(err)
	}

	storageService, err := mongodb.NewStorageService(mongodbClient)
	if err != nil {
		log.Fatalln(err)
	}

	apiServer, err := server.NewServer(":8080", storageService)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(apiServer.ListenAndServe())
}
