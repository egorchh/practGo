package main
import (
	"strconv"
    "bufio"
    "fmt"
	//"reflect"
    "log"
    "os"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:@/test")
     
    if err != nil {
        panic(err)
    } 
    defer db.Close()
	
	// открытие из браузера корневого каталога.
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	
		viewSelect(w, db)
    })

	// сохранение отправленных значений через поля формы.
	http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request){
     
        val1 := r.FormValue("col1")
        val2 := r.FormValue("col2")
        val3 := r.FormValue("col3")
		sQuery := "INSERT INTO myarttable (text, description, keywords) VALUES ('"+val1+"', '"+val2+"', '"+val3+"')"
 
		fmt.Println(sQuery)
 
		rows, err := db.Query(sQuery)
		if err != nil {
			panic(err)
		}		
		defer rows.Close()
		
		viewSelect(w, db)
    })
	

    fmt.Println("Server is listening on http://localhost:8181/")
    http.ListenAndServe(":8181", nil)	
}


// отправка в браузер заголовка таблицы.
func viewHeadQuery(w http.ResponseWriter, db *sql.DB, sShow string) {
	type sHead struct {
		clnme string
	}
    rows, err := db.Query(sShow)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	fmt.Fprintf(w, "<tr>")
	i := 0
     for i < 4 {
		rows.Next()
        p := sHead{}
        err := rows.Scan(&p.clnme)
        if err != nil{
            fmt.Println(err)
            continue
        }
		fmt.Fprintf(w, "<td>"+p.clnme+"</td>")
		i++
    }
	fmt.Fprintf(w, "</tr>")
}

// отправка в браузер строк из таблицы.
func viewSelectQuery(w http.ResponseWriter, db *sql.DB, sSelect string) {
	type test struct {
		id int
		text string
		description string
		keywords string
	}
	tests := []test{}
	//fmt.Println(reflect.TypeOf(tests))

	// получение значений в массив tests из струкрур типа test.
    rows, err := db.Query(sSelect)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
     
    for rows.Next(){
        p := test{}
        err := rows.Scan(&p.id, &p.text, &p.description, &p.keywords)
        if err != nil{
            fmt.Println(err)
            continue
        }
        tests = append(tests, p)
    }
	
	// перебор массива из БД.
	for _, p := range tests {
		fmt.Fprintf(w, "<tr><td>"+strconv.Itoa(p.id)+"</td><td>"+p.text+"</td><td>"+p.description+"</td><td>"+p.keywords+"</td></tr>")
	}
}
	
// отправка в браузер версии базы данных.
func viewSelectVerQuery (w http.ResponseWriter, db *sql.DB, sSelect string) {
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
        if err != nil{
            fmt.Println(err)
            continue
        }
		fmt.Fprintf(w, p.ver)
    }
}

// главная функция для показа таблицы в браузере, которая показывается при любом запросе.
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
			viewHeadQuery(w, db, "select COLUMN_NAME AS clnme from information_schema.COLUMNS where TABLE_NAME='myarttable'")
			viewSelectQuery(w, db, "SELECT * FROM myarttable WHERE id>14 ORDER BY id DESC")
		}
		if scanner.Text() == "@ver" {
			viewSelectVerQuery(w, db, "SELECT VERSION() AS ver")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

