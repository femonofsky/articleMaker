# Go Article REST API
A RESTful API for simple Article application with Go

This making simple RESTful API with Go using **gorilla/mux** (A nice mux library) and **gorm** (An ORM for Go)

## Installation & Run

```bash
# Download this project
 go get github.com/femonofsky/articleMaker/article
```

Before running API server, you should set the database config with yours or set the your database config with my values in the config folder
and  [config.json](https://github.com/femonofsky/articleMaker/blob/master/article/config/config.json)
```go
{
  "server": {
    "host": "127.0.0.1",
    "port": "8080"
  },
  "db": {
    "driver": "postgres",
    "host": "127.0.0.1",
    "user": "postgres",
    "password": "",
    "name": "articledb",
    "port": "5432"
  }
}
```

```bash
# Build and Run
cd article
go build 
./article

# API Endpoint : http://127.0.0.1:8000
```

## Structure
```
├── article
|   ├── config              // Configuration Folder
│   │   ├── config.go       // Helps with the manipulation of config 
│   │   ├── config.json     // sample config file
│   │   ├── doc.go          // documentation for config package/module
│   ├── controller          // Our API core handlers 
│   │   ├── article.go      // APIs for Article Handlers
│   │   ├── controller.go   // Common response functions and loading for all handlers
|   ├── model               // Models for our application
│   │   ├── article.go      // Article Model
│   │   ├── category.go     // Category Model
│   │   ├── publisher.go    // Publisher Model
│   ├── .gitignore
│   ├── go.mod          // Dependenies 
│   ├── main.go         // entry point
└────- README.md

```

## API

#### /article
* `GET` : Get all articles
* `POST` : Create a new article

#### /article/:id
* `GET` : Get a article by id
* `PUT` : Update a article id 
* `DELETE` : Delete a article id
