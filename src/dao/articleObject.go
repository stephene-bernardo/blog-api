package dao

type ArticleObject struct {
  Id int `json:"id"`
  Title string `json:"title"`
  Content string `json:"content"`
  Author string `json:"author"`
}
