package recon

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type TransactionType string

const (
	Debit  TransactionType = "debit"
	Credit TransactionType = "credit"
)

type Transaction struct {
	Id     string
	Amount float64
	Type   TransactionType
	Time   time.Time
}

type TransactionStorage struct {
	destinationFileNamePath string
	destinationSheetName    string

	excelWriterFactory ExcelWriterFactory
}

func NewTransactionStorage(destinationFileNamePath string, destinationSheetName string, excelWriterFactory ExcelWriterFactory) TransactionStorage {
	return TransactionStorage{
		destinationFileNamePath: destinationFileNamePath,
		destinationSheetName:    destinationSheetName,
		excelWriterFactory:      excelWriterFactory,
	}
}

func (t TransactionStorage) StoreTransactions(transactions []Transaction) error {
	f, err := t.excelWriterFactory.New(t.destinationFileNamePath)
	if err != nil {
		return err
	}

	index, err := f.GetSheetIndex(t.destinationSheetName)
	if err != nil {
		return err
	}

	if index == -1 {
		f.NewSheet(t.destinationSheetName)
	}

	// header
	headers := []string{"Id", "Amount", "Type", "Time"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(t.destinationSheetName, cell, h)
	}

	// rows
	for row, tx := range transactions {
		values := []any{
			tx.Id,
			tx.Amount,
			string(tx.Type),
			tx.Time.Format(time.RFC3339),
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
			f.SetCellValue(t.destinationSheetName, cell, v)
		}
	}

	return f.SaveAs(t.destinationFileNamePath)
}

func (TransactionStorage) GetTransactions(filename string, startDate time.Time, endDate time.Time) ([]Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data rows found")
	}

	var transactions []Transaction
	// Skip header (records[0])
	for _, row := range records[1:] {
		amount, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount in row: %v", row)
		}

		t, err := time.Parse(time.RFC3339, row[3])
		if err != nil {
			return nil, fmt.Errorf("invalid time format in row: %v", row)
		}

		if t.Before(startDate) || t.After(endDate.Add(24*time.Hour)) {
			continue
		}

		tx := Transaction{
			Id:     row[0],
			Amount: amount,
			Type:   TransactionType(row[2]),
			Time:   t,
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}
