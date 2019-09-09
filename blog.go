package main

import (
  "fmt"
  "database/sql"
  _ "github.com/lib/pq"
  "log"
  "blog-api/src/dao"
)

func main(){
  fmt.Println("Hello blog API")
  connStr:="user=postgres dbname=blog password=abc123 sslmode=disable"
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }

  articleService := dao.ArticleDao{Db: db}
  fmt.Println(articleService.FindById(1))
  fmt.Println(articleService.FindAll())
  articleService.Insert( "Shining", "some horror thing", "Stephen King")
  fmt.Println(articleService.FindAll())
 // articleService.Delete(2)
 // fmt.Println(articleService.FindAll())
}
