package dao

import (
	"database/sql"
  "errors"
  "fmt"
	"log"
  "strconv"
)

type ArticleDao struct {
	Db    *sql.DB
	Table string
}

func (articleDao *ArticleDao) FindAll() ([]ArticleObject, error) {
	findAllQuery := fmt.Sprintf("SELECT * FROM %s", articleDao.Table)
	rows, err := articleDao.Db.Query(findAllQuery)
	articles := make([]ArticleObject, 0)
	if err != nil {
		log.Println(err)
		return articles, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int
			var title, content, author string
			rows.Scan(&id, &title, &content, &author)
			articles = append(articles, ArticleObject{Id: id, Title: title, Content: content, Author: author})
		}
	}
	return articles, err
}

func (articleDao *ArticleDao) FindById(id int) (ArticleObject, error) {
	findByIdQuery := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", articleDao.Table, id)
	rows, err := articleDao.Db.Query(findByIdQuery)
	if err != nil {
		log.Println(err)
		return ArticleObject{}, err
	} else {
		var title, content, author string
		defer rows.Close()
		if rows.Next() {
			rows.Scan(&id, &title, &content, &author)
			return ArticleObject{Id: id, Title: title, Content: content, Author: author}, err
		}
	}
	return ArticleObject{}, errors.New("Unable to find id:" + strconv.Itoa(id))
}

func (articleDao *ArticleDao) Insert(title, content, author string) (int, error) {
	var userid int
	queryString := fmt.Sprintf("INSERT INTO %s( title, content, author) VALUES ($1, $2, $3) RETURNING id", articleDao.Table)
	insert, err := articleDao.Db.Query(queryString, title, content, author)
	if err != nil {
		log.Println(err)
	} else {
		insert.Next()
		insert.Scan(&userid)
		log.Println(fmt.Sprintf("Inserted id: %d in the database", userid))
	}
	return userid, err
}

func (articleDao *ArticleDao) Delete(id int) (string, error) {
	deleteUserQuery := fmt.Sprintf("DELETE FROM %s WHERE id = %d RETURNING title", articleDao.Table, id)
	var title string
	rows, err := articleDao.Db.Query(deleteUserQuery)
	if err != nil {
		log.Println(err)
	} else {
		rows.Next()
		rows.Scan(&title)
		log.Println(fmt.Sprintf("Deleted id: %d in the database", id))
	}
	return title, err
}
