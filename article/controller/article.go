package controller

import (
	"fmt"
	"github.com/femonofsky/ArticleMaker/article/model"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

// ArticleController Handler
type ArticleController struct {
	logger *log.Logger
}

// TODO: Filter by Date,
// GetAll Handler: handle get all articles and can be filter by category,publisher, created_at, published_at
func (ac *ArticleController) GetAll(w io.Writer, r *http.Request) (interface{}, int, error) {
	article := model.Article{}
	if category := r.FormValue("category"); category != "" {
		article.CategoryName = category
	}
	if publisher := r.FormValue("publisher"); publisher != "" {
		article.PublisherName = publisher
	}
	//if createdAt := r.FormValue("created_at"); createdAt != "" {
	//
	//	createdAtTime, err := time.Parse(DateTimeLayout, createdAt)
	//	if err != nil {
	//		return nil, http.StatusBadRequest , fmt.Errorf("unable to convert time %v", err)
	//	}
	//
	//	article.CreatedAt = createdAtTime
	//}
	//
	//if publishedAt := r.FormValue("published_at"); publishedAt != "" {
	//	//convert publishedAt to time object
	//	publishedAtTime, err :=  time.Parse(DateTimeLayout, publishedAt)
	//	if err != nil {
	//		return nil, http.StatusBadRequest, err
	//	}
	//	article.PublishedAt = publishedAtTime
	//}

	articles, err := model.GetArticles(article)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return articles, http.StatusOK, nil
}

// Create Handler: Create a new Article
func (ac *ArticleController) Create(w io.Writer, r *http.Request) (interface{}, int, error) {

	article, err := model.Serialize(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err  = article.Validate(); err != nil {
		return nil, http.StatusBadRequest, err
	}

	err = model.CreateArticle(article)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return article, http.StatusCreated, nil
}

// Get Handler: Get article using ID
func (ac *ArticleController) Get(w io.Writer, r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid Id got: %v", id)
	}

	article, err := model.GetArticle(model.Article{ID: uint(id)})
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return article, http.StatusOK, nil

}

// Delete Handler: delete article by ID
func (ac *ArticleController) Delete(w io.Writer, r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid Id got: %v, %v", id, err)
	}
	if err = model.DeleteArticle(id); err != nil {
		return  nil, http.StatusBadRequest, err
	}

	return nil, http.StatusNoContent, nil
}

// TODO: still giving error
// PUT Handler: Update article
func (ac *ArticleController) Put(w io.Writer, r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid Id got: %v", id)
	}
	article, err := model.Serialize(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if err  = article.Validate(); err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = model.UpdateArticle(id, article)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("unable to upload article got: %v", err)
	}
	return &article, http.StatusOK, nil

}


// newArticle creates a new Article Handle
func newArticle(logger *log.Logger) *ArticleController {
	return &ArticleController{logger:logger}
}
