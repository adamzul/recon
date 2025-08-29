package main

import (
	"flag"
	"log"
	"recon/recon"
	"time"
)

const reconPath = "data/recon.xlsx"

func main() {
	var transactionPath, bankStatementPaths string
	var startDateStr, endDateStr string
	flag.StringVar(&transactionPath, "transaction-path", "transaction.csv", "transactions CSV file path")
	flag.StringVar(&bankStatementPaths, "bank-statement-paths", "bca.csv,bri.csv", "bank statements CSV file path")
	flag.StringVar(&startDateStr, "start-date", time.Now().Format("2006-01-02"), "bank statements CSV file path")
	flag.StringVar(&endDateStr, "end-date", time.Now().Format("2006-01-02"), "bank statements CSV file path")
	flag.Parse()

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Println(err)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Println(err)
		return
	}

	excelFactory := recon.ExcelFactory{}
	csvReaderFactory := recon.CSVReaderFactory{}

	reconExecutor := recon.NewReconExecutor(
		recon.NewTransactionStorage(reconPath, "Transaction", excelFactory, csvReaderFactory),
		recon.NewBankStatementStorage(reconPath),
		recon.NewSummaryStorage(reconPath, "Summary", excelFactory),
	)

	reconExecutor.Execute(transactionPath, bankStatementPaths, startDate, endDate)
}
