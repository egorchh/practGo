package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	db, err := sql.Open("mysql", "root:password@/bank")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		viewSelect(w, db)
	})

	// сохранение отправленных значений через поля формы.
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {

		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		patronymic := r.FormValue("patronymic")
		passport := r.FormValue("passport")
		tin := r.FormValue("tin")
		snils := r.FormValue("snils")
		driverLicense := r.FormValue("driver_license")
		additionalDocuments := r.FormValue("additional_documents")
		notes := r.FormValue("notes")
		borrowerId := r.FormValue("borrower_id")

		sQuery := ""
		var rows *sql.Rows
		var err error

		if borrowerId == "" {
			sQuery = "INSERT INTO individuals (first_name, last_name, patronymic, passport, tin, snils, driver_license, additional_documents, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
			rows, err = db.Query(sQuery, firstName, lastName, patronymic, passport, tin, snils, driverLicense, additionalDocuments, notes)
		} else {
			sQuery = "INSERT INTO individuals (first_name, last_name, patronymic, passport, tin, snils, driver_license, additional_documents, notes, borrower_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			rows, err = db.Query(sQuery, firstName, lastName, patronymic, passport, tin, snils, driverLicense, additionalDocuments, notes, borrowerId)
		}

		fmt.Println(sQuery)

		if err != nil {
			panic(err)
		}
		defer rows.Close()

		viewSelect(w, db)
	})

	fmt.Println("Server is listening on http://localhost:8181/")
	http.ListenAndServe(":8181", nil)
}

func viewHeadQuery(w http.ResponseWriter, db *sql.DB, sShow string) {
	type sHead struct {
		clnme string
	}
	rows, err := db.Query(sShow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, "<tr>")
	for rows.Next() {
		var p sHead
		err := rows.Scan(&p.clnme)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "<td>%s</td>", p.clnme)
	}
	fmt.Fprintf(w, "</tr>")

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewSelectQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
	type individual struct {
		id                  int
		firstName           sql.NullString
		lastName            sql.NullString
		patronymic          sql.NullString
		passport            sql.NullString
		tin                 sql.NullString
		snils               sql.NullString
		driverLicense       sql.NullString
		additionalDocuments sql.NullString
		notes               sql.NullString
		borrowerId          sql.NullInt64
	}
	individuals := []individual{}

	// получение значений в массив individuals из структуры типа individual.
	rows, err := db.Query(sSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := individual{}
		err := rows.Scan(&p.id, &p.firstName, &p.lastName, &p.patronymic, &p.passport, &p.tin, &p.snils, &p.driverLicense, &p.additionalDocuments, &p.notes, &p.borrowerId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		individuals = append(individuals, p)
	}

	// перебор массива из БД.
	for _, p := range individuals {
		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			p.id,
			nullStringToEmpty(p.firstName),
			nullStringToEmpty(p.lastName),
			nullStringToEmpty(p.patronymic),
			nullStringToEmpty(p.passport),
			nullStringToEmpty(p.tin),
			nullStringToEmpty(p.snils),
			nullStringToEmpty(p.driverLicense),
			nullStringToEmpty(p.additionalDocuments),
			nullStringToEmpty(p.notes),
			nullInt64ToEmpty(p.borrowerId))
	}
}

func nullStringToEmpty(ns sql.NullString) string {
	if ns.Valid && ns.String != "" {
		return ns.String
	}
	return "EMPTY"
}

func nullInt64ToEmpty(ni sql.NullInt64) string {
	if ni.Valid {
		return strconv.FormatInt(ni.Int64, 10)
	}
	return "EMPTY"
}

func viewSelectVerQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
	type sVer struct {
		ver string
	}
	rows, err := db.Query(sSelect)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p sVer
		err := rows.Scan(&p.ver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, p.ver)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewSelect(w http.ResponseWriter, db *sql.DB) {
	// чтение шаблона.
	file, err := os.Open("select.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//	кодовая фраза для вставки значений из БД.
		if scanner.Text() != "@tr" && scanner.Text() != "@ver" {
			fmt.Fprintf(w, scanner.Text())
		}
		if scanner.Text() == "@tr" {
			viewHeadQuery(w, db, "select COLUMN_NAME AS clnme from information_schema.COLUMNS where TABLE_NAME='individuals' ORDER BY ORDINAL_POSITION")
			viewSelectQuery(w, db, "SELECT * FROM individuals ORDER BY id ASC")
		}
		if scanner.Text() == "@ver" {
			viewSelectVerQuery(w, db, "SELECT VERSION() AS ver")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
