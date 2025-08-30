package recon

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

type reconExecutorSuite struct {
	mockTransactionStorage       *MockTransactionStorageProvider
	mockBankStatementRepoStorage *MockBankStatementStorageProvider
	mockSummaryRepoStorage       *MockSummaryStorageProvider
	reconExecutor                ReconExecutor
}

func getReconExecutorSuite(ctrl *gomock.Controller) reconExecutorSuite {
	mockTransactionStorage := NewMockTransactionStorageProvider(ctrl)
	mockBankStatementRepoStorage := NewMockBankStatementStorageProvider(ctrl)
	mockSummaryRepoStorage := NewMockSummaryStorageProvider(ctrl)

	return reconExecutorSuite{
		mockTransactionStorage:       mockTransactionStorage,
		mockBankStatementRepoStorage: mockBankStatementRepoStorage,
		mockSummaryRepoStorage:       mockSummaryRepoStorage,
		reconExecutor:                NewReconExecutor(mockTransactionStorage, mockBankStatementRepoStorage, mockSummaryRepoStorage),
	}
}

func TestReconExecutor_Execute(t *testing.T) {
	transactionPath := "transaction.xlsx"
	bankStatementPaths := []string{"bca.xlsx", "bri.xlsx"}
	startDate, _ := time.Parse(time.DateOnly, "2025-08-01")
	endDate, _ := time.Parse(time.DateOnly, "2025-08-30")

	t.Run("should execute recon successfully", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getReconExecutorSuite(ctrl)

		transactions := []Transaction{
			{ID: "1", Amount: 100.0, Type: Credit, Time: startDate},
			{ID: "2", Amount: 200.0, Type: Debit, Time: startDate},
			{ID: "3", Amount: 250.0, Type: Debit, Time: startDate},
		}

		bankStatementsBCA := []BankStatement{
			{Bank: "BCA", Amount: 100.0, Time: startDate},
			{Bank: "BCA", Amount: 300.0, Time: startDate},
		}
		bankStatementsBRI := []BankStatement{
			{Bank: "BRI", Amount: 200.0, Time: startDate},
			{Bank: "BRI", Amount: 400.0, Time: startDate},
		}

		suite.mockTransactionStorage.EXPECT().GetTransactions(transactionPath, startDate, endDate).Return(transactions, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bca.xlsx", startDate, endDate).Return(bankStatementsBCA, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bri.xlsx", startDate, endDate).Return(bankStatementsBRI, nil)

		expectedSummary := Summary{
			TotalAmountBankStatements: 1000.0,
			TotalAmountTransactions:   550.0,
			TotalMatched:              2,
			TotalUnmatched:            3,
			TotalProcessed:            5,
		}
		suite.mockSummaryRepoStorage.EXPECT().StoreSummary(expectedSummary).Return(nil)
		suite.mockTransactionStorage.EXPECT().StoreTransactions([]Transaction{{ID: "3", Amount: 250.0, Type: Debit, Time: startDate}}).Return(nil)

		suite.mockBankStatementRepoStorage.EXPECT().StoreBankStatements([]BankStatement{{Bank: "BCA", Amount: 300.0, Time: startDate}}, "BCA").Return(nil)
		suite.mockBankStatementRepoStorage.EXPECT().StoreBankStatements([]BankStatement{{Bank: "BRI", Amount: 400.0, Time: startDate}}, "BRI").Return(nil)

		err := suite.reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
		g.Expect(err).Should(BeNil())
	})

	t.Run("should return error when GetTransactions fails", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getReconExecutorSuite(ctrl)

		suite.mockTransactionStorage.EXPECT().GetTransactions(transactionPath, startDate, endDate).Return(nil, fmt.Errorf("get transactions error"))

		err := suite.reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
		g.Expect(err).Should(Not(BeNil()))
	})

	t.Run("should return error when GetBankStatements fails", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getReconExecutorSuite(ctrl)

		suite.mockTransactionStorage.EXPECT().GetTransactions(transactionPath, startDate, endDate).Return([]Transaction{}, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bca.xlsx", startDate, endDate).Return(nil, fmt.Errorf("get bank statements error"))

		err := suite.reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
		g.Expect(err).Should(Not(BeNil()))
	})

	t.Run("should return error when StoreSummary fails", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getReconExecutorSuite(ctrl)

		transactions := []Transaction{
			{ID: "1", Amount: 100.0, Type: Credit, Time: startDate},
		}

		bankStatementsBCA := []BankStatement{
			{Bank: "BCA", Amount: 100.0, Time: startDate},
		}

		suite.mockTransactionStorage.EXPECT().GetTransactions(transactionPath, startDate, endDate).Return(transactions, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bca.xlsx", startDate, endDate).Return(bankStatementsBCA, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bri.xlsx", startDate, endDate).Return([]BankStatement{}, nil)

		expectedSummary := Summary{
			TotalProcessed:            1,
			TotalAmountBankStatements: 100.0,
			TotalAmountTransactions:   100.0,
			TotalMatched:              1,
			TotalUnmatched:            0,
		}

		suite.mockSummaryRepoStorage.EXPECT().StoreSummary(gomock.Eq(expectedSummary)).Return(fmt.Errorf("store summary error"))

		err := suite.reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
		g.Expect(err).Should(Not(BeNil()))
	})

	t.Run("should return error when StoreTransactions fails", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getReconExecutorSuite(ctrl)

		transactions := []Transaction{
			{ID: "1", Amount: 100.0, Type: Credit, Time: startDate},
		}

		bankStatementsBCA := []BankStatement{
			{Bank: "BCA", Amount: 100.0, Time: startDate},
		}

		suite.mockTransactionStorage.EXPECT().GetTransactions(transactionPath, startDate, endDate).Return(transactions, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bca.xlsx", startDate, endDate).Return(bankStatementsBCA, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bri.xlsx", startDate, endDate).Return([]BankStatement{}, nil)

		expectedSummary := Summary{
			TotalProcessed:            1,
			TotalAmountBankStatements: 100.0,
			TotalAmountTransactions:   100.0,
			TotalMatched:              1,
			TotalUnmatched:            0,
		}
		suite.mockSummaryRepoStorage.EXPECT().StoreSummary(gomock.Eq(expectedSummary)).Return(nil)
		suite.mockTransactionStorage.EXPECT().StoreTransactions(gomock.Eq([]Transaction{})).Return(fmt.Errorf("store transactions error"))

		err := suite.reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
		g.Expect(err).Should(Not(BeNil()))
	})

	t.Run("should return error when StoreBankStatements fails", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		suite := getReconExecutorSuite(ctrl)

		transactions := []Transaction{
			{ID: "1", Amount: 100.0, Type: Credit, Time: startDate},
		}

		bankStatementsBCA := []BankStatement{
			{Bank: "BCA", Amount: 100.0, Time: startDate},
			{Bank: "BCA", Amount: 100.0, Time: startDate},
		}

		suite.mockTransactionStorage.EXPECT().GetTransactions(transactionPath, startDate, endDate).Return(transactions, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bca.xlsx", startDate, endDate).Return(bankStatementsBCA, nil)
		suite.mockBankStatementRepoStorage.EXPECT().GetBankStatements("bri.xlsx", startDate, endDate).Return([]BankStatement{}, nil)

		expectedSummary := Summary{
			TotalProcessed:            2,
			TotalAmountBankStatements: 200.0,
			TotalAmountTransactions:   100.0,
			TotalMatched:              1,
			TotalUnmatched:            1,
		}
		suite.mockSummaryRepoStorage.EXPECT().StoreSummary(gomock.Eq(expectedSummary)).Return(nil)
		suite.mockTransactionStorage.EXPECT().StoreTransactions(gomock.Eq([]Transaction{})).Return(nil)
		suite.mockBankStatementRepoStorage.EXPECT().StoreBankStatements(gomock.Eq([]BankStatement{{Bank: "BCA", Amount: 100.0, Time: startDate}}), "BCA").Return(fmt.Errorf("store bank statements error"))

		err := suite.reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
		g.Expect(err).Should(Not(BeNil()))
	})
}
