package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/thedatashed/xlsxreader"
)

const baseUri = "https://www.jaredhettinger.io/lit/txt/"
var errors []ErrorDatum

// Represents the details of a download error
type ErrorDatum struct {
	CaughtError any `json:"caughtError"`
	Message string `json:"message"`
	RowData RowData `json:"rowData"`
}

// Represents a row of data from the CSV file
type RowData struct {
	WorkTitle string `json:"workTitle"`
	AuthorLastName string `json:"authorLastName"`
	AuthorFirstName string `json:"authorFirstName"`
	PublicationYear int `json:"publicationYear"`
}

// Represents the details of an HTTP response
type ResponseError struct {
	Status string `json:"status"`
	StatusCode int `json:"statusCode"`
	Header http.Header `json:"header"`
}

/*
 * Adds the argument to the program's errors array
 * @param err The incoming error details
 */
 func addToErrors(err ErrorDatum) {
	errors = append(errors, err)
}

/*
 * Download the file according to the given parameters. If there is an error 
 * with the HTTP request then pass that along to the HTTP error handler 
 * function instead.
 *
 * @param uri The URI for the file
 * @param saveAs The filename to save the file under
 * @param context The details for the download from the CSV file
 * @return true if the download is successful, false if not
 */
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

/*
 * @param res The response details from an HTTP request
 * @param err The nil/error object from the the same HTTP request
 * @param rowData The details for the download from the CSV file
 * @return true if the HTTP res indicates a failure or non-200 result, otherwise false
 */
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

/*
 * Check for an error and log it if found
 * @param err The possibly nil error object
 */
 func printIfError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

/*
 * Save the errors array into a file
 */
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

	// Open and read the XLSX file
	xl, err := xlsxreader.OpenFile("in.xlsx")
	printIfError(err)
	defer xl.Close()

	// Loop through the XLSX file and create an array of the row data
	var rows = []RowData{}
	for row := range xl.ReadRows(xl.Sheets[0]) {
		parsedPublicationYear, err := strconv.Atoi(row.Cells[3].Value)
		printIfError(err)
		newRow := RowData{
			WorkTitle: row.Cells[0].Value,
			AuthorLastName: row.Cells[1].Value,
			AuthorFirstName: row.Cells[2].Value,
			PublicationYear: parsedPublicationYear,
		}
		rows = append(rows, newRow)
	}

	// fmt.Println("_jhdb: rows:", rows)

	// Loop through the row data and download the files
	for i := 0; i < len(rows); i++ {
		downloadSlug := rows[i].AuthorLastName + " - " + rows[i].WorkTitle + ".txt"
		downloadUri := baseUri + downloadSlug;
		fmt.Println("Downloading file from " + downloadUri)
		download(downloadUri, rows[i].AuthorLastName, rows[i])
	}

	// If any errors occurred, save them
	if len(errors) > 0 {
		errorCount := strconv.Itoa(len(errors))
		fmt.Println("\n" + errorCount + " error(s) found, saving to errors.json")
		saveErrors()
	}

	fmt.Print("\nClosing\n")
}