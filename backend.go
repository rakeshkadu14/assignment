package main

import (
	"encoding/csv"
	"io"
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
			//fmt.Println("endday=>", endDay)
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

func analyzeMonth(reservations []Reservation, month string) (float64, int, error) {
	//fmt.Println(month)
	monthStart, err := time.Parse("2006-01", month)

	if err != nil {
		return 0, 0, err
	}
	monthEnd := monthStart.AddDate(0, 1, -1)
	//fmt.Println("monthEnd", monthEnd)
	daysInMonth := monthEnd.Day()
	//fmt.Println("daysInMonth", daysInMonth)
	//fmt.Println("-----------------")
	totalRevenue := 0.0
	unreservedCapacity := 0

	for _, res := range reservations {
		// Check overlap
		//fmt.Println("*****************")
		//fmt.Println("res=>", res)
		//fmt.Println("*****************")
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
		//fmt.Println("start", start)
		//fmt.Println("end", end)
		//fmt.Println("daysReserved", daysReserved)
		//fmt.Println("res.MonthlyPrice", res.MonthlyPrice)
		proratedRevenue := res.MonthlyPrice * float64(daysReserved) / float64(daysInMonth)
		//fmt.Println("proratedRevenue", proratedRevenue)
		totalRevenue += proratedRevenue
		//fmt.Println("totalRevenue", totalRevenue)
	}

	return totalRevenue, unreservedCapacity, nil
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
