package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run client.go <csv-filepath> <YYYY-MM>")
		return
	}

	filepath := os.Args[1]
	month := os.Args[2]

	reservations, err := parseCSV(filepath)
	if err != nil {
		fmt.Println("Error parsing CSV:", err)
		return
	}

	revenue, unreservedCap, err := analyzeMonth(reservations, month)
	if err != nil {
		fmt.Println("Error analyzing month:", err)
		return
	}

	fmt.Printf("* %s: expected revenue: $%.0f, expected total capacity of the unreserved offices: %d\n", month, revenue, unreservedCap)
}
