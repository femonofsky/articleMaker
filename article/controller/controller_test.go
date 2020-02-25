package controller

import (
	"fmt"
	"github.com/femonofsky/articleMaker/article/config"
	"github.com/femonofsky/articleMaker/article/model"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	cfg := config.Config{
		DB: config.DB{
			Driver: "sqlite3",
			Name:   "articleTest.db",
		},
	}
	log.Println("loading Database")
	db, err := model.New(&cfg)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	defer db.Close()
	db.Debug().AutoMigrate(&model.Article{}, &model.Category{}, &model.Publisher{})

	log.Println("Finished loading Database")
	srv := httptest.NewServer(New(nil))
	defer srv.Close()
	server = srv
	if err := refreshAllTable(db); err != nil {
		log.Fatal("unable to refreshTable")
	}
	ret := m.Run()

	os.Exit(ret)

}

// Clear all DB tables
func refreshAllTable(Db *gorm.DB) error {
	// Drop Table if Exists
	err := Db.Debug().DropTableIfExists(&model.Article{}, &model.Category{}, &model.Publisher{}).Error
	if err != nil {
		return err
	}

	// Migrate All table
	err = Db.Debug().AutoMigrate(model.Article{}, &model.Category{}, &model.Publisher{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func TestNewArticleController_Create(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{"case 01", `{ "title": "Tommy test","body": "Andela is the best office to work in ajjfdsfdfskjfv ",
								"category": "Extras","publisher": "Femonofsky"}`, false},
		{"case 02", `{ "title": "Tommy test","body": "Andela is the best office to work in ajjfdsfdfskjfv ",
								"category": "Extras","publisher": "Femonofsky"}`, true},
		{"case 03", `{ "body": "Andela is the best office to work in ajjfdsfdfskjfv ",
								"category": "Extras","publisher": "Femonofsky"}`, true},
		{"case 04", `{ "title": "Tommy test","category": "Extras","publisher": "Femonofsky"}`, true},
		{"case 05", `{ "title": "Money in the bank","body": "Andela is the best office to work in ajjfdsfdfskjfv ",
								"category": "Extras","publisher": "Femonofsky"}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := http.Post(fmt.Sprintf("%s/article/", server.URL),
				"application/json", strings.NewReader(tt.args))
			if err != nil {
				t.Fatalf("could not send POST request: %v", err)
			}
			if res.StatusCode != http.StatusCreated && !tt.wantErr {
				t.Errorf("expected status created; got %v", res.Status)
			}
			if tt.wantErr {
				return
			}

			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}

func TestNewArticleController_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		want    int
		wantErr bool
	}{
		{"case 01", `?category=Extras`, 1, false},
		{"case 02", `?title=Extras`, 1, false},
		{"case 03", `?published_at=Tommy`, 0, true},
		{"case 04", `?publisher=Femonofsky`, 1, false},
		{"case 05", ``, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/article/%s", server.URL, tt.filter)
			res, err := http.Get(url)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			if res.StatusCode != http.StatusOK && !tt.wantErr {
				t.Errorf("expected status ok; got %v", res.Status)
			}
			if tt.wantErr {
				return
			}

			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}

func TestNewArticleController_Get(t *testing.T) {
	tests := []struct {
		name    string
		params  string
		wantErr bool
	}{
		{"case 01", `1`, false},
		{"case 02", `990`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/article/%s", server.URL, tt.params)
			res, err := http.Get(url)
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			if res.StatusCode != http.StatusOK && !tt.wantErr {
				t.Errorf("expected status ok; got %v", res.Status)
			}
			if tt.wantErr {
				return
			}

			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}

func TestNewArticleController_Put(t *testing.T) {
	tests := []struct {
		name    string
		params  string
		body    string
		wantErr bool
	}{
		{"case 01", `1`, `{"body": "money in the bank"}`, false},
		{"case 02", `990`, `{"body": "money in the bank"}`, true},
		{"case 03", `2`, `{"title": "Tommy test"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}
			url := fmt.Sprintf("%s/article/%s", server.URL, tt.params)
			req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("could not send PUT request: %v", err)
			}
			res, err := client.Do(req)
			if res.StatusCode != http.StatusOK && !tt.wantErr {
				t.Errorf("expected status ok; got %v", res.Status)
			}
			if tt.wantErr {
				return
			}
			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}

func TestNewArticleController_Delete(t *testing.T) {
	tests := []struct {
		name    string
		params  string
		wantErr bool
	}{
		{"case 01", `1`, false},
		{"case 02", `990`, true},
		{"case 03", `2`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := http.Client{}
			url := fmt.Sprintf("%s/article/%s", server.URL, tt.params)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			if err != nil {
				t.Fatalf("could not send DELETE request: %v", err)
			}
			res, err := client.Do(req)
			if res.StatusCode != http.StatusNoContent && !tt.wantErr {
				t.Errorf("expected status ok; got %v", res.Status)
			}
			if tt.wantErr {
				return
			}
			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}
