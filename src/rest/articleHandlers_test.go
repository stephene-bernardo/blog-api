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
var articleObject dao.ArticleObject

func TestMain(m *testing.M){
	mockArticleDao = new(mocks.ArticleDaoInterface)
	articleObject = dao.ArticleObject{1, "Java Lang", "Some Content", "Mr.Java"}
	m.Run()
}

func TestArticleHandler_GetAllHandler(t *testing.T) {
	mockArticleDao.On("FindAll").Return([]dao.ArticleObject{articleObject}, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("GET", "/articles", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetAllHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleGetAllResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, http.StatusOK, response.Status)
	assert.EqualValues(t, HttpResponseSuccessMessage, response.Message)
	assertArticle(t, articleObject, response.Data[0])
}

func TestArticleHandler_GetByIdHandler(t *testing.T) {
	mockArticleDao.On("FindById", articleObject.Id).Return(articleObject, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleGetIdResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, http.StatusOK, response.Status)
	assert.EqualValues(t, HttpResponseSuccessMessage, response.Message)
	assertArticle(t, articleObject, response.Data)
}

func assertArticle(t *testing.T, expectedArticle dao.ArticleObject, actualArticle dao.ArticleObject){
	assert.EqualValues(t, expectedArticle.Id, actualArticle.Id)
	assert.EqualValues(t, expectedArticle.Title, actualArticle.Title)
	assert.EqualValues(t, expectedArticle.Content, actualArticle.Content)
	assert.EqualValues(t, expectedArticle.Author, actualArticle.Author)
}

func TestArticleHandler_InsertHandler(t *testing.T) {
	mockArticleDao.On(
		"Insert",
		articleObject.Title,
		articleObject.Content,
		articleObject.Author).Return(articleObject.Id, nil)
	articleHandlers := ArticleHandler{mockArticleDao}
	jsonArticle, _ := json.Marshal(articleObject)

	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.InsertHandler)
	handler.ServeHTTP(rr, req)

	var response ArticlePostResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
	assert.EqualValues(t, http.StatusCreated, response.Status)
	assert.EqualValues(t, HttpResponseCreatedMessage, response.Message)
	assert.EqualValues(t, articleObject.Id, response.Data.Id)
}


func TestArticleHandler_DeleteHandler(t *testing.T) {
	mockArticleDao.On("Delete", articleObject.Id).Return(articleObject.Title, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("DELETE", "/articles/1", nil)
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.RemoveHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleDeleteResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, http.StatusOK, response.Status)
	assert.EqualValues(t, HttpResponseSuccessMessage, response.Message)
	assert.EqualValues(t, articleObject.Title, response.Data.Title)
}