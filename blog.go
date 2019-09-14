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
	"os"
)

func main(){

	baseurl := "localhost"
	port := "8080"
	postgresHost := "localhost"
	if os.Getenv("BASE_URL") != "" {
		baseurl = os.Getenv("BASE_URL")
	}
	if os.Getenv("POSTGRES_HOST") != "" {
		postgresHost = os.Getenv("POSTGRES_HOST")
	}
	connStr:=fmt.Sprintf("user=postgres dbname=blog password=abc123 port=5432 sslmode=disable host=%s", postgresHost)
	db, err := sql.Open("postgres", connStr)
	var databaseTable = "article"
	var queryStringForCreateTable = fmt.Sprintf(`CREATE TABLE %s (id serial not null primary key, TITLE VARCHAR,
    content VARCHAR, author VARCHAR)`, databaseTable)
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

	fmt.Println(fmt.Sprintf("listening in port %s...", port))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", baseurl, port), r)
	if err !=nil {
		log.Fatal(err)
	}
}
