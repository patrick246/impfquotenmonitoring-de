package mongodb

import (
	"context"
	"github.com/patrick246/impfquotenmonitoring-de/pkg/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type StorageService struct {
	client *Client
}

const collectionName = "vaccinemetrics"

func NewStorageService(mongodbClient *Client) (*StorageService, error) {
	return &StorageService{client: mongodbClient}, nil
}

func (s *StorageService) Store(date time.Time, state string, data persistence.VaccinationData) error {
	month := date.Format("01-2006")
	docId := month + "-" + state
	day := date.Format("2")

	query := bson.D{{
		"_id", docId,
	}}

	upsertUpdate := bson.D{{
		"$set", bson.D{{
			"month", month,
		}, {
			"state", state,
		}},
	}, {
		"$setOnInsert", bson.D{{
			"days", bson.A{},
		}},
	}}

	updateOptions := options.Update().SetUpsert(true)

	_, err := s.client.Collection(collectionName).UpdateOne(context.Background(), query, upsertUpdate, updateOptions)
	if err != nil {
		return err
	}

	dataInsertUpdate := bson.D{{
		"$set", bson.D{{
			"days." + day, data,
		}},
	}}

	_, err = s.client.Collection(collectionName).UpdateOne(context.Background(), query, dataInsertUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) GetMonths(from time.Time, to time.Time) ([]persistence.VaccinationMonthMetric, error) {
	monthFrom := from.Format("01-2006")
	monthTo := to.Format("01-2006")

	query := bson.D{{
		"month", bson.D{{
			"$gte", monthFrom,
		}, {
			"$lte", monthTo,
		}},
	}}

	result, err := s.client.Collection(collectionName).Find(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var data []VaccineMonthMetricModel
	err = result.All(context.Background(), &data)
	if err != nil {
		return nil, err
	}

	var generalData []persistence.VaccinationMonthMetric
	for _, d := range data {
		generalData = append(generalData, toVaccineMonthMetric(d))
	}

	return generalData, nil
}
