package rest

import (
	"blog-api/src/dao"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	HttpResponseSuccessMessage = "Success"
	HttpResponseCreatedMessage = "Created"
	HttpResponseErrorPathParameterMessage = "must specify integer value in path paramater"
	HttpResponseErrorArticleNotFound = "article id:%d not found"
	HttpResponseIncompleteRequestMessage = "should specify title, content, and author"
)

type ArticleHandler struct {
	ArticleDao dao.ArticleDaoInterface
}

func (a *ArticleHandler)InsertHandler(w http.ResponseWriter, r *http.Request){
	articleObject := dao.ArticleObject{}
	err := json.NewDecoder(r.Body).Decode(&articleObject)
	var postResponse ArticlePostResponse
	if err != nil {
		postResponse = ArticlePostResponse{http.StatusBadRequest,
			err.Error(), nil}
		w.WriteHeader(http.StatusBadRequest)
	} else if articleObject.Title == "" || articleObject.Author == "" || articleObject.Content == "" {
		postResponse = ArticlePostResponse{http.StatusBadRequest,
			HttpResponseIncompleteRequestMessage,
			nil}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		userid, err := a.ArticleDao.Insert(articleObject.Title, articleObject.Content, articleObject.Author)
		if err != nil {
			postResponse = ArticlePostResponse{http.StatusFailedDependency, err.Error(), nil}
			w.WriteHeader(http.StatusFailedDependency)
		} else {
			postData := ArticlePostData{userid}
			postResponse = ArticlePostResponse{http.StatusCreated,
				HttpResponseCreatedMessage, &postData}
			w.WriteHeader(http.StatusCreated)
		}
	}
	json.NewEncoder(w).Encode(postResponse)
}

func (a *ArticleHandler)GetAllHandler(w http.ResponseWriter, r *http.Request){
	objects, err := a.ArticleDao.FindAll()
	var articleResponse ArticleGetAllResponse
	if err != nil {
		articleResponse = ArticleGetAllResponse{http.StatusFailedDependency, err.Error(), nil}
		w.WriteHeader(http.StatusFailedDependency)
	} else {
		articleResponse = ArticleGetAllResponse{http.StatusOK, HttpResponseSuccessMessage, &objects}
	}
	json.NewEncoder(w).Encode(articleResponse)
}

func (a *ArticleHandler)GetByIdHandler(w http.ResponseWriter, r *http.Request){
	arr := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(arr[len(arr)-1])
	var articleResponse ArticleGetIdResponse
	if err != nil {
		articleResponse = ArticleGetIdResponse{http.StatusBadRequest,
			HttpResponseErrorPathParameterMessage,
			nil}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		object, err := a.ArticleDao.FindById(id)
		if err != nil {
			articleResponse = ArticleGetIdResponse{http.StatusFailedDependency, err.Error(), nil}
			w.WriteHeader(http.StatusFailedDependency)
		} else if object.Id == 0 {
			articleResponse = ArticleGetIdResponse{http.StatusNotFound,
				fmt.Sprintf(HttpResponseErrorArticleNotFound, id), nil}
			w.WriteHeader(http.StatusNotFound)
		} else {
			articleResponse = ArticleGetIdResponse{http.StatusOK, HttpResponseSuccessMessage, &object}
		}
	}
	json.NewEncoder(w).Encode(articleResponse)
}

func (a *ArticleHandler)RemoveHandler(w http.ResponseWriter, r *http.Request){
	arr := strings.Split(r.URL.Path, "/")
	id, err:= strconv.Atoi(arr[len(arr)-1])
	var articleResponse ArticleDeleteResponse
	if err != nil {
		articleResponse = ArticleDeleteResponse{http.StatusBadRequest,
			HttpResponseErrorPathParameterMessage,
			nil}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		title, err:= a.ArticleDao.Delete(id)
		articleData := ArticleDeleteData{title}
		if err != nil {
			articleResponse = ArticleDeleteResponse{http.StatusFailedDependency,
				err.Error(), nil}
			w.WriteHeader(http.StatusFailedDependency)
		} else {
			articleResponse = ArticleDeleteResponse{http.StatusOK, HttpResponseSuccessMessage, &articleData}

		}
	}
	json.NewEncoder(w).Encode(articleResponse)
}