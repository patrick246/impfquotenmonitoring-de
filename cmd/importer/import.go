package main

import (
	"github.com/patrick246/impfquotenmonitoring-de/pkg/downloader"
	"github.com/patrick246/impfquotenmonitoring-de/pkg/persistence/mongodb"
	"log"
	"os"
	"time"
)

func main() {
	client, err := downloader.NewDownloader("https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Daten/Impfquotenmonitoring.xlsx?__blob=publicationFile")
	if err != nil {
		log.Fatalln(err)
	}

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

	data, err := client.Download()
	if err != nil {
		log.Fatalln(err)
	}

	quotes, err := downloader.Extract(data)
	if err != nil {
		log.Fatalln(err)
	}

	today := time.Now()
	for _, q := range quotes {
		err := storageService.Store(today, q.State, q.VaccinationData)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
