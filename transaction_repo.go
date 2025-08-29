package main

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

type TransactionRepo struct {
}

func (TransactionRepo) WriteTransactions(path string, transactions []Transaction, sheet string) error {
	f, err := excelize.OpenFile(path) // open existing file
	if err != nil {
		// if not exist, create new
		f = excelize.NewFile()
	}

	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return err
	}

	if index == -1 {
		f.NewSheet(sheet)
	}

	// header
	headers := []string{"Id", "Amount", "Type", "Time"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// rows
	for row, tx := range transactions {
		values := []interface{}{
			tx.Id,
			tx.Amount,
			string(tx.Type),
			tx.Time.Format(time.RFC3339),
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	return f.SaveAs(path)
}
func (TransactionRepo) GetTransactions(filename string, startDate time.Time, endDate time.Time) ([]Transaction, error) {
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
