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

	bankStatementMap := map[float64]*BankStatementGroup{}
	bankStatementPathArray := strings.Split(bankStatementPaths, ",")
	for _, path := range bankStatementPathArray {
		statements, err := readBankStatementsFromCSV(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, statement := range statements {
			if _, ok := bankStatementMap[statement.Amount]; !ok {
				bankStatementMap[statement.Amount] = &BankStatementGroup{}
			}
			bankStatementMap[statement.Amount].Add(statement)
		}
	}

	//recon
	transactions = lo.Filter(transactions, func(t Transaction, i int) bool {
		bankStatements := bankStatementMap[t.Amount]
		if bankStatements != nil && len(bankStatements.BankStatements) > 0 {
			bankStatementMap[t.Amount].Shift()
			// remove the transaction if it has a matching bank statement
			return false
		}
		return true
	})

	bankStatementDisrepancy := map[string]*bankStatementDisrepancyGroup{}
	for _, group := range bankStatementMap {
		if len(group.BankStatements) > 0 {
			for _, statement := range group.BankStatements {
				if _, ok := bankStatementDisrepancy[statement.Bank]; !ok {
					bankStatementDisrepancy[statement.Bank] = &bankStatementDisrepancyGroup{}
				}
				bankStatementDisrepancy[statement.Bank].Add(statement)
			}
		}
	}

	writeTransactionsToExcel("recon.xlsx", transactions, "transaction")
	for bank, group := range bankStatementDisrepancy {
		writeBankStatementsToExcel("recon.xlsx", group.Statements, bank)
	}
}

type bankStatementDisrepancyGroup struct {
	Statements []BankStatement
}

func (b *bankStatementDisrepancyGroup) Add(statement BankStatement) {
	b.Statements = append(b.Statements, statement)
}

func writeTransactionsToExcel(path string, transactions []Transaction, sheet string) error {
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

func writeBankStatementsToExcel(path string, statements []BankStatement, sheet string) error {
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

	// Write header row
	headers := []string{"Bank", "ID", "Amount", "Time"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // row 1
		f.SetCellValue(sheet, cell, h)
	}

	// Write each BankStatement
	for row, s := range statements {
		values := []interface{}{
			s.Bank,
			s.ID,
			s.Amount,
			s.Time.Format(time.RFC3339), // store as formatted string
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2) // data starts at row 2
			f.SetCellValue(sheet, cell, v)
		}
	}

	// Save file
	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}
