package dao

import (
  "log"
  "database/sql"
  "errors"
  "strconv"
  "fmt"
)

type ArticleDao struct {
  Db *sql.DB
}




func (articleDao *ArticleDao) FindAll()([]ArticleObject, error){
  findAllQuery := "SELECT * FROM article"
  rows, err :=  articleDao.Db.Query(findAllQuery)
  if err != nil {
    log.Fatal(err)
  }
  articles := make([]ArticleObject, 0)
  defer rows.Close()
  for rows.Next(){
    var id int
    var title, content, author string
    rows.Scan(&id, &title, &content, &author)
    articles = append(articles, ArticleObject{Id: id, Title: title, Content: content, Author: author})
  }
  return articles, err
}

func (articleDao *ArticleDao) FindById(id int)(ArticleObject, error){
  findByIdQuery := fmt.Sprintf("SELECT * FROM article WHERE id = %d", id)
  rows, err :=  articleDao.Db.Query(findByIdQuery)
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

    var title, content, author string
    if rows.Next(){
      rows.Scan(&id, &title, &content, &author)
      return ArticleObject{Id: id, Title: title, Content: content, Author: author}, err
    }

   return ArticleObject{}, errors.New("Unable to find id:" + strconv.Itoa(id))
}

func (articleDao *ArticleDao) Insert(title , content, author string) (error){
  insertUserQuery := "INSERT INTO article( title, content, author)VALUES($1, $2, $3)"
  _, err :=  articleDao.Db.Query(insertUserQuery, title, content, author)
  if err != nil {
    log.Fatal(err)
  }
  return err
}

func (articleDao *ArticleDao) Delete(id int) (error){
  deleteUserQuery := fmt.Sprintf("DELETE FROM article WHERE id = %d", id)
  _, err := articleDao.Db.Query(deleteUserQuery)
  if err != nil {
    log.Fatal(err)
  }
  return err
}
