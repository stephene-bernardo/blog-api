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

func main() {
	baseurl := "localhost"
	port := "8080"
	postgresHost := "localhost"
	postgresPort := "5432"
	postgresDbName := "postgres"
	postgresPassword := "abc123"
	postgresUser := "postgres"
	if os.Getenv("BASE_URL") != "" {
		baseurl = os.Getenv("BASE_URL")
	}
	if os.Getenv("POSTGRES_HOST") != "" {
		postgresHost = os.Getenv("POSTGRES_HOST")
	}
	if os.Getenv("POSTGRES_PORT") != "" {
		postgresPort = os.Getenv("POSTGRES_PORT")
	}
	if os.Getenv("POSTGRES_DB_NAME") != "" {
		postgresDbName = os.Getenv("POSTGRES_DB_NAME")
	}
	if os.Getenv("POSTGRES_USER") != "" {
		postgresUser = os.Getenv("POSTGRES_USER")
	}
	if os.Getenv("POSTGRES_PASSWORD") != "" {
		postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	}
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable host=%s", postgresUser,
		postgresDbName,
		postgresPassword,
		postgresPort,
		postgresHost)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	var databaseTable = "article"
	var queryStringForCreateTable = fmt.Sprintf(`CREATE TABLE %s (id serial not null primary key, TITLE VARCHAR,
    content VARCHAR, author VARCHAR)`, databaseTable)
	_, err = db.Query(queryStringForCreateTable)
	if err != nil {
		log.Fatal(err)
	}
	articleService := dao.ArticleDao{Db: db, Table: databaseTable}
	articleHandler := rest.ArticleHandler{ArticleDao: &articleService}

	r := mux.NewRouter()
	r.HandleFunc("/articles", articleHandler.InsertHandler).Methods("POST")
	r.HandleFunc("/articles", articleHandler.GetAllHandler).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}", articleHandler.GetByIdHandler).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}", articleHandler.RemoveHandler).Methods("DELETE")

	fmt.Println(fmt.Sprintf("listening in port %s...", port))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", baseurl, port), r)
	if err != nil {
		log.Fatal(err)
	}
}
