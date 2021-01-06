package mongodb

import "github.com/patrick246/impfquotenmonitoring-de/pkg/persistence"

type VaccineMonthMetricModel struct {
	ID    string                         `bson:"_id"`
	State string                         `json:"state"`
	Month string                         `json:"month"`
	Days  []*persistence.VaccinationData `json:"days"`
}

func fromVaccineMonthMetric(vmm persistence.VaccinationMonthMetric) VaccineMonthMetricModel {
	id := vmm.Month + "-" + vmm.State
	return VaccineMonthMetricModel{
		ID:    id,
		State: vmm.State,
		Month: vmm.Month,
		Days:  vmm.Days,
	}
}

func toVaccineMonthMetric(model VaccineMonthMetricModel) persistence.VaccinationMonthMetric {
	return persistence.VaccinationMonthMetric{
		State: model.State,
		Month: model.Month,
		Days:  model.Days,
	}
}
