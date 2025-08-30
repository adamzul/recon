package recon

import (
	"fmt"
	"time"

	"github.com/samber/lo"
)

type bankStatementDisrepancyGroup struct {
	Statements []BankStatement
}

func (b *bankStatementDisrepancyGroup) Add(statement BankStatement) {
	b.Statements = append(b.Statements, statement)
}

type ReconExecutor struct {
	transactionStorage       TransactionStorageProvider
	bankStatementRepoStorage BankStatementStorageProvider
	summaryRepoStorage       SummaryStorageProvider
}

func NewReconExecutor(transactionRepo TransactionStorageProvider, bankStatementRepo BankStatementStorageProvider, summaryRepo SummaryStorageProvider) ReconExecutor {
	return ReconExecutor{
		transactionStorage:       transactionRepo,
		bankStatementRepoStorage: bankStatementRepo,
		summaryRepoStorage:       summaryRepo,
	}
}

func (r ReconExecutor) Execute(transactionPath string, bankStatementPathArray []string, startDate time.Time, endDate time.Time) error {
	transactions, err := r.transactionStorage.GetTransactions(transactionPath, startDate, endDate)
	if err != nil {
		return fmt.Errorf("get transactions error: %w", err)
	}

	total := Summary{TotalProcessed: len(transactions)}

	bankStatementMap := map[float64]*BankStatementGroup{}
	for _, path := range bankStatementPathArray {
		statements, err := r.bankStatementRepoStorage.GetBankStatements(path, startDate, endDate)
		if err != nil {
			return fmt.Errorf("get bank statements error: %w", err)
		}

		for _, statement := range statements {
			if _, ok := bankStatementMap[statement.Amount]; !ok {
				bankStatementMap[statement.Amount] = &BankStatementGroup{}
			}
			bankStatementMap[statement.Amount].Add(statement)
			total.TotalAmountBankStatements += statement.Amount
		}
	}

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
				total.TotalProcessed++
				total.TotalUnmatched++
			}
		}
	}

	err = r.summaryRepoStorage.StoreSummary(total)
	if err != nil {
		return fmt.Errorf("store summary error: %w", err)
	}

	err = r.transactionStorage.StoreTransactions(transactionDiscrepancies)
	if err != nil {
		return fmt.Errorf("store transactions error: %w", err)
	}

	for bank, group := range bankStatementDisrepancy {
		err = r.bankStatementRepoStorage.StoreBankStatements(group.Statements, bank)
		if err != nil {
			return fmt.Errorf("store bank statements error: %w", err)
		}
	}

	return nil
}
