package recon

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

func TestTransactionStorage_StoreTransactions(t *testing.T) {
	t.Run("should store transactions successfully", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
		mockExcelWriter := NewMockExcelWriter(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Sheet1"

		transactionStorage := NewTransactionStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)

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

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(-1, nil)
		mockExcelWriter.EXPECT().NewSheet(destinationSheetName).Return(1, nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Id").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", "Amount").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C1", "Type").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D1", "Time").Return(nil)

		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", transactions[0].Id).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", transactions[0].Amount).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C2", string(transactions[0].Type)).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D2", transactions[0].Time.Format(time.RFC3339)).Return(nil)

		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A3", transactions[1].Id).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B3", transactions[1].Amount).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C3", string(transactions[1].Type)).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D3", transactions[1].Time.Format(time.RFC3339)).Return(nil)

		mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(nil)

		err := transactionStorage.StoreTransactions(transactions)

		g.Expect(err).Should(BeNil())
	})

	t.Run("should return error when excelWriterFactory.New returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Sheet1"

		transactionStorage := NewTransactionStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)

		transactions := []Transaction{}

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(nil, fmt.Errorf("new error"))

		err := transactionStorage.StoreTransactions(transactions)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("new error"))
	})

	t.Run("should return error when f.SaveAs returns error", func(t *testing.T) {
		g := NewGomegaWithT(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)
		mockExcelWriter := NewMockExcelWriter(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Sheet1"

		transactionStorage := NewTransactionStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)

		transactions := []Transaction{
			{
				Id:     "1",
				Amount: 100.0,
				Type:   Credit,
				Time:   time.Now(),
			},
		}

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(-1, nil)
		mockExcelWriter.EXPECT().NewSheet(destinationSheetName).Return(1, nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Id").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", "Amount").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C1", "Type").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D1", "Time").Return(nil)

		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", transactions[0].Id).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", transactions[0].Amount).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "C2", string(transactions[0].Type)).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "D2", transactions[0].Time.Format(time.RFC3339)).Return(nil)

		mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(fmt.Errorf("save error"))

		err := transactionStorage.StoreTransactions(transactions)

		g.Expect(err).Should(Not(BeNil()))
		g.Expect(err.Error()).Should(Equal("save error"))
	})
}
