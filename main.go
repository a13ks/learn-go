package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Database(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

type Handler struct {
	db *sql.DB
}

// getAlbums responds with the list of all albums as JSON.
func (h Handler) getAlbums(c *gin.Context) {
	var albums = getAlbums(h.db)

	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func (h Handler) postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	lastInsertId := createAlbum(h.db, &newAlbum)

	c.Header("Location", fmt.Sprintf("/albums/%d", lastInsertId))
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func (h Handler) getAlbumByID(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.ParseInt(param, 10, 64)
	CheckError(err)

	if alb := getAlbumByID(h.db, id); alb.ID > 0 {
		c.IndentedJSON(http.StatusOK, alb)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var psqlconn = os.Getenv("DATABASE_URL")
	if psqlconn == "" {
		psqlconn = "postgres://user:password@localhost:5432/go_tips?sslmode=disable"
	}

	println("psqlconn = ", psqlconn)
	db, err := sql.Open("postgres", psqlconn)

	h := Handler{db}

	CheckError(err)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// r.Use(Database(db))
	r.GET("/albums", h.getAlbums)
	r.GET("/albums/:id", h.getAlbumByID)
	r.POST("/albums", h.postAlbums)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
