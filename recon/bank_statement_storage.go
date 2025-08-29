package recon

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type BankStatement struct {
	Bank   string
	ID     string
	Amount float64
	Time   time.Time
}

type BankStatementGroup struct {
	BankStatements []BankStatement
	AppearMultiple bool
}

func (b *BankStatementGroup) Add(statement BankStatement) {
	if len(b.BankStatements) == 0 {
		b.BankStatements = []BankStatement{}
	} else {
		b.SetAppearMultiple()
	}

	b.BankStatements = append(b.BankStatements, statement)
}

func (b *BankStatementGroup) Shift() {
	b.BankStatements = b.BankStatements[1:]
}

func (b *BankStatementGroup) SetAppearMultiple() {
	b.AppearMultiple = true
}

type BankStatementStorage struct {
	destinationFileNamePath string

	excelWriterFactory ExcelWriterFactory
	readerFactory      ReaderFactory
}

func NewBankStatementStorage(destinationFileNamePath string, excelWriterFactory ExcelWriterFactory, readerFactory ReaderFactory) BankStatementStorage {
	return BankStatementStorage{
		destinationFileNamePath: destinationFileNamePath,
		excelWriterFactory:      excelWriterFactory,
		readerFactory:           readerFactory,
	}
}

func (b BankStatementStorage) GetBankStatements(filename string, startDate time.Time, endDate time.Time) ([]BankStatement, error) {
	reader, err := b.readerFactory.NewReader(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data rows found in %s", filename)
	}

	bankName := filepath.Base(filename) // extract filename only, e.g. "bca.csv"
	bankName = strings.TrimSuffix(bankName, filepath.Ext(bankName))

	var statements []BankStatement
	for _, row := range records[1:] { // skip header
		if len(row) < 3 {
			continue
		}

		amount, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount in row: %v", row)
		}

		t, err := time.Parse(time.RFC3339, row[2])
		if err != nil {
			return nil, fmt.Errorf("invalid time format in row: %v", row)
		}

		if t.Before(startDate) || t.After(endDate.Add(24*time.Hour)) {
			continue
		}

		statements = append(statements, BankStatement{
			Bank:   bankName,
			ID:     row[0],
			Amount: amount,
			Time:   t,
		})
	}

	return statements, nil
}

func (b BankStatementStorage) StoreBankStatements(statements []BankStatement, appearMultiple bool, bankName string) error {
	f, err := b.excelWriterFactory.New(b.destinationFileNamePath) // open existing file
	if err != nil {
		return err
	}

	index, err := f.GetSheetIndex(bankName)
	if err != nil {
		return err
	}

	if index == -1 {
		_, errSheet := f.NewSheet(bankName)
		if errSheet != nil {
			return errSheet
		}
	}

	// Write header row
	headers := []string{"Bank", "ID", "Amount", "Time", "Appear Multiple Time"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // row 1
		f.SetCellValue(bankName, cell, h)
	}

	// Write each BankStatement
	for row, s := range statements {
		values := []interface{}{
			s.Bank,
			s.ID,
			s.Amount,
			s.Time.Format(time.RFC3339), // store as formatted string
			appearMultiple,
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2) // data starts at row 2
			f.SetCellValue(bankName, cell, v)
		}
	}

	// Save file
	if err := f.SaveAs(b.destinationFileNamePath); err != nil {
		return err
	}
	return nil
}
