package rest

import "blog-api/src/dao"

type ArticlePostResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data ArticlePostData `json:"data"`
}

type ArticlePostData struct {
	Id int`json:"id"`
}

type ArticleGetIdResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data dao.ArticleObject `json:"data"`
}


type ArticleGetAllResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data []dao.ArticleObject `json:"data"`
}

type ArticleDeleteResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data ArticleDeleteData `json:"data"`
}

type ArticleDeleteData struct {
	Title string`json:"title"`
}


