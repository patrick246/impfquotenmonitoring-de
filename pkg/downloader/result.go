package downloader

import "github.com/patrick246/impfquotenmonitoring-de/pkg/persistence"

type Impfquote struct {
	State string
	persistence.VaccinationData
}
