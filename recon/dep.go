//go:generate mockgen -typed -source=dep.go -destination=dep_mocks.go -package=recon
package recon

import (
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelWriterFactory interface {
	New(path string) (ExcelWriter, error)
}

type ExcelWriter interface {
	SetCellValue(sheet, axis string, value interface{}) error
	GetSheetIndex(name string) (int, error)
	NewSheet(name string) (int, error)
	SaveAs(name string, options ...excelize.Options) error
}

type ReaderFactory interface {
	NewReader(filename string) (Reader, error)
}

type Reader interface {
	ReadAll() ([][]string, error)
	Close() error
}

type TransactionStorageProvider interface {
	StoreTransactions(transactions []Transaction) error
	GetTransactions(filename string, startDate time.Time, endDate time.Time) ([]Transaction, error)
}

type BankStatementStorageProvider interface {
	StoreBankStatements(statements []BankStatement, bankName string) error
	GetBankStatements(filename string, startDate time.Time, endDate time.Time) ([]BankStatement, error)
}

type SummaryStorageProvider interface {
	StoreSummary(total Summary) error
}
