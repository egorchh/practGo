package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sql.Open("mysql", "root:password@/space_observatory")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		viewSelect(w, db)
	})

	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request) {

		objType := r.FormValue("type")
		accuracy := r.FormValue("accuracy")
		quantity := r.FormValue("quantity")
		time := r.FormValue("time")
		date := r.FormValue("date")
		notes := r.FormValue("notes")

		sQuery := "INSERT INTO objects (type, accuracy, quantity, time, date, notes) VALUES (?, ?, ?, ?, ?, ?)"

		_, err := db.Exec(sQuery, objType, accuracy, quantity, time, date, notes)

		if err != nil {
			panic(err)
		}

		viewSelect(w, db)
	})

	fmt.Println("Server is listening on http://localhost:8181/")
	http.ListenAndServe(":8181", nil)
}

func mapObjectType(objType string) string {
	var objTypeMap = map[string]string{
		"planet":    "Планета",
		"star":      "Звезда",
		"satellite": "Спутник",
		"asteroid":  "Астероид",
		"comet":     "Комета",
		"meteorite": "Метеорит",
	}

	return objTypeMap[objType]
}

func viewSelect(w http.ResponseWriter, db *sql.DB) {
	file, err := os.Open("select.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "@tr" {
			viewTableData(w, db)
		} else if scanner.Text() == "@ver" {
			viewSelectVerQuery(w, db, "SELECT VERSION() AS ver")
		} else {
			fmt.Fprintf(w, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func viewTableData(w http.ResponseWriter, db *sql.DB) {
	type object struct {
		id       int
		objType  string
		accuracy string
		quantity string
		time     string
		date     string
		notes    string
	}

	sQuery := "SELECT id, type, accuracy, quantity, time, date, notes FROM objects"
	rows, err := db.Query(sQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fmt.Fprintf(w, "<table>")
	// Write table header
	fmt.Fprintf(w, "<tr><th>ID</th><th>Type</th><th>Accuracy</th><th>Quantity</th><th>Time</th><th>Date</th><th>Notes</th></tr>")

	for rows.Next() {
		var p object
		err := rows.Scan(&p.id, &p.objType, &p.accuracy, &p.quantity, &p.time, &p.date, &p.notes)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", p.id, mapObjectType(p.objType), p.accuracy, p.quantity, p.time, p.date, p.notes)
	}
	fmt.Fprintf(w, "</table>")
}

func viewSelectVerQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
	type sVer struct {
		ver string
	}
	rows, err := db.Query(sSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		p := sVer{}
		err := rows.Scan(&p.ver)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(w, p.ver)
	}
}
