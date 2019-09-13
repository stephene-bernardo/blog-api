package rest

import (
	"blog-api/src/dao"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	HttpResponseSuccessMessage = "Success"
	HttpResponseCreatedMessage = "Created"
)

type ArticleHandler struct {
	ArticleDao dao.ArticleDaoInterface
}

func (a *ArticleHandler)InsertHandler(w http.ResponseWriter, r *http.Request){
	articleObject := dao.ArticleObject{}
	json.NewDecoder(r.Body).Decode(&articleObject)
	userid, _ := a.ArticleDao.Insert(articleObject.Title, articleObject.Content, articleObject.Author)
	postData := ArticlePostData{userid}
	postResponse := ArticlePostResponse{201, HttpResponseCreatedMessage, postData}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(postResponse)
}

func (a *ArticleHandler)GetAllHandler(w http.ResponseWriter, r *http.Request){
	objects, _ := a.ArticleDao.FindAll()
	articleResponse := ArticleGetAllResponse{200, HttpResponseSuccessMessage, objects}
	json.NewEncoder(w).Encode(articleResponse)
}

func (a *ArticleHandler)GetByIdHandler(w http.ResponseWriter, r *http.Request){
	arr := strings.Split(r.URL.Path, "/")
	id, _:= strconv.Atoi(arr[len(arr)-1])
	object, _ := a.ArticleDao.FindById(id)
	articleResponse := ArticleGetIdResponse{200, HttpResponseSuccessMessage, object}
	json.NewEncoder(w).Encode(articleResponse)
}

func (a *ArticleHandler)RemoveHandler(w http.ResponseWriter, r *http.Request){
	arr := strings.Split(r.URL.Path, "/")
	id, _:= strconv.Atoi(arr[len(arr)-1])
	title, _:= a.ArticleDao.Delete(id)
	articleData := ArticleDeleteData{title}
	articleResponse := ArticleDeleteResponse{200, HttpResponseSuccessMessage, articleData}
	json.NewEncoder(w).Encode(articleResponse)
}