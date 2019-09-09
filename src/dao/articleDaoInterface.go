package dao

type ArticleDaoInterface interface {
  FindAll()([]ArticleObject, error)
  FindById(id int)(ArticleObject, error)
  Insert(title , content, author string) (error)
  Delete(id int) (error)
}
