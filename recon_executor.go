package main

import (
	"log"
	"strings"
	"time"

	"github.com/samber/lo"
)

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
