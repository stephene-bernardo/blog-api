package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type article struct {
	title   string
	content string
	author  string
}

var goArticle = article{"Go lang", "some go lang meaningful content", "Mr. Go"}
var javaArticle = article{"Java lang", "some JAVA meaningful content", "Mr. JAVA"}
var perlArticle = article{"Perl", "for deletion testing", "Foo Bar"}

func TestArticleDao(t *testing.T) {
	host := "localhost"
	fmt.Println(os.Getenv("POSTGRES_HOST"))
	if os.Getenv("POSTGRES_HOST") != ""{
		host = os.Getenv("POSTGRES_HOST")
	}
	connStr := fmt.Sprintf("user=postgres dbname=postgres password=abc123 sslmode=disable host=%s", host)
	db, err := sql.Open("postgres", connStr)
	assert.Nil(t, err)
	defer db.Close()
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
		_, err := articleDao.FindById(1000)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "Unable to find id:1000")
	})

	t.Run("delete article", func(t *testing.T) {
		title, _ := articleDao.Delete(3)

		assert.Equal(t, perlArticle.title, title)
	})

	db.Query(queryStringForDropTable)
}

func TestArticleDao_WithInvalidDbConnection(t *testing.T) {
	invalidDbConnectionStr := "user=NotExisting dbname=WrongDb password=WrongPassword sslmode=disable"
	db, err := sql.Open("postgres", invalidDbConnectionStr)
	assert.Nil(t, err)
	defer db.Close()

	articleDao := ArticleDao{db, "articlefortesting"}

	t.Run("findAll", func(t *testing.T) {
		_, err := articleDao.FindAll()

		assert.Error(t, err)
	})

	t.Run("findById", func(t *testing.T) {
		_, err := articleDao.FindById(1)

		assert.Error(t, err)
	})

	t.Run("delete article", func(t *testing.T) {
		_, err := articleDao.Delete(3)

		assert.Error(t, err)
	})

	t.Run("insert article", func(t *testing.T) {
		_, err := articleDao.Insert(goArticle.title, goArticle.content, goArticle.author)

		assert.Error(t, err)
	})
}
