package recon

import "github.com/xuri/excelize/v2"

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
}

func NewSummaryStorage(destinationFileNamePath string, destinationSheetName string) SummaryStorage {
	return SummaryStorage{
		destinationFileNamePath: destinationFileNamePath,
		destinationSheetName:    destinationSheetName,
	}
}

func (s SummaryStorage) StoreSummary(total Summary) error {
	f, err := excelize.OpenFile(s.destinationFileNamePath)
	if err != nil {
		f = excelize.NewFile()
	}

	index, err := f.GetSheetIndex(s.destinationSheetName)
	if err != nil {
		return err
	}

	if index == -1 {
		f.NewSheet(s.destinationSheetName)
	}

	// key-value pairs
	rows := [][]interface{}{
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
