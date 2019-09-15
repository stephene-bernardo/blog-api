package dao

import "blog-api/src"

type ArticleDaoInterface interface {
	FindAll() ([]src.ArticleObject, error)
	FindById(id int) (src.ArticleObject, error)
	Insert(title, content, author string) (int, error)
	Delete(id int) (string, error)
}
