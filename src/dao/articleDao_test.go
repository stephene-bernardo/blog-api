package dao

import (
  "testing"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/stretchr/testify/assert"
  "fmt"
)
type article struct {
  title string
  content string
  author string
}

var goArticle = article{"Go lang", "some go lang meaningful content", "Mr. Go"}
var javaArticle = article{"Java lang", "some JAVA meaningful content", "Mr. JAVA"}
var perlArticle = article{"Perl", "for deletion testing", "Foo Bar"}

func TestArticleDao(t *testing.T){
  connStr:="user=postgres dbname=blog password=abc123 sslmode=disable"
  db, err := sql.Open("postgres", connStr)
  defer db.Close()

  assert.Nil(t, err)
  var databaseTable = "articlefortesting"
  var queryStringForDropTable = fmt.Sprintf("DROP TABLE IF EXISTs %s", databaseTable)
  var queryStringForCreateTable = fmt.Sprintf(`
    CREATE TABLE %s (
    id serial not null primary key,
    TITLE VARCHAR,
    content VARCHAR,
    author VARCHAR)
  `, databaseTable)

  db.Query(queryStringForDropTable)
  db.Query(queryStringForCreateTable)

  articleDao := ArticleDao{db, "articlefortesting"}
  articleDao.Insert(goArticle.title, goArticle.content, goArticle.author)
  articleDao.Insert(javaArticle.title, javaArticle.content, javaArticle.author)
  articleDao.Insert(perlArticle.title, perlArticle.content, perlArticle.author)

  t.Run("findAll", func(t *testing.T) {
      article, _ := articleDao.FindAll()

      assert.Equal(t, goArticle.title, article[0].Title)
      assert.Equal(t, javaArticle.title, article[1].Title)

  })

  t.Run("findById", func(t *testing.T) {
      article, _ := articleDao.FindById(1)

      assert.Equal(t, goArticle.title, article.Title)
  })

  t.Run("findById article not found", func(t *testing.T) {
    _, err:= articleDao.FindById(1000)

    assert.Equal(t, err.Error(), "Unable to find id:1000")
  })

  t.Run("delete article", func(t *testing.T) {
    articleDao.Delete(3)

    articles, _ := articleDao.FindAll()

    for _, elem := range articles {
      assert.NotEqual(t, perlArticle.title, elem)
    }
  })

  db.Query(queryStringForDropTable)
}
