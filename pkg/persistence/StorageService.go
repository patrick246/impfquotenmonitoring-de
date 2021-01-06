package persistence

import (
	"time"
)

type StorageService interface {
	Store(day time.Time, state string, data VaccinationData) error
	GetMonths(from time.Time, to time.Time) ([]VaccinationMonthMetric, error)
}

type VaccinationMonthMetric struct {
	State string             `json:"state"`
	Month string             `json:"month"`
	Days  []*VaccinationData `json:"days"`
}

type VaccinationData struct {
	Total                 int64   `json:"total"`
	DifferenceToYesterday int64   `json:"difference_to_yesterday"`
	Per1K                 float64 `json:"per_1k"`
	BecauseAge            int64   `json:"because_age"`
	BecauseJob            int64   `json:"because_job"`
	BecauseMedicalReasons int64   `json:"because_medical_reasons"`
	BecauseInCare         int64   `json:"because_in_care"`
}
