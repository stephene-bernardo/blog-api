package rest

import (
	"blog-api/mocks"
	"blog-api/src/dao"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockArticleDao *mocks.ArticleDaoInterface
func TestMain(m *testing.M){
	mockArticleDao = new(mocks.ArticleDaoInterface)
	m.Run()
}

func TestArticleHandler_GetAllHandler(t *testing.T) {
	article1 := dao.ArticleObject{1, "Java Lang", "SomeContent", "Mr.Java"}
	mockArticleDao.On("FindAll").Return([]dao.ArticleObject{article1}, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetAllHandler)
	handler.ServeHTTP(rr, req)

	mockArticleDao.AssertCalled(t, "FindAll")
	var list ArticleGetAllResponse
	json.NewDecoder(rr.Body).Decode(&list)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, http.StatusOK, list.Status)
	assert.EqualValues(t, "Success", list.Message)
	assert.EqualValues(t, 1, list.Data[0].Id)
	assert.EqualValues(t, "Mr.Java", list.Data[0].Author)
	assert.EqualValues(t, "SomeContent", list.Data[0].Content)
	assert.EqualValues(t, "Java Lang", list.Data[0].Title)
}

func TestArticleHandler_GetByIdHandler(t *testing.T) {
	mockArticleDao := new(mocks.ArticleDaoInterface)
	article1 := dao.ArticleObject{1, "Java Lang", "SomeContent", "Mr.Java"}
	mockArticleDao.On("FindById", 1).Return(article1, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, err := http.NewRequest("GET", "/articles/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
	handler.ServeHTTP(rr, req)

	mockArticleDao.AssertCalled(t, "FindById", 1)
	var list ArticleGetIdResponse
	json.NewDecoder(rr.Body).Decode(&list)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, http.StatusOK, list.Status)
	assert.EqualValues(t, "Success", list.Message)
	assert.EqualValues(t, 1, list.Data.Id)
	assert.EqualValues(t, "Mr.Java", list.Data.Author)
	assert.EqualValues(t, "SomeContent", list.Data.Content)
	assert.EqualValues(t, "Java Lang", list.Data.Title)
}

func TestArticleHandler_InsertHandler(t *testing.T) {
	mockArticleDao := new(mocks.ArticleDaoInterface)
	article1 := dao.ArticleObject{1, "Java Lang", "SomeContent", "Mr.Java"}
	mockArticleDao.On("Insert", "Java Lang", "SomeContent", "Mr.Java").Return(1, nil)
	articleHandlers := ArticleHandler{mockArticleDao}
	s, _ := json.Marshal(article1)
	b := bytes.NewBuffer(s)
	req, err := http.NewRequest("POST", "/articles", b)
	if err != nil {
		t.Fatal(err)
	}

	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.InsertHandler)
	handler.ServeHTTP(rr, req)

	mockArticleDao.AssertCalled(t, "Insert", "Java Lang", "SomeContent", "Mr.Java")
	var list ArticlePostResponse
	json.NewDecoder(rr.Body).Decode(&list)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
	assert.EqualValues(t, http.StatusCreated, list.Status)
	assert.EqualValues(t, "Success", list.Message)
	assert.EqualValues(t, 1, list.Data.Id)
}


func TestArticleHandler_DeleteHandler(t *testing.T) {
	mockArticleDao := new(mocks.ArticleDaoInterface)
	mockArticleDao.On("Delete", 1).Return("Java Lang", nil)
	articleHandlers := ArticleHandler{mockArticleDao}
	req, err := http.NewRequest("DELETE", "/articles/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.RemoveHandler)
	handler.ServeHTTP(rr, req)

	mockArticleDao.AssertCalled(t, "Delete", 1)
	var list ArticleDeleteResponse
	json.NewDecoder(rr.Body).Decode(&list)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, http.StatusOK, list.Status)
	assert.EqualValues(t, "Success", list.Message)
	assert.EqualValues(t, "Java Lang", list.Data.Title)
}