//go:generate mockgen -typed -source=dep.go -destination=dep_mocks.go -package=recon
package recon

import "github.com/xuri/excelize/v2"

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
