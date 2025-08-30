package recon

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Summary struct {
	TotalAmountTransactions   float64
	TotalAmountBankStatements float64
	TotalMatched              int
	TotalUnmatched            int
	TotalProcessed            int
}

type SummaryStorage struct {
	destinationFileNamePath string
	destinationSheetName    string
	excelWriterFactory      ExcelWriterFactory
}

func NewSummaryStorage(destinationFileNamePath string, destinationSheetName string, excelWriterFactory ExcelWriterFactory) SummaryStorage {
	return SummaryStorage{
		destinationFileNamePath: destinationFileNamePath,
		destinationSheetName:    destinationSheetName,
		excelWriterFactory:      excelWriterFactory,
	}
}

func (s SummaryStorage) StoreSummary(total Summary) error {
	f, err := s.excelWriterFactory.New(s.destinationFileNamePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	index, err := f.GetSheetIndex(s.destinationSheetName)
	if err != nil {
		return fmt.Errorf("failed to get sheet index: %w", err)
	}

	if index == -1 {
		_, err = f.NewSheet(s.destinationSheetName)
		if err != nil {
			return fmt.Errorf("failed to create sheet: %w", err)
		}
	}

	// key-value pairs
	rows := [][]any{
		{"Total Amount Transactions", total.TotalAmountTransactions},
		{"Total Amount Bank Statements", total.TotalAmountBankStatements},
		{"Total Matched", total.TotalMatched},
		{"Total Unmatched", total.TotalUnmatched},
		{"Total Processed", total.TotalProcessed},
		{"Total Amount Dicrepancy", total.TotalAmountTransactions - total.TotalAmountBankStatements},
	}

	for i, row := range rows {
		for j, v := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			f.SetCellValue(s.destinationSheetName, cell, v)
		}
	}

	return f.SaveAs(s.destinationFileNamePath)
}
