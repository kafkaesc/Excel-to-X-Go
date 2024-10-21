package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InRow struct {
	workTitle string
	authorLastName string
	authorFirstName string
	publicationYear int
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	fmt.Print("Running Download-via-CSV...\n\n")

	data, err := os.ReadFile("in.csv")
	checkError(err)

	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll();
	checkError(err)

	var rows = []InRow{};
	for i := 1; i < len(records); i++ {
		parsedPublicationYear, err := strconv.Atoi(records[i][3])
		checkError(err)
		newRow := InRow{
			workTitle: records[i][0],
			authorFirstName: records[i][1],
			authorLastName: records[i][2],
			publicationYear: parsedPublicationYear,
		};
		rows = append(rows, newRow)
	}

	for i := 0; i < len(rows); i++ {
		fmt.Println("Row ", i+1, ": ", rows[i])
	}

	fmt.Print("Closing\n")
}