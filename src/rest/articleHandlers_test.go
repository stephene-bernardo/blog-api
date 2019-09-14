package rest

import (
	"blog-api/mocks"
	"blog-api/src/dao"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockArticleDao *mocks.ArticleDaoInterface
var articleObject dao.ArticleObject
const databaseConnectionErrorMessage = "connection to database failed"

func TestMain(m *testing.M){
	articleObject = dao.ArticleObject{1, "Java Lang", "Some Content", "Mr.Java"}
	m.Run()
}

func TestArticleHandler_GetAllHandler(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
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
	assertArticle(t, articleObject, (*response.Data)[0])
}

func TestArticleHandler_GetAllHandlerDatabaseError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	mockArticleDao.On("FindAll").Return(nil, errors.New(databaseConnectionErrorMessage))
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("GET", "/articles", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetAllHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleGetAllResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusFailedDependency, rr.Code)
	assert.EqualValues(t, http.StatusFailedDependency, response.Status)
	assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_GetByIdHandler(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
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
	assertArticle(t, articleObject, *response.Data)
}

func TestArticleHandler_GetByIdHandlerErrorInPathParameterType(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	mockArticleDao.On("FindById", articleObject.Id).Return(dao.ArticleObject{}, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("GET", "/articles/wrongtype", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleGetIdResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	assert.EqualValues(t, http.StatusBadRequest, response.Status)
	assert.EqualValues(t, HttpResponseErrorPathParameterMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_GetByIdHandlerDatabaseError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	mockArticleDao.On("FindById", articleObject.Id).Return(articleObject,
		errors.New(databaseConnectionErrorMessage))
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleGetIdResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusFailedDependency, rr.Code)
	assert.EqualValues(t, http.StatusFailedDependency, response.Status)
	assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_GetByIdHandlerArticleNotFound(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	const missingArticle = 2
	mockArticleDao.On("FindById", missingArticle).Return(dao.ArticleObject{}, nil)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("GET", "/articles/2", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleGetIdResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusNotFound, rr.Code)
	assert.EqualValues(t, http.StatusNotFound, response.Status)
	assert.EqualValues(t, fmt.Sprintf(HttpResponseErrorArticleNotFound, missingArticle), response.Message)
	assert.Nil(t, response.Data)
}


func assertArticle(t *testing.T, expectedArticle dao.ArticleObject, actualArticle dao.ArticleObject){
	assert.EqualValues(t, expectedArticle.Id, actualArticle.Id)
	assert.EqualValues(t, expectedArticle.Title, actualArticle.Title)
	assert.EqualValues(t, expectedArticle.Content, actualArticle.Content)
	assert.EqualValues(t, expectedArticle.Author, actualArticle.Author)
}

func TestArticleHandler_InsertHandler(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
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

func TestArticleHandler_InsertHandlerDatabaseError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	mockArticleDao.On(
		"Insert",
		articleObject.Title,
		articleObject.Content,
		articleObject.Author).Return(0, errors.New(databaseConnectionErrorMessage))
	articleHandlers := ArticleHandler{mockArticleDao}
	jsonArticle, _ := json.Marshal(articleObject)

	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.InsertHandler)
	handler.ServeHTTP(rr, req)

	var response ArticlePostResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusFailedDependency, rr.Code)
	assert.EqualValues(t, http.StatusFailedDependency, response.Status)
	assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_InsertHandlerMissingTitleError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	articleHandlers := ArticleHandler{mockArticleDao}
	jsonArticle, _ := json.Marshal(dao.ArticleObject{0, "",
		articleObject.Content, articleObject.Author})

	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.InsertHandler)
	handler.ServeHTTP(rr, req)

	var response ArticlePostResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	assert.EqualValues(t, http.StatusBadRequest, response.Status)
	assert.EqualValues(t, HttpResponseIncompleteRequestMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_InsertHandlerMissingContentError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	articleHandlers := ArticleHandler{mockArticleDao}
	jsonArticle, _ := json.Marshal(dao.ArticleObject{0, articleObject.Title,
		"", articleObject.Author})

	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.InsertHandler)
	handler.ServeHTTP(rr, req)

	var response ArticlePostResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	assert.EqualValues(t, http.StatusBadRequest, response.Status)
	assert.EqualValues(t, HttpResponseIncompleteRequestMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_InsertHandlerMissingAuthorError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	articleHandlers := ArticleHandler{mockArticleDao}
	jsonArticle, _ := json.Marshal(dao.ArticleObject{0, articleObject.Title,
		articleObject.Content, ""})

	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.InsertHandler)
	handler.ServeHTTP(rr, req)

	var response ArticlePostResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	assert.EqualValues(t, http.StatusBadRequest, response.Status)
	assert.EqualValues(t, HttpResponseIncompleteRequestMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_DeleteHandler(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
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

func TestArticleHandler_DeleteHandlerDatabaseError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	mockArticleDao.On("Delete", articleObject.Id).Return("",
		errors.New(databaseConnectionErrorMessage))
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("DELETE", "/articles/1", nil)
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.RemoveHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleDeleteResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusFailedDependency, rr.Code)
	assert.EqualValues(t, http.StatusFailedDependency, response.Status)
	assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestArticleHandler_DeleteHandlerPathParameterTypeError(t *testing.T) {
	mockArticleDao = new(mocks.ArticleDaoInterface)
	articleHandlers := ArticleHandler{mockArticleDao}

	req, _ := http.NewRequest("DELETE", "/articles/wrongType", nil)
	rr :=  httptest.NewRecorder()
	handler := http.HandlerFunc(articleHandlers.RemoveHandler)
	handler.ServeHTTP(rr, req)

	var response ArticleDeleteResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	assert.EqualValues(t, http.StatusBadRequest, response.Status)
	assert.EqualValues(t, HttpResponseErrorPathParameterMessage, response.Message)
	assert.Nil(t, response.Data)
}