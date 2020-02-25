package model

import (
	"bytes"
	"github.com/femonofsky/articleMaker/article/config"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cfg := config.Config{
		DB: config.DB{
			Driver: "sqlite3",
			Name:   "articleTest.db",
		},
	}
	log.Println("loading Database")
	db, err := New(&cfg)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	defer db.Close()
	db.Debug().AutoMigrate(&Article{}, &Category{}, &Publisher{})

	log.Println("Finished loading Database")
	if err := refreshAllTable(); err != nil {
		log.Fatal("unable to refreshTable")
	}
	ret := m.Run()

	os.Exit(ret)

}

// Clear all DB tables
func refreshAllTable() error {
	// Drop Table if Exists
	err := Db.Debug().DropTableIfExists(&Article{}, &Category{}, &Publisher{}).Error
	if err != nil {
		return err
	}

	// Migrate All table
	err = Db.Debug().AutoMigrate(Article{}, &Category{}, &Publisher{}).Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil
}

//Test Checks Validation
func TestChecksValidation(t *testing.T) {

	tests := []struct {
		name    string
		args    Article
		wantErr bool
	}{
		{"case 01", Article{ID: 0, Title: "", Body: "",
			Category: Category{}, CategoryName: "", Publisher: Publisher{}, PublisherName: "",
			CreatedAt: time.Time{}, PublishedAt: time.Time{}, UpdatedAt: time.Time{}}, true},
		{"case 02", Article{ID: 1, Title: "Love of Money", Body: "Love of Money",
			Category: Category{}, CategoryName: "Money", Publisher: Publisher{}, PublisherName: "tunde",
			CreatedAt: time.Time{}, PublishedAt: time.Time{}, UpdatedAt: time.Time{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()

			if err != nil && !tt.wantErr {
				t.Errorf("unable to validate data:%v", err)
			}
			if tt.wantErr {
				return
			}
		})
	}
}

// Serialize
func TestSerialize(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    Article
		wantErr bool
	}{
		{"case 01", []byte(`{ "id": 0, "title": "Money", "body":"Money is good", "category": "social",
			"publisher": "femonofsky"
		}`),
			Article{Title: "Money", Body: "Money is good", CategoryName: "social",
				PublisherName: "femonofsky"}, false},
		{"case 02", []byte(`{ "id": 0, "title": "Money", "body":"Money is good", "category": "social",
			"publisher": "femonofsky
		}`),
			Article{Title: "Money", Body: "Money is good", CategoryName: "social",
				PublisherName: "femonofsky"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := ioutil.NopCloser(bytes.NewBuffer(tt.args))
			got, err := Serialize(body)
			if err != nil && !tt.wantErr {
				t.Errorf("unable validate data:%v", err)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("Serialize() = %v, want %v", got, &tt.want)
			}
		})
	}
}

// Create Article
func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name    string
		args    Article
		wantErr bool
	}{
		{"case 01", Article{Title: "Money", Body: "Money is good", CategoryName: "social",
			PublisherName: "femonofsky"}, false},
		{"case 02", Article{Title: "Love of Money", Body: "Love of Money", CategoryName: "Money",
			PublisherName: "tunde"}, false},
		{"case 03", Article{Title: "Love of Money", Body: "Love of Money", CategoryName: "Money",
			PublisherName: "tunde"}, true},
		{"case 04", Article{Title: "", Body: "Love of Money", CategoryName: "Money",
			PublisherName: "tunde"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateArticle(&tt.args)
			if err != nil && !tt.wantErr {
				t.Errorf("unable validate data:%v", err)
			}
		})
	}
}

// Delete Article *
func TestDeleteArticle(t *testing.T) {
	if err := refreshAllTable(); err != nil {
		log.Fatal("unable to refreshTable")
	}
	tests := []struct {
		name    string
		args    Article
		wantErr bool
	}{
		{"case 01", Article{Title: "Money", Body: "Money is good", CategoryName: "social",
			PublisherName: "femonofsky"}, false},
		{"case 02", Article{Title: "Love of Money", Body: "Love of Money", CategoryName: "Money",
			PublisherName: "tunde"}, false},
		{"case 03", Article{}, true},
		{"case 04", Article{Title: "", Body: "Love of Money", CategoryName: "Money",
			PublisherName: "tunde"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arti := tt.args
			CreateArticle(&arti)
			err := DeleteArticle(int(arti.ID))
			if err != nil && !tt.wantErr {
				t.Errorf("unable validate data:%v", err)
			}
		})
	}
}

// Update Article
func TestUpdateArticle(t *testing.T) {

	if err := refreshAllTable(); err != nil {
		log.Fatal("unable to refreshTable")
	}
	tests := []struct {
		name    string
		args    Article
		want    Article
		wantErr bool
	}{
		{"case 01", Article{Title: "Money", Body: "Money is good", CategoryName: "social",
			PublisherName: "femonofsky"}, Article{Title: "Money", Body: "Money is good sjdnaj", CategoryName: "social",
			PublisherName: "femonofsky"}, false},
		{"case 02", Article{Title: "Love of Money", Body: "Love of Money", CategoryName: "Money",
			PublisherName: "tunde"}, Article{Title: "Money", Body: "Money is good", CategoryName: "social",
			PublisherName: "femonofsky"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arti := tt.args
			CreateArticle(&arti)
			err := UpdateArticle(int(arti.ID), &tt.want)
			if err != nil && !tt.wantErr {
				t.Errorf("unable validate data:%v", err)
			}
			if tt.wantErr {
				return
			}

		})
	}
}

// Get All Articles
func TestGetArticles(t *testing.T) {
	if err := refreshAllTable(); err != nil {
		t.Error("unable to refreshTable")
	}
	// Create Dummy Data
	articles := Articles{&Article{Title: "Money", Body: "Money is good", CategoryName: "social",
		PublisherName: "femonofsky"}, &Article{Title: "Love of Money", Body: "Love of Money", CategoryName: "Money",
		PublisherName: "tunde"}}
	for _, article := range articles {
		err := CreateArticle(article)
		if err != nil {
			t.Errorf("unable to create new article %v", err)
		}
	}

	tests := []struct {
		name  string
		arg   Article
		count int
	}{
		{"case 01", Article{}, 2},
		{"case 02", Article{Title: "Love of Money"}, 1},
		{"case 03", Article{CategoryName: "Money"}, 1},
		{"case 04", Article{CategoryName: "Money23"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetArticles(tt.arg)
			if err != nil {
				t.Errorf("unable to validate data:%v", err)
			}
			if len(got) != tt.count {
				t.Errorf("filter is not working got : %v want: %v", got, tt.count)
			}

		})
	}
}
