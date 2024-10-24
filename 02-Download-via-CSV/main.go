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
)

const baseUri = "https://www.jaredhettinger.io/lit/txt/"
var errors []ErrorDatum

type ErrorDatum struct {
	CaughtError any `json:"caughtError"`
	Message string `json:"message"`
	RowData RowData `json:"rowData"`
}

type RowData struct {
	WorkTitle string `json:"workTitle"`
	AuthorLastName string `json:"authorLastName"`
	AuthorFirstName string `json:"authorFirstName"`
	PublicationYear int `json:"publicationYear"`
}

type ResponseError struct {
	Status string `json:"status"`
	StatusCode int `json:"statusCode"`
	Header http.Header `json:"header"`
}

func addToErrors(err ErrorDatum) {
	errors = append(errors, err)
}

func printIfError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func download(uri string, saveAs string, context RowData) bool {
	res, err := http.Get(uri)
	httpError := handleHttpError(res, err, context)
	defer res.Body.Close()

	if httpError {
		fmt.Println("! Error downloading " + uri + " as " + saveAs)
		return false
	}

	out, err := os.Create("out/" + saveAs +".txt")
	printIfError(err)
  defer out.Close()

	_, err = io.Copy(out, res.Body)
	printIfError(err)
	return true
}

func handleHttpError(res *http.Response, err error, rowData RowData) bool {
	if err != nil {
		newError := ErrorDatum{
			CaughtError: err,
			Message: "Error occurred during the HTML request",
			RowData: rowData,
		}
		addToErrors(newError)
		return true
	} else if res.StatusCode != 200 {
		newCaughtError := ResponseError{
			Status: res.Status,
			StatusCode: res.StatusCode,
			Header: res.Header,
		}
		newError := ErrorDatum{
			CaughtError: newCaughtError,
			Message: "HTTP response status code was not 200",
			RowData: rowData,
		}
		addToErrors(newError)
		return true
	}
	return false
}

func saveErrors() {
	// Create the errors.json file
	file, err := os.Create("errors.json")
	printIfError(err)
	defer file.Close()

	// Convert the errors array to JSON
	jsonData, err := json.Marshal(errors)
	printIfError(err)

	// Save the errors array to errors.json
	file.Write(jsonData)
	printIfError(err)
}

func main() {
	fmt.Print("Running Download-via-CSV...\n\n")

	data, err := os.ReadFile("in.csv")
	printIfError(err)

	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll();
	printIfError(err)

	var rows = []RowData{};
	for i := 1; i < len(records); i++ {
		parsedPublicationYear, err := strconv.Atoi(records[i][3])
		printIfError(err)
		newRow := RowData{
			WorkTitle: records[i][0],
			AuthorLastName: records[i][1],
			AuthorFirstName: records[i][2],
			PublicationYear: parsedPublicationYear,
		};
		rows = append(rows, newRow)
	}

	for i := 0; i < len(rows); i++ {
		downloadSlug := rows[i].AuthorLastName + " - " + rows[i].WorkTitle + ".txt"
		downloadUri := baseUri + downloadSlug;
		fmt.Println("Downloading file from " + downloadUri)
		download(downloadUri, rows[i].AuthorLastName, rows[i])
	}

	if len(errors) > 0 {
		errorCount := strconv.Itoa(len(errors))
		fmt.Println("\n" + errorCount + " error(s) found, saving to errors.json")
		saveErrors()
	}

	fmt.Print("\nClosing\n")
}