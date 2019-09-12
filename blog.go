package main

import (
	"blog-api/src/dao"
	"blog-api/src/rest"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	r := mux.NewRouter()
	r.HandleFunc("/articles",  func(w http.ResponseWriter, r *http.Request){
		articleObject := dao.ArticleObject{}
		json.NewDecoder(r.Body).Decode(&articleObject)
		userid, _ := articleService.Insert(articleObject.Title, articleObject.Content, articleObject.Author)
		postData := rest.ArticlePostData{userid}
		postResponse := rest.ArticlePostResponse{201, "Success", postData}
		json.NewEncoder(w).Encode(postResponse)

	}).Methods("POST")

	r.HandleFunc("/articles",  func(w http.ResponseWriter, r *http.Request){
		objects, _ := articleService.FindAll()
		articleResponse := rest.ArticleGetAllResponse{200, "Success", objects}
		json.NewEncoder(w).Encode(articleResponse)
	}).Methods("GET")

	r.HandleFunc("/articles/{id:[0-9]+}",  func(w http.ResponseWriter, r *http.Request){
		arr := strings.Split(r.RequestURI, "/")
		id, _:= strconv.Atoi(arr[len(arr)-1])
		object, _ := articleService.FindById(id)
		articleResponse := rest.ArticleGetIdResponse{200, "Success", object}
		json.NewEncoder(w).Encode(articleResponse)
	}).Methods("GET")

	r.HandleFunc("/articles/{id:[0-9]+}",  func(w http.ResponseWriter, r *http.Request){
		arr := strings.Split(r.RequestURI, "/")
		id, _:= strconv.Atoi(arr[len(arr)-1])
		title, _:= articleService.Delete(id)
		articleData := rest.ArticleDeleteData{title}
		articleResponse := rest.ArticleDeleteResponse{200, "Success", articleData}
		json.NewEncoder(w).Encode(articleResponse)
	}).Methods("DELETE")

	fmt.Println("listening ...")
	http.ListenAndServe("localhost:8080", r)
}
