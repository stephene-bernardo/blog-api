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

func TestMain(m *testing.M) {
	articleObject = dao.ArticleObject{1, "Java Lang", "Some Content", "Mr.Java"}
	m.Run()
}

func TestArticleHandler_GetAllHandler(t *testing.T) {
	t.Run("getAll", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On("FindAll").Return([]dao.ArticleObject{articleObject}, nil)
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("GET", "/articles", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.GetAllHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleGetAllResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusOK, rr.Code, response.Status)
		assertArticle(t, articleObject, (*response.Data)[0])
		assert.EqualValues(t, HttpResponseSuccessMessage, response.Message)
	})
	t.Run("should not get all when database error occur", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On("FindAll").Return(nil, errors.New(databaseConnectionErrorMessage))
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("GET", "/articles", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.GetAllHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleGetAllResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusFailedDependency, rr.Code, response.Status)
		assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
		assert.Nil(t, response.Data)
	})
}

func TestArticleHandler_GetByIdHandler(t *testing.T){
	t.Run("getById", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On("FindById", articleObject.Id).Return(articleObject, nil)
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("GET", "/articles/1", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleGetIdResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusOK, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseSuccessMessage, response.Message)
		assertArticle(t, articleObject, *response.Data)
	})

	t.Run("should not getById when wrong path parameter type", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On("FindById", articleObject.Id).Return(dao.ArticleObject{}, nil)
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("GET", "/articles/wrongtype", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.GetByIdHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleGetIdResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusBadRequest, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseErrorPathParameterMessage, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("should not getById when database error occur", func(t *testing.T) {
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
		assertStatusCode(t, http.StatusFailedDependency, rr.Code, response.Status)
		assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("should get response when article not found", func(t *testing.T) {
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
		assertStatusCode(t, http.StatusNotFound, rr.Code, response.Status)
		assert.EqualValues(t, fmt.Sprintf(HttpResponseErrorArticleNotFound, missingArticle), response.Message)
		assert.Nil(t, response.Data)
	})
}

func assertArticle(t *testing.T, expectedArticle dao.ArticleObject, actualArticle dao.ArticleObject) {
	assert.EqualValues(t, expectedArticle.Id, actualArticle.Id)
	assert.EqualValues(t, expectedArticle.Title, actualArticle.Title)
	assert.EqualValues(t, expectedArticle.Content, actualArticle.Content)
	assert.EqualValues(t, expectedArticle.Author, actualArticle.Author)
}

func TestArticleHandler_InsertHandler(t *testing.T) {
	t.Run("should insert article", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On(
			"Insert",
			articleObject.Title,
			articleObject.Content,
			articleObject.Author).Return(articleObject.Id, nil)
		articleHandlers := ArticleHandler{mockArticleDao}
		jsonArticle, _ := json.Marshal(articleObject)

		req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.InsertHandler)
		handler.ServeHTTP(rr, req)

		var response ArticlePostResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusCreated, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseCreatedMessage, response.Message)
		assert.EqualValues(t, articleObject.Id, response.Data.Id)
	})

	t.Run("should handle insert database error", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On(
			"Insert",
			articleObject.Title,
			articleObject.Content,
			articleObject.Author).Return(0, errors.New(databaseConnectionErrorMessage))
		articleHandlers := ArticleHandler{mockArticleDao}
		jsonArticle, _ := json.Marshal(articleObject)

		req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.InsertHandler)
		handler.ServeHTTP(rr, req)

		var response ArticlePostResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusFailedDependency, rr.Code, response.Status)
		assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
		assert.Nil(t, response.Data)
	})
	t.Run("should handle missing title error", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		articleHandlers := ArticleHandler{mockArticleDao}
		jsonArticle, _ := json.Marshal(dao.ArticleObject{0, "",
			articleObject.Content, articleObject.Author})

		req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.InsertHandler)
		handler.ServeHTTP(rr, req)

		var response ArticlePostResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusBadRequest, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseIncompleteRequestMessage, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("should handle missing content error", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		articleHandlers := ArticleHandler{mockArticleDao}
		jsonArticle, _ := json.Marshal(dao.ArticleObject{0, articleObject.Title,
			"", articleObject.Author})

		req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.InsertHandler)
		handler.ServeHTTP(rr, req)

		var response ArticlePostResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusBadRequest, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseIncompleteRequestMessage, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("should handle missing author error", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		articleHandlers := ArticleHandler{mockArticleDao}
		jsonArticle, _ := json.Marshal(dao.ArticleObject{0, articleObject.Title,
			articleObject.Content, ""})

		req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonArticle))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.InsertHandler)
		handler.ServeHTTP(rr, req)

		var response ArticlePostResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusBadRequest, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseIncompleteRequestMessage, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("should handle invalid json format", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer( []byte(`{invalidjsonformat`)))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.InsertHandler)
		handler.ServeHTTP(rr, req)

		var response ArticlePostResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusBadRequest, rr.Code, response.Status)
		assert.NotEmpty(t,  response.Message)
		assert.Nil(t, response.Data)
	})
}


func TestArticleHandler_DeleteHandler(t *testing.T) {
	t.Run("should handle delete article", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On("Delete", articleObject.Id).Return(articleObject.Title, nil)
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("DELETE", "/articles/1", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.RemoveHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleDeleteResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusOK, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseSuccessMessage, response.Message)
		assert.EqualValues(t, articleObject.Title, response.Data.Title)
	})

	t.Run("should handle delete database error", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		mockArticleDao.On("Delete", articleObject.Id).Return("",
			errors.New(databaseConnectionErrorMessage))
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("DELETE", "/articles/1", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.RemoveHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleDeleteResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusFailedDependency, rr.Code, response.Status)
		assert.EqualValues(t, databaseConnectionErrorMessage, response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("should handle delete path parameter type error", func(t *testing.T) {
		mockArticleDao = new(mocks.ArticleDaoInterface)
		articleHandlers := ArticleHandler{mockArticleDao}

		req, _ := http.NewRequest("DELETE", "/articles/wrongType", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(articleHandlers.RemoveHandler)
		handler.ServeHTTP(rr, req)

		var response ArticleDeleteResponse
		json.NewDecoder(rr.Body).Decode(&response)
		assertStatusCode(t, http.StatusBadRequest, rr.Code, response.Status)
		assert.EqualValues(t, HttpResponseErrorPathParameterMessage, response.Message)
		assert.Nil(t, response.Data)
	})
}

func assertStatusCode(t *testing.T, expected ,actualHeaderStatusCode, actualBodyStatusCode int) {
	assert.EqualValues(t, expected, actualHeaderStatusCode)
	assert.EqualValues(t, expected, actualBodyStatusCode)
}