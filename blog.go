package main

import (
	"blog-api/src/dao"
	"blog-api/src/rest"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main(){
	fmt.Println("Hello blog API")
	connStr:="user=postgres dbname=blog password=abc123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	var databaseTable = "article"
	//var queryStringForDropTable = fmt.Sprintf("DROP TABLE IF EXISTs %s", databaseTable)
	var queryStringForCreateTable = fmt.Sprintf(`CREATE TABLE %s (id serial not null primary key, TITLE VARCHAR,
    content VARCHAR, author VARCHAR)`, databaseTable)
	//db.Query(queryStringForDropTable)
	db.Query(queryStringForCreateTable)
	if err != nil {
		log.Fatal(err)
	}

	articleService := dao.ArticleDao{db, databaseTable}
	articleHandler := rest.ArticleHandler{&articleService}

	r := mux.NewRouter()
	r.HandleFunc("/articles",  articleHandler.InsertHandler).Methods("POST")
	r.HandleFunc("/articles",  articleHandler.GetAllHandler).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}", articleHandler.GetByIdHandler).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}",  articleHandler.RemoveHandler).Methods("DELETE")

	fmt.Println("listening ...")
	http.ListenAndServe("localhost:8080", r)
}
