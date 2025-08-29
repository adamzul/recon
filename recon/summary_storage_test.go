package recon

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

type SummaryStorageSuite struct {
	mockExcelWriter        *MockExcelWriter
	mockExcelWriterFactory *MockExcelWriterFactory
	summaryStorage         SummaryStorage
}

func summaryStorageSuite(ctrl *gomock.Controller) SummaryStorageSuite {
	mockExcelWriter := NewMockExcelWriter(ctrl)
	mockExcelWriterFactory := NewMockExcelWriterFactory(ctrl)

	return SummaryStorageSuite{
		mockExcelWriter:        mockExcelWriter,
		mockExcelWriterFactory: mockExcelWriterFactory,
		summaryStorage:         NewSummaryStorage("test.xlsx", "Summary", mockExcelWriterFactory),
	}
}

func TestSummaryStorage_StoreSummary(t *testing.T) {
	destinationFileNamePath := "test.xlsx"
	destinationSheetName := "Summary"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)
		suite := summaryStorageSuite(ctrl)

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(suite.mockExcelWriter, nil)
		suite.mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(1, nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Total Amount Transactions").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", summary.TotalAmountTransactions).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", "Total Amount Bank Statements").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", summary.TotalAmountBankStatements).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A3", "Total Matched").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B3", summary.TotalMatched).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A4", "Total Unmatched").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B4", summary.TotalUnmatched).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A5", "Total Processed").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B5", summary.TotalProcessed).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A6", "Total Amount Dicrepancy").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B6", summary.TotalAmountTransactions-summary.TotalAmountBankStatements).Return(nil)
		suite.mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(nil)

		err := suite.summaryStorage.StoreSummary(summary)

		g.Expect(err).Should(BeNil())
	})

	t.Run("excelize open file error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)
		suite := summaryStorageSuite(ctrl)

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(suite.mockExcelWriter, errors.New("open file error"))

		err := suite.summaryStorage.StoreSummary(summary)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(err.Error()).Should(Equal("open file error"))
	})

	t.Run("get sheet index error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)
		suite := summaryStorageSuite(ctrl)

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(suite.mockExcelWriter, nil)
		suite.mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(0, errors.New("get sheet index error"))

		err := suite.summaryStorage.StoreSummary(summary)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(err.Error()).Should(Equal("get sheet index error"))
	})

	t.Run("save as error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		g := NewGomegaWithT(t)
		suite := summaryStorageSuite(ctrl)

		summary := Summary{
			TotalAmountTransactions:   100.0,
			TotalAmountBankStatements: 90.0,
			TotalMatched:              50,
			TotalUnmatched:            50,
			TotalProcessed:            100,
		}

		suite.mockExcelWriterFactory.EXPECT().New(destinationFileNamePath).Return(suite.mockExcelWriter, nil)
		suite.mockExcelWriter.EXPECT().GetSheetIndex(destinationSheetName).Return(1, nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A1", "Total Amount Transactions").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B1", summary.TotalAmountTransactions).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A2", "Total Amount Bank Statements").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B2", summary.TotalAmountBankStatements).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A3", "Total Matched").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B3", summary.TotalMatched).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A4", "Total Unmatched").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B4", summary.TotalUnmatched).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A5", "Total Processed").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B5", summary.TotalProcessed).Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "A6", "Total Amount Dicrepancy").Return(nil)
		suite.mockExcelWriter.EXPECT().SetCellValue(destinationSheetName, "B6", summary.TotalAmountTransactions-summary.TotalAmountBankStatements).Return(nil)
		suite.mockExcelWriter.EXPECT().SaveAs(destinationFileNamePath).Return(errors.New("save as error"))

		err := suite.summaryStorage.StoreSummary(summary)

		g.Expect(err).ShouldNot(BeNil())
		g.Expect(err.Error()).Should(Equal("save as error"))
	})
}
