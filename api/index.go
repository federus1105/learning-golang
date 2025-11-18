package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var App *gin.Engine
var db *sql.DB

func Handler(w http.ResponseWriter, r *http.Request) {
	if App == nil {
		App = setupApp()
	}
	App.ServeHTTP(w, r)
}

func setupApp() *gin.Engine {
	// Connect DB
	psqlInfo := os.Getenv("DATABASE_URL")
	if psqlInfo == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic("Failed to open database connection: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	fmt.Println("Successfully connected to the database!")

	// Gin router
	r := gin.New()
	r.Use(gin.Recovery())

	// LOAD TEMPLATES
	r.LoadHTMLGlob("templates/*")

	// ROUTES
	r.GET("/", func(c *gin.Context) {
		users, err := getUsersFromDB()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"users": users})
	})

	r.GET("/users", func(c *gin.Context) {
		users, err := getUsersFromDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	return r
}
