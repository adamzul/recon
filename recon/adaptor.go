package recon

import (
	"encoding/csv"
	"os"

	"github.com/xuri/excelize/v2"
)

type ExcelWrapper struct {
	*excelize.File
}

type ExcelFactory struct{}

func (ExcelFactory) New(path string) (ExcelWriter, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return &ExcelWrapper{excelize.NewFile()}, nil
	}
	return &ExcelWrapper{f}, nil
}

type CSVReader struct {
	*csv.Reader
	*os.File
}

type CSVReaderFactory struct{}

func (CSVReaderFactory) NewReader(filename string) (Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	return &CSVReader{reader, file}, nil
}
