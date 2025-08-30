package main

import (
	"flag"
	"log"
	"recon/recon"
	"strings"
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

	bankStatementPathArray := strings.Split(bankStatementPaths, ",")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Panic(err)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		log.Panic(err)
	}

	excelFactory := recon.ExcelFactory{}
	csvReaderFactory := recon.CSVReaderFactory{}

	reconExecutor := recon.NewReconExecutor(
		recon.NewTransactionStorage(reconPath, "Transaction", excelFactory, csvReaderFactory),
		recon.NewBankStatementStorage(reconPath, excelFactory, csvReaderFactory),
		recon.NewSummaryStorage(reconPath, "Summary", excelFactory),
	)

	err = reconExecutor.Execute(transactionPath, bankStatementPathArray, startDate, endDate)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Recon completed successfully")
}
