//go:generate mockgen -typed -source=summary_storage.go -destination=summary_storage_mocks.go -package=recon
package recon

import "github.com/xuri/excelize/v2"

type ExcelWriterFactory interface {
	New(path string) (ExcelWriter, error)
}

type ExcelWriter interface {
	SetCellValue(sheet, axis string, value interface{}) error
	GetSheetIndex(name string) (int, error)
	NewSheet(name string) (int, error)
	SaveAs(name string, options ...excelize.Options) error
}



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
		return err
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
