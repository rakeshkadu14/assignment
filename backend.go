package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Reservation struct {
	Capacity     int
	MonthlyPrice float64
	StartDay     time.Time
	EndDay       *time.Time
}

type Report struct {
	Month             string  `json:"month"`
	ExpectedRevenue   float64 `json:"expectedRevenue"`
	UnreservedOffices int     `json:"unreservedOffices"`
}

func parseCSV(filepath string) ([]Reservation, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var reservations []Reservation

	reader := csv.NewReader(file)
	_, _ = reader.Read() // skip header (1st line)

	for {
		record, err := reader.Read()
		//fmt.Println(record)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		capacity, _ := strconv.Atoi(strings.TrimSpace(record[0]))
		price, _ := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		startDay, _ := time.Parse("2006-01-02", strings.TrimSpace(record[2]))
		var endDay *time.Time
		if strings.TrimSpace(record[3]) != "" {
			e, _ := time.Parse("2006-01-02", record[3])
			endDay = &e
		}

		reservations = append(reservations, Reservation{
			Capacity:     capacity,
			MonthlyPrice: price,
			StartDay:     startDay,
			EndDay:       endDay,
		})
		//fmt.Println(reservations)
	}
	return reservations, nil
}

func analyzeMonth(w http.ResponseWriter, r *http.Request) {

	filepath := os.Args[1]
	reservations, err := parseCSV(filepath)
	if err != nil {
		fmt.Println("Error parsing CSV:", err)
		os.Exit(0)
	}

	month := r.URL.Query().Get("month")
	//fmt.Println(month)
	monthStart, err := time.Parse("2006-01", month)

	if err != nil {
		os.Exit(0)
	}
	monthEnd := monthStart.AddDate(0, 1, -1)
	daysInMonth := monthEnd.Day()
	totalRevenue := 0.0
	unreservedCapacity := 0

	for _, res := range reservations {
		resEnd := res.EndDay
		if resEnd == nil {
			tmp := time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
			resEnd = &tmp
		}
		if res.StartDay.After(monthEnd) || resEnd.Before(monthStart) {
			unreservedCapacity += res.Capacity
			continue
		}

		// Calculate overlap days
		start := maxDate(res.StartDay, monthStart)
		end := minDate(*resEnd, monthEnd)
		daysReserved := int(end.Sub(start).Hours()/24) + 1
		proratedRevenue := res.MonthlyPrice * float64(daysReserved) / float64(daysInMonth)
		totalRevenue += proratedRevenue
	}

	report := Report{
		Month:             month,
		ExpectedRevenue:   totalRevenue,       // Dummy data
		UnreservedOffices: unreservedCapacity, // Dummy data
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func minDate(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
