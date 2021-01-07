package downloader

import (
	"fmt"
	"github.com/patrick246/impfquotenmonitoring-de/pkg/persistence"
	"github.com/tealeg/xlsx/v3"
)

func Extract(data []byte) ([]Impfquote, error) {
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}

	if len(file.Sheets) < 2 {
		return nil, fmt.Errorf("too few sheets to extract data, got %d, expected at least 2", len(file.Sheets))
	}

	dataSheet := file.Sheets[1]

	var quotes []Impfquote

	for i := 1; i < 17; i++ {
		state, err := readCellString(dataSheet, i, 1)
		if err != nil {
			return nil, err
		}

		total, err := readCellInt64(dataSheet, i, 2)
		if err != nil {
			return nil, err
		}

		diffYesterday, err := readCellInt64(dataSheet, i, 3)
		if err != nil {
			return nil, err
		}

		per1k, err := readCellFloat64(dataSheet, i, 4)
		if err != nil {
			return nil, err
		}

		becauseAge, err := readCellInt64(dataSheet, i, 5)
		if err != nil {
			return nil, err
		}

		becauseJob, err := readCellInt64(dataSheet, i, 6)
		if err != nil {
			return nil, err
		}

		medReason, err := readCellInt64(dataSheet, i, 7)
		if err != nil {
			return nil, err
		}

		inCare, err := readCellInt64(dataSheet, i, 8)
		if err != nil {
			return nil, err
		}

		extractedRow := Impfquote{
			State: state,
			VaccinationData: persistence.VaccinationData{
				Total:                 total,
				DifferenceToYesterday: diffYesterday,
				Per1K:                 per1k,
				BecauseAge:            becauseAge,
				BecauseJob:            becauseJob,
				BecauseMedicalReasons: medReason,
				BecauseInCare:         inCare,
			},
		}
		quotes = append(quotes, extractedRow)
	}
	return quotes, nil
}

func readCellString(sheet *xlsx.Sheet, row, col int) (string, error) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return "", err
	}
	return cell.String(), nil
}

func readCellInt64(sheet *xlsx.Sheet, row, col int) (int64, error) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return 0, err
	}

	if cell.String() == "" {
		return -1, nil
	}

	value, err := cell.Int64()
	if err != nil {
		return 0, err
	}

	return value, nil
}

func readCellFloat64(sheet *xlsx.Sheet, row, col int) (float64, error) {
	cell, err := sheet.Cell(row, col)
	if err != nil {
		return 0, err
	}

	value, err := cell.Float()
	if err != nil {
		return 0, err
	}

	return value, nil
}
