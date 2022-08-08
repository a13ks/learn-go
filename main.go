package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// album represents data about a record album.
type album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

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
	rows, err := h.db.Query("SELECT * FROM album")
	CheckError(err)

	defer rows.Close()

	var albums = []album{}
	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			panic(err)
		}

		albums = append(albums, alb)
	}

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

	// Add the new album to the slice.
	// albums = append(albums, newAlbum)

	lastInsertId := int64(0)
	err := h.db.QueryRow("INSERT into album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&lastInsertId)

	if err != nil {
		log.Fatalf("An error occurred while executing query: %v", err)
	}

	newAlbum.ID = lastInsertId

	c.Header("Location", fmt.Sprintf("/albums/%d", lastInsertId))
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func (h Handler) getAlbumByID(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.ParseInt(param, 10, 64)

	CheckError(err)

	rows, err := h.db.Query("SELECT * FROM album WHERE id = $1", id)
	CheckError(err)

	defer rows.Close()

	if rows == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			panic(err)
		}

		c.IndentedJSON(http.StatusOK, alb)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	return
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
