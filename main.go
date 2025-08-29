package main

import (
	"flag"
	"log"
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

type ReconExecutor struct {
	transactionRepo   TransactionRepo
	bankStatementRepo BankStatementRepo
	summaryRepo       SummaryRepo
}

func (r ReconExecutor) Execute(transactionPath, bankStatementPaths string, startDate, endDate time.Time) {
	transactions, err := r.transactionRepo.GetTransactions(transactionPath, startDate, endDate)
	if err != nil {
		log.Println(err)
		return
	}

	total := Summary{TotalProcessed: len(transactions)}

	bankStatementMap := map[float64]*BankStatementGroup{}
	bankStatementPathArray := strings.Split(bankStatementPaths, ",")
	for _, path := range bankStatementPathArray {
		statements, err := r.bankStatementRepo.GetBankStatements(path, startDate, endDate)
		if err != nil {
			log.Println(err)
			return
		}

		for _, statement := range statements {
			if _, ok := bankStatementMap[statement.Amount]; !ok {
				bankStatementMap[statement.Amount] = &BankStatementGroup{}
			}
			bankStatementMap[statement.Amount].Add(statement)
			total.TotalAmountBankStatements += statement.Amount
		}
	}

	//recon
	transactionDiscrepancies := lo.Filter(transactions, func(t Transaction, i int) bool {
		total.TotalAmountTransactions += t.Amount
		bankStatements := bankStatementMap[t.Amount]
		if bankStatements != nil && len(bankStatements.BankStatements) > 0 {
			bankStatementMap[t.Amount].Shift()
			total.TotalMatched++
			// remove the transaction if it has a matching bank statement
			return false
		}
		total.TotalUnmatched++
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
				total.TotalProcessed++
				total.TotalUnmatched++
			}
		}
	}

	r.summaryRepo.WriteSummary("recon.xlsx", "summary", total)

	r.transactionRepo.WriteTransactions("recon.xlsx", transactionDiscrepancies, "transaction")
	for bank, group := range bankStatementDisrepancy {
		r.bankStatementRepo.WriteBankStatements("recon.xlsx", group.Statements, group.AppearMultiple, bank)
	}
}
