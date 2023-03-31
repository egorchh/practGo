package main
import (
	"strconv"
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type test struct{
    id int
    text string
    description string
    keywords string
}
func main() {
    db, err := sql.Open("mysql", "root:@/test")
     
    if err != nil {
        panic(err)
    } 
    defer db.Close()
    rows, err := db.Query("select * from test.myarttable")
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
    for _, p := range tests{
        fmt.Println(p.id, p.text, p.description, p.keywords)
    }
	

// server
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        http.ServeFile(w, r, "user.html")
    })

    http.HandleFunc("/postform", func(w http.ResponseWriter, r *http.Request){
     
        name := r.FormValue("username")
        age := r.FormValue("userage")
         
        fmt.Fprintf(w, "Имя: %s Возраст: %s", name, age)
		
		
		for _, p := range tests{
			fmt.Fprintf(w, "<br />"+strconv.Itoa(p.id)+"<br />"+p.text+"<br />"+p.description+"<br />"+p.keywords)
		}
		
    })
    fmt.Println("Server is listening...")
    http.ListenAndServe(":8181", nil)	
}