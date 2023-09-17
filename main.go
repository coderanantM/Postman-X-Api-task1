package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var branches = map[string]string{
	"AA": "ECE",
	"AB": "Manu",
	"A1": "Chemical",
	"A2": "Civil",
	"A3": "EEE",
	"A4": "Mech",
	"A5": "Pharma",
	"A7": "CSE",
	"A8": "ENI",
	"B1": "MSc BIO",
	"B2": "MSc Chem",
	"B3": "MSc Eco",
	"B4": "MSc Mathematics",
	"B5": "Msc Physics",
}

type Student struct {
	Name             string `json:"name"`
	BITSID           string `json:"BITS ID"`
	BITSEmailAddress string `json:"BITS email address"`
	Branch           string `json:"Branch"`
}

func main() {
	r := gin.Default()

	r.GET("/students", func(c *gin.Context) {
		data, err := readExcelFile("BITS_Students_Info.xlsx")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		students := []Student{}
		for _, row := range data {
			idNo := row[1][8:12]
			year := row[1][0:4]
			branch := branches[row[1][4:6]]
			email := fmt.Sprintf("f%s%s@pilani.bits-pilani.ac.in", year, idNo)

			if row[1][6:8] != "PS" {
				branch += " + " + branches[row[1][6:8]]
			}

			student := Student{
				Name:             row[0],
				BITSID:           row[1],
				BITSEmailAddress: email,
				Branch:           branch,
			}
			students = append(students, student)
		}

		c.JSON(http.StatusOK, students)
	})

	r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

func readExcelFile(filename string) ([][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheetName := f.GetSheetName(1)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
