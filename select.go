package main
import (
	"strconv"
    "bufio"
    "fmt"
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
	

    fmt.Println("Server is listening...")
    http.ListenAndServe(":8181", nil)	
}


// главная функция для показа таблицы в браузере, которая показывается при любом запросе.
func viewSelect(w http.ResponseWriter, db *sql.DB) {
	type test struct {
		id int
		text string
		description string
		keywords string
	}
	
	// получение значений в массив tests из струкрур типа test.
    rows, err := db.Query("SELECT * FROM myarttable WHERE id>14 ORDER BY id DESC")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    tests := []test{}
     
    for rows.Next(){
        p := test{}
        err := rows.Scan(&p.id, &p.text, &p.description, &p.keywords)
        if err != nil{
            fmt.Println(err)
            continue
        }
        tests = append(tests, p)
    }

	// чтение файла и выгрузка значений из массива.
	file, err := os.Open("select.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//	кодовая фраза для вставки значений из БД.
		if scanner.Text() != "@tr" {
			fmt.Fprintf(w, scanner.Text())
		} else {
			// перебор массива из БД.
			for _, p := range tests {
				fmt.Fprintf(w, "<tr><td>"+strconv.Itoa(p.id)+"</td><td>"+p.text+"</td><td>"+p.description+"</td><td>"+p.keywords+"</td></tr>")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

