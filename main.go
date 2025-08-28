package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
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

type BankStatement struct {
	Bank   string
	ID     string
	Amount float64
	Time   time.Time
}

func main() {
	var transactionPath, bankStatementPaths string
	flag.StringVar(&transactionPath, "transaction-path", "", "transactions CSV file path")
	flag.StringVar(&bankStatementPaths, "bank-statement-paths", "", "bank statements CSV file path")

	flag.Parse()

	transactions, err := readTransactionsFromCSV(transactionPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	bankStatementMap := map[float64][]BankStatement{}
	bankStatementPathArray := strings.Split(bankStatementPaths, ",")
	for _, path := range bankStatementPathArray {
		statements, err := readBankStatementsFromCSV(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, statement := range statements {
			if bankStatementMap[statement.Amount] == nil {
				bankStatementMap[statement.Amount] = []BankStatement{}
			}
			bankStatementMap[statement.Amount] = append(bankStatementMap[statement.Amount], statement)
		}
	}

	//recon
	transactions = lo.Filter(transactions, func(t Transaction, i int) bool {
		bankStatements := bankStatementMap[t.Amount]
		if len(bankStatements) > 0 {
			bankStatements = bankStatements[1:]
			bankStatementMap[t.Amount] = bankStatements
			// remove the transaction if it has a matching bank statement
			return false
		}
		return true
	})
	fmt.Println(transactions)
	fmt.Println(bankStatementMap)
}

func readTransactionsFromCSV(filename string) ([]Transaction, error) {
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

func readBankStatementsFromCSV(filename string) ([]BankStatement, error) {
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
			// fallback example: try "2006-01-02 15:04:05"
			t, err = time.Parse("2006-01-02 15:04:05", row[2])
			if err != nil {
				return nil, fmt.Errorf("invalid time format in row: %v", row)
			}
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
