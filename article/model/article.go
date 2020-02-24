package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
	"time"
)

// Article Defines the structure for an API article
type Article struct {
	ID            uint      `gorm:"primary_key;auto_increment"`
	Title         string    `sql:"unique;unique_index;not null" json:"title" validate:"required"`
	Body          string    `sql:"not null" json:"body" validate:"required"`
	Category      Category  `gorm:"association_foreignKey:CategoryName" json:"-"`
	CategoryName  string    `json:"category" validate:"required"`
	Publisher     Publisher `gorm:"association_foreignKey:PublisherName" json:"-"`
	PublisherName string    `json:"publisher" validate:"required"`
	CreatedAt     time.Time `json:"created_at"`
	PublishedAt   time.Time `json:"published_at" `
	UpdatedAt     time.Time `json:"-"`
}

// UnmarshalJSON parses the json string in the custom format
func (article *Article) UnmarshalJSON(data []byte) (err error) {
	var auxArticle struct {
		Title         string `json:"title"`
		Body          string `json:"body"`
		CategoryName  string `json:"category" `
		PublisherName string `json:"publisher"`
		PublishedAt   string `json:"published_at" `
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(&auxArticle); err != nil {
		return fmt.Errorf("unable to decode %v", err)
	}
	article.Title = auxArticle.Title
	article.Body = auxArticle.Body
	article.CategoryName = auxArticle.CategoryName
	article.PublisherName = auxArticle.PublisherName
	if auxArticle.PublishedAt != "" {
		publishedAt, err := time.Parse(DateTimeLayout, auxArticle.PublishedAt)
		if err != nil {
			return fmt.Errorf("invalid format use (%v format) : %v", DateTimeLayout, err)
		}
		article.PublishedAt = publishedAt
	}
	return nil

}

//MarshalJSON writes a quotes string in the custom format
func (article *Article) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID            uint   `json:"id"`
		Title         string `json:"title"`
		Body          string `json:"body"`
		CategoryName  string `json:"category"`
		PublisherName string `json:"publisher"`
		CreatedAt     string `json:"created_at"`
		PublishedAt   string `json:"published_at"`
	}{
		ID:            article.ID,
		Title:         article.Title,
		Body:          article.Body,
		CategoryName:  article.CategoryName,
		PublisherName: article.PublisherName,
		CreatedAt:     article.CreatedAt.Format(DateTimeLayout),
		PublishedAt:   article.PublishedAt.Format(DateTimeLayout),
	})
}

// Validate: check if all Article fields met requirements
func (article *Article) Validate() error {
	validate := validator.New()
	return validate.Struct(article)
}

// BeforeSave Gorm trigger before saving the article
func (article *Article) BeforeSave() error {
	category := Category{}
	if err := Db.FirstOrCreate(&category, Category{Name: article.CategoryName}).Error; err != nil {
		return fmt.Errorf("could not reference category got: %v", err)
	}

	article.Category = category
	publisher := Publisher{}
	if err := Db.FirstOrCreate(&publisher, Publisher{Name: article.PublisherName}).Error; err != nil {
		return fmt.Errorf("could not reference category got: %v", err)
	}

	article.Publisher = publisher

	return nil

}

// Articles of a collection of Article
type Articles []*Article

// GetArticles returns a slice of articles
func GetArticles(article Article) (Articles, error) {
	articles := Articles{}
	if err := Db.Find(&articles, article).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// CreateArticle create new  Article
func CreateArticle(article *Article) error {
	arr, err := GetArticle(Article{Title: article.Title})
	if arr != nil || err == nil {
		return fmt.Errorf("title aleady exists")
	}
	if err := Db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

// UpdateArticle update Article using ID to check if the article exist
func UpdateArticle(id int, article *Article) error {
	arr, err := GetArticle(Article{ID: uint(id)})
	if arr == nil || err != nil {
		return err
	}

	if err = Db.Debug().Model(arr).Update(article).Error; err != nil {
		return err
	}
	return nil
}

// DeleteArticle Article using article id
func DeleteArticle(id int) error {
	articles, err := GetArticle(Article{ID: uint(id)})
	if articles == nil || err != nil {
		return err
	}
	if err = Db.Delete(&articles).Error; err != nil {
		return err
	}
	return nil
}

// Serialize convert request into Article object
func Serialize(r io.ReadCloser) (*Article, error) {

	article := &Article{}

	// read the JSON-encoded value and decode into slice of Articles
	if err := json.NewDecoder(r).Decode(article); err != nil {
		return nil, fmt.Errorf("unable to decode json request body: %v", err)
	}
	defer r.Close()
	return article, nil

}

// ErrArticleNotFound article not found error
var ErrArticleNotFound = fmt.Errorf("article not found")

// GetArticle get article by ID
func GetArticle(query interface{}) (*Article, error) {
	articles := &Article{}
	if err := Db.First(articles, query).Error; err != nil {
		return nil, ErrArticleNotFound
	}

	return articles, nil
}
