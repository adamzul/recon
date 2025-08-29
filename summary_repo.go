package main

import "github.com/xuri/excelize/v2"

type Summary struct {
	TotalAmountTransactions   float64
	TotalAmountBankStatements float64
	TotalMatched              int
	TotalUnmatched            int
	TotalProcessed            int
}

type SummaryRepo struct {
	fileNamePath string
	sheetName    string
}

func (s SummaryRepo) WriteSummary(total Summary) error {
	f, err := excelize.OpenFile(s.fileNamePath)
	if err != nil {
		f = excelize.NewFile()
	}

	index, err := f.GetSheetIndex(s.sheetName)
	if err != nil {
		return err
	}

	if index == -1 {
		f.NewSheet(s.sheetName)
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
			f.SetCellValue(s.sheetName, cell, v)
		}
	}

	return f.SaveAs(s.fileNamePath)
}
