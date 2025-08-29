package main

import "github.com/xuri/excelize/v2"

type SummaryRepo struct {
}

func (SummaryRepo) WriteSummary(path, sheet string, total Summary) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		f = excelize.NewFile()
	}

	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return err
	}

	if index == -1 {
		f.NewSheet(sheet)
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
			f.SetCellValue(sheet, cell, v)
		}
	}

	return f.SaveAs(path)
}
