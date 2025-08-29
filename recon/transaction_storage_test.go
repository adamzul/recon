package recon

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

type transactionStorageSuite struct {
	mockExcelWriter        *MockExcelWriter
	mockExcelWriterFactory *MockExcelWriterFactory
	mockReader             *MockReader
	mockReaderFactory      *MockReaderFactory
	transactionStorage     TransactionStorage
}

func getTransactionStorageSuite(ctrl *gomock.Controller) transactionStorageSuite {
	mockExcelWriter := NewMockExcelWriter(ctrl)
	mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
	MockReaderFactory := NewMockReaderFactory(ctrl)
	MockReader := NewMockReader(ctrl)

	return transactionStorageSuite{
		mockExcelWriter:        mockExcelWriter,
		mockExcelWriterFactory: mockExcelWriterFactory,
		mockReader:             MockReader,
		mockReaderFactory:      MockReaderFactory,
		transactionStorage:     NewTransactionStorage("test.xlsx", "Transaction", mockExcelWriterFactory, MockReaderFactory),
	}
}

func TestTransactionStorage_StoreTransactions(t *testing.T) {
	destinationFileNamePath := "test.xlsx"
	destinationSheetName := "Transaction"
	t.Run("should store transactions successfully", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		transactions := []Transaction{
			{
				Id:     "1",
				Amount: 100.0,
				Type:   Credit,
				Time:   time.Now(),
			},
			{
				Id:     "2",
				Amount: 200.0,
				Type:   Debit,
				Time:   time.Now(),
			},
		}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(suite.mockExcelWriter, nil)
		suite.mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(-1, nil)
		suite.mockExcelWriter.EXPECT().NewSheet(destinationSheetName).Return(1, nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Id").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", "Amount").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C1", "Type").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D1", "Time").Return(nil)

		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", transactions[0].Id).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", transactions[0].Amount).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C2", string(transactions[0].Type)).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D2", transactions[0].Time.Format(time.RFC3339)).Return(nil)

		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A3", transactions[1].Id).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B3", transactions[1].Amount).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C3", string(transactions[1].Type)).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D3", transactions[1].Time.Format(time.RFC3339)).Return(nil)

		suite.mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(nil)

		err := suite.transactionStorage.StoreTransactions(transactions)

		g.Expect(err).Should(BeNil())
	})

	t.Run("should return error when excelWriterFactory.New returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		transactions := []Transaction{}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(nil, fmt.Errorf("new error"))

		err := suite.transactionStorage.StoreTransactions(transactions)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("new error"))
	})

	t.Run("should return error when f.SaveAs returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		transactions := []Transaction{
			{
				Id:     "1",
				Amount: 100.0,
				Type:   Credit,
				Time:   time.Now(),
			},
		}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(suite.mockExcelWriter, nil)
		suite.mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(-1, nil)
		suite.mockExcelWriter.EXPECT().NewSheet(destinationSheetName).Return(1, nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Id").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", "Amount").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C1", "Type").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D1", "Time").Return(nil)

		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", transactions[0].Id).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", transactions[0].Amount).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C2", string(transactions[0].Type)).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D2", transactions[0].Time.Format(time.RFC3339)).Return(nil)

		suite.mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(fmt.Errorf("save error"))

		err := suite.transactionStorage.StoreTransactions(transactions)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("save error"))
	})
}

func TestTransactionStorage_GetTransactions(t *testing.T) {
	filename := "test.xlsx"
	startDate, _ := time.Parse(time.DateOnly, "2025-08-28")
	endDate, _ := time.Parse(time.DateOnly, "2025-08-29")

	t.Run("should get transactions successfully", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		filename := "test.xlsx"

		mockRecords := [][]string{
			{"Id", "Amount", "Type", "Time"},
			{"1", "100.0", "credit", startDate.Format(time.RFC3339)},
			{"2", "200.0", "debit", endDate.Format(time.RFC3339)},
		}

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(suite.mockReader, nil)
		suite.mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		suite.mockReader.EXPECT().Close().Return(nil)

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(BeNil())
		g.Expect(transactions).Should(HaveLen(2))
		g.Expect(transactions[0].Id).Should(Equal("1"))
		g.Expect(transactions[0].Amount).Should(Equal(100.0))
		g.Expect(transactions[0].Type).Should(Equal(TransactionType("credit")))
		g.Expect(transactions[1].Id).Should(Equal("2"))
		g.Expect(transactions[1].Amount).Should(Equal(200.0))
		g.Expect(transactions[1].Type).Should(Equal(TransactionType("debit")))
	})

	t.Run("should return error when readerFactory.NewReader returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(nil, fmt.Errorf("new reader error"))

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("new reader error"))
		g.Expect(transactions).Should(BeNil())
	})

	t.Run("should return error when reader.ReadAll returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(suite.mockReader, nil)
		suite.mockReader.EXPECT().ReadAll().Return(nil, fmt.Errorf("read all error"))
		suite.mockReader.EXPECT().Close().Return(nil)

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("read all error"))
		g.Expect(transactions).Should(BeNil())
	})

	t.Run("should return empty slice when no data rows found", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		mockRecords := [][]string{
			{"Id", "Amount", "Type", "Time"},
		}

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(suite.mockReader, nil)
		suite.mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		suite.mockReader.EXPECT().Close().Return(nil)

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("no data rows found"))
		g.Expect(transactions).Should(BeNil())
	})

	t.Run("should return error when invalid amount in row", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		mockRecords := [][]string{
			{"Id", "Amount", "Type", "Time"},
			{"1", "invalid", "credit", startDate.Format(time.RFC3339)},
		}

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(suite.mockReader, nil)
		suite.mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		suite.mockReader.EXPECT().Close().Return(nil)

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(ContainSubstring("invalid amount in row"))
		g.Expect(transactions).Should(BeNil())
	})

	t.Run("should return error when invalid time format in row", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		mockRecords := [][]string{
			{"Id", "Amount", "Type", "Time"},
			{"1", "100.0", "credit", "invalid"},
		}

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(suite.mockReader, nil)
		suite.mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		suite.mockReader.EXPECT().Close().Return(nil)

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(ContainSubstring("invalid time format in row"))
		g.Expect(transactions).Should(BeNil())
	})

	t.Run("should skip transactions outside the date range", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getTransactionStorageSuite(ctrl)

		filename := "test.xlsx"
		outsideDate := startDate.Add(-time.Hour * 48)

		mockRecords := [][]string{
			{"Id", "Amount", "Type", "Time"},
			{"1", "100.0", "credit", outsideDate.Format(time.RFC3339)},
			{"2", "200.0", "debit", endDate.Format(time.RFC3339)},
		}

		suite.mockReaderFactory.EXPECT().NewReader(filename).Return(suite.mockReader, nil)
		suite.mockReader.EXPECT().ReadAll().Return(mockRecords, nil)
		suite.mockReader.EXPECT().Close().Return(nil)

		transactions, err := suite.transactionStorage.GetTransactions(filename, startDate, endDate)

		g.Expect(err).Should(BeNil())
		g.Expect(transactions).Should(HaveLen(1))
		g.Expect(transactions[0].Id).Should(Equal("2"))
	})
}
