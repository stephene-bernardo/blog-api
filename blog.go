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
	"time"
)

func main() {
	hostname := "localhost"
	port := "8080"
	postgresHost := "localhost"
	postgresPort := "5432"
	postgresDbName := "postgres"
	postgresPassword := "abc123"
	postgresUser := "postgres"
	if env := os.Getenv("BASE_URL"); env != "" {
		hostname = env
	}
	if env := os.Getenv("POSTGRES_HOST"); env != "" {
		postgresHost = env
	}
	if env := os.Getenv("POSTGRES_PORT"); env != "" {
		postgresPort = env
	}
	if env := os.Getenv("POSTGRES_DB_NAME"); env != "" {
		postgresDbName = env
	}
	if env := os.Getenv("POSTGRES_USER"); env != "" {
		postgresUser = env
	}
	if env := os.Getenv("POSTGRES_PASSWORD"); env != "" {
		postgresPassword = env
	}
	db := connectToDb(postgresUser, postgresDbName, postgresPassword, postgresPort, postgresHost)
	var tableName = "article"
	createDatabaseTable(tableName, db)
	articleService := dao.ArticleDao{Db: db, Table: tableName}
	articleHandler := rest.ArticleHandler{ArticleDao: &articleService}
	articleRouter := generateArticleRouter(articleHandler)
	log.Println(fmt.Sprintf("listening in port %s...", port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", hostname, port), articleRouter)
	if err != nil {
		log.Fatal(err)
	}
}

func connectToDb(postgresUser string, postgresDbName string,
	postgresPassword string,
	postgresPort string,
	postgresHost string) *sql.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable host=%s", postgresUser,
		postgresDbName,
		postgresPassword,
		postgresPort,
		postgresHost)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return db
}
func waitDbConnection(postgresUser string,
	postgresDbName string,
	postgresPassword string,
	postgresPort string,
	postgresHost string) *sql.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable host=%s", postgresUser,
		postgresDbName,
		postgresPassword,
		postgresPort,
		postgresHost)

	var db *sql.DB
	for {
		log.Println("waiting for database connection...")
		db, _ = sql.Open("postgres", connStr)
		_, err := db.Query(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'",postgresDbName))
		log.Println(err)
		if err != nil {
			break
		}
		time.Sleep(time.Second * 5)
	}
	return db
}

func createDatabaseTable(tableName string, db *sql.DB) error {
	var queryStringForCreateTable = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id serial not null primary key, 
		TITLE VARCHAR, content VARCHAR, author VARCHAR)`, tableName)
	_, err := db.Query(queryStringForCreateTable)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func generateArticleRouter(articleHandler rest.ArticleHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/articles", articleHandler.InsertHandler).Methods("POST")
	r.HandleFunc("/articles", articleHandler.GetAllHandler).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}", articleHandler.GetByIdHandler).Methods("GET")
	r.HandleFunc("/articles/{id:[0-9]+}", articleHandler.RemoveHandler).Methods("DELETE")
	return r
}
