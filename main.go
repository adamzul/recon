package main

import (
	"flag"
	"log"
	"time"
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

type Summary struct {
	TotalAmountTransactions   float64
	TotalAmountBankStatements float64
	TotalMatched              int
	TotalUnmatched            int
	TotalProcessed            int
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

func main() {
	var transactionPath, bankStatementPaths string
	var startDateStr, endDateStr string
	flag.StringVar(&transactionPath, "transaction-path", "transaction.csv", "transactions CSV file path")
	flag.StringVar(&bankStatementPaths, "bank-statement-paths", "bca.csv,bri.csv", "bank statements CSV file path")
	flag.StringVar(&startDateStr, "start-date", time.Now().Format("2006-01-02"), "bank statements CSV file path")
	flag.StringVar(&endDateStr, "end-date", time.Now().Format("2006-01-02"), "bank statements CSV file path")
	flag.Parse()

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Println(err)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Println(err)
		return
	}

	reconExecutor := ReconExecutor{
		transactionRepo:   TransactionRepo{},
		bankStatementRepo: BankStatementRepo{},
		summaryRepo:       SummaryRepo{},
	}
	reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
}
