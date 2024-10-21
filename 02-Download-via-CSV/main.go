package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const baseUri = "https://www.jaredhettinger.io/lit/txt/";

/* type DownloadContext struct {
	row int
	rowData InRow
} */

type InRow struct {
	workTitle string
	authorLastName string
	authorFirstName string
	publicationYear int
}

func logIfError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func download(uri string, saveAs string/*, context DownloadContext*/) {
	res, err := http.Get(uri)
	logIfError(err)
	defer res.Body.Close()

	out, err := os.Create("out/" + saveAs +".txt")
	logIfError(err)
  defer out.Close()

	_, err = io.Copy(out, res.Body)
	logIfError(err)
}

func main() {
	fmt.Print("Running Download-via-CSV...\n\n")

	data, err := os.ReadFile("in.csv")
	logIfError(err)

	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll();
	logIfError(err)

	var rows = []InRow{};
	for i := 1; i < len(records); i++ {
		parsedPublicationYear, err := strconv.Atoi(records[i][3])
		logIfError(err)
		newRow := InRow{
			workTitle: records[i][0],
			authorLastName: records[i][1],
			authorFirstName: records[i][2],
			publicationYear: parsedPublicationYear,
		};
		rows = append(rows, newRow)
	}

	for i := 0; i < len(rows); i++ {
		downloadSlug := rows[i].authorLastName + " - " + rows[i].workTitle + ".txt"
		downloadUri := baseUri + downloadSlug;
		fmt.Println(downloadUri)
		download(downloadUri, rows[i].authorLastName)
	}

	fmt.Print("\nClosing\n")
}