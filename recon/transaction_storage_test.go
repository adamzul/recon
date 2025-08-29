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
