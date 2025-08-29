package recon

import "github.com/xuri/excelize/v2"

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
