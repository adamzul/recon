package main

import (
	"flag"
	"fmt"
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
				bankStatementDisrepancy[statement.Bank].SetAppearMultiple(group.AppearMultiple)
			}
		}
	}

	writeTransactionsToExcel("recon.xlsx", transactions, "transaction")
	for bank, group := range bankStatementDisrepancy {
		writeBankStatementsToExcel("recon.xlsx", group.Statements, group.AppearMultiple, bank)
	}
}

type bankStatementDisrepancyGroup struct {
	Statements     []BankStatement
	AppearMultiple bool
}

func (b *bankStatementDisrepancyGroup) Add(statement BankStatement) {
	b.Statements = append(b.Statements, statement)
}

func (b *bankStatementDisrepancyGroup) SetAppearMultiple(isAppearMultipleTime bool) {
	b.AppearMultiple = isAppearMultipleTime
}
