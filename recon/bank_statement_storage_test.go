package recon

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

func TestBankStatementStorage_GetBankStatements(t *testing.T) {
	filename := "test.csv"
	startDate, _ := time.Parse(time.DateOnly, "2025-08-28")
	endDate, _ := time.Parse(time.DateOnly, "2025-08-29")

	t.Run("should get bank statements successfully", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)
		mockReader := NewMockReader(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		mockRecords := [][]string{
			{"ID", "Amount", "Time"},
			{"1", "100.0", startDate.Format(time.RFC3339)},
			{"2", "200.0", endDate.Format(time.RFC3339)},
		}

		mockReaderFactory.EXPECT().NewReader(filename).Return(mockReader, nil)
		mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		mockReader.EXPECT().Close().Return(nil)

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).Should(BeNil())
		g.Expect(statements).Should(HaveLen(2))
		g.Expect(statements[0].ID).Should(Equal("1"))
		g.Expect(statements[0].Amount).Should(Equal(100.0))
		g.Expect(statements[1].ID).Should(Equal("2"))
		g.Expect(statements[1].Amount).Should(Equal(200.0))
	})

	t.Run("should return error when readerFactory.NewReader returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		mockReaderFactory.EXPECT().NewReader(filename).Return(nil, fmt.Errorf("new reader error"))

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(statements).Should(BeNil())
	})

	t.Run("should return error when reader.ReadAll returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)
		mockReader := NewMockReader(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		mockReaderFactory.EXPECT().NewReader(filename).Return(mockReader, nil)
		mockReader.EXPECT().ReadAll().Return(nil, fmt.Errorf("read all error"))
		mockReader.EXPECT().Close().Return(nil)

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(statements).Should(BeNil())
	})

	t.Run("should return error when no data rows found", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)
		mockReader := NewMockReader(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		mockRecords := [][]string{
			{"ID", "Amount", "Time"},
		}

		mockReaderFactory.EXPECT().NewReader(filename).Return(mockReader, nil)
		mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		mockReader.EXPECT().Close().Return(nil)

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(statements).Should(BeNil())
	})

	t.Run("should return error when invalid amount in row", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)
		mockReader := NewMockReader(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		mockRecords := [][]string{
			{"ID", "Amount", "Time"},
			{"1", "invalid", startDate.Format(time.RFC3339)},
		}

		mockReaderFactory.EXPECT().NewReader(filename).Return(mockReader, nil)
		mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		mockReader.EXPECT().Close().Return(nil)

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(statements).Should(BeNil())
	})

	t.Run("should return error when invalid time format in row", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)
		mockReader := NewMockReader(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		mockRecords := [][]string{
			{"ID", "Amount", "Time"},
			{"1", "100.0", "invalid"},
		}

		mockReaderFactory.EXPECT().NewReader(filename).Return(mockReader, nil)
		mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		mockReader.EXPECT().Close().Return(nil)

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(statements).Should(BeNil())
	})

	t.Run("should skip transactions outside the date range", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockReaderFactory := NewMockReaderFactory(ctrl)
		mockReader := NewMockReader(ctrl)

		bankStatementStorage := NewBankStatementStorage("test.xlsx", nil, mockReaderFactory)

		outsideDate := startDate.Add(-time.Hour * 48)

		mockRecords := [][]string{
			{"ID", "Amount", "Time"},
			{"1", "100.0", outsideDate.Format(time.RFC3339)},
			{"2", "200.0", endDate.Format(time.RFC3339)},
		}

		mockReaderFactory.EXPECT().NewReader(filename).Return(mockReader, nil)
		mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		mockReader.EXPECT().Close().Return(nil)

		statements, err := bankStatementStorage.GetBankStatements(filename, startDate, endDate)

		g.Expect(err).Should(BeNil())
		g.Expect(statements).Should(HaveLen(1))
		g.Expect(statements[0].ID).Should(Equal("2"))
	})
}

func TestBankStatementStorage_StoreBankStatements(t *testing.T) {
	destinationFileNamePath := "test.xlsx"

	t.Run("should store bank statements successfully", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
		mockExcelWriter := NewMockExcelWriter(ctrl)

		bankStatementStorage := NewBankStatementStorage(destinationFileNamePath, mockExcelWriterFactory, nil)

		statements := []BankStatement{
			{
				Bank:   "BankA",
				ID:     "1",
				Amount: 100.0,
				Time:   time.Now(),
			},
			{
				Bank:   "BankB",
				ID:     "2",
				Amount: 200.0,
				Time:   time.Now(),
			},
		}
		bankName := "BankA"

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(bankName).Return(-1, nil)
		mockExcelWriter.EXPECT().NewSheet(bankName).Return(1, nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "A1", "Bank").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "B1", "ID").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "C1", "Amount").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "D1", "Time").Return(nil)

		mockExcelWriter.EXPECT().SetCellValue(bankName, "A2", statements[0].Bank).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "B2", statements[0].ID).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "C2", statements[0].Amount).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "D2", statements[0].Time.Format(time.RFC3339)).Return(nil)

		mockExcelWriter.EXPECT().SetCellValue(bankName, "A3", statements[1].Bank).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "B3", statements[1].ID).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "C3", statements[1].Amount).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "D3", statements[1].Time.Format(time.RFC3339)).Return(nil)

		mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(nil)

		err := bankStatementStorage.StoreBankStatements(statements, bankName)

		g.Expect(err).Should(BeNil())
	})

	t.Run("should return error when excelWriterFactory.New returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

		bankStatementStorage := NewBankStatementStorage(destinationFileNamePath, mockExcelWriterFactory, nil)

		statements := []BankStatement{}
		bankName := "BankA"

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(nil, fmt.Errorf("new error"))

		err := bankStatementStorage.StoreBankStatements(statements, bankName)

		g.Expect(err).ShouldNot(BeNil())
	})

	t.Run("should return error when f.GetSheetIndex returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
		mockExcelWriter := NewMockExcelWriter(ctrl)

		bankStatementStorage := NewBankStatementStorage(destinationFileNamePath, mockExcelWriterFactory, nil)

		statements := []BankStatement{}
		bankName := "BankA"

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(bankName).Return(-1, fmt.Errorf("get sheet index error"))

		err := bankStatementStorage.StoreBankStatements(statements, bankName)

		g.Expect(err).ShouldNot(BeNil())
	})

	t.Run("should return error when f.NewSheet returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
		mockExcelWriter := NewMockExcelWriter(ctrl)

		bankStatementStorage := NewBankStatementStorage(destinationFileNamePath, mockExcelWriterFactory, nil)

		statements := []BankStatement{}
		bankName := "BankA"

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(bankName).Return(-1, nil)
		mockExcelWriter.EXPECT().NewSheet(bankName).Return(1, fmt.Errorf("new sheet error"))

		err := bankStatementStorage.StoreBankStatements(statements, bankName)

		g.Expect(err).ShouldNot(BeNil())
	})

	t.Run("should return error when f.SaveAs returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
		mockExcelWriter := NewMockExcelWriter(ctrl)

		bankStatementStorage := NewBankStatementStorage(destinationFileNamePath, mockExcelWriterFactory, nil)

		statements := []BankStatement{
			{
				Bank:   "BankA",
				ID:     "1",
				Amount: 100.0,
				Time:   time.Now(),
			},
		}
		bankName := "BankA"

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(bankName).Return(-1, nil)
		mockExcelWriter.EXPECT().NewSheet(bankName).Return(1, nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "A1", "Bank").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "B1", "ID").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "C1", "Amount").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "D1", "Time").Return(nil)

		mockExcelWriter.EXPECT().SetCellValue(bankName, "A2", statements[0].Bank).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "B2", statements[0].ID).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "C2", statements[0].Amount).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(bankName, "D2", statements[0].Time.Format(time.RFC3339)).Return(nil)

		mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(fmt.Errorf("save error"))

		err := bankStatementStorage.StoreBankStatements(statements, bankName)

		g.Expect(err).ShouldNot(BeNil())
	})
}
