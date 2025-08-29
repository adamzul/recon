package recon

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestSummaryStorage_StoreSummary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)

		mockExcelWriter := NewMockExcelWriter(ctrl)
		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Summary"

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(1, nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Total Amount Transactions").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", summary.TotalAmountTransactions).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", "Total Amount Bank Statements").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", summary.TotalAmountBankStatements).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A3", "Total Matched").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B3", summary.TotalMatched).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A4", "Total Unmatched").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B4", summary.TotalUnmatched).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A5", "Total Processed").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B5", summary.TotalProcessed).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A6", "Total Amount Dicrepancy").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B6", summary.TotalAmountTransactions-summary.TotalAmountBankStatements).Return(nil)
		mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(nil)

		summaryStorage := NewSummaryStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)
		err := summaryStorage.StoreSummary(summary)

		g.Expect(err).Should(BeNil())
	})

	t.Run("excelize open file error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)

		mockExcelWriter := NewMockExcelWriter(ctrl)
		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Summary"

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, errors.New("open file error"))

		summaryStorage := NewSummaryStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)
		err := summaryStorage.StoreSummary(summary)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(err.Error()).Should(Equal("open file error"))
	})

	t.Run("get sheet index error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)

		mockExcelWriter := NewMockExcelWriter(ctrl)
		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Summary"

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(0, errors.New("get sheet index error"))

		summaryStorage := NewSummaryStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)
		err := summaryStorage.StoreSummary(summary)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(err.Error()).Should(Equal("get sheet index error"))
	})

	t.Run("save as error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)

		mockExcelWriter := NewMockExcelWriter(ctrl)
		mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

		destinationFileNamePath := "test.xlsx"
		destinationSheetName := "Summary"

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(mockExcelWriter, nil)
		mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(1, nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Total Amount Transactions").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", summary.TotalAmountTransactions).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", "Total Amount Bank Statements").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", summary.TotalAmountBankStatements).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A3", "Total Matched").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B3", summary.TotalMatched).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A4", "Total Unmatched").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B4", summary.TotalUnmatched).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A5", "Total Processed").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B5", summary.TotalProcessed).Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A6", "Total Amount Dicrepancy").Return(nil)
		mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B6", summary.TotalAmountTransactions-summary.TotalAmountBankStatements).Return(nil)
		mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(errors.New("save as error"))

		summaryStorage := NewSummaryStorage(destinationFileNamePath, destinationSheetName, mockExcelWriterFactory)
		err := summaryStorage.StoreSummary(summary)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(err.Error()).Should(Equal("save as error"))
	})
}
