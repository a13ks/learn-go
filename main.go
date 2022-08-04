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

// albums slice to seed record album data.
var albums = []album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func Database(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	// albums = append(albums, newAlbum)

	lastInsertId := int64(0)
	db := c.MustGet("db").(*sql.DB)
	err := db.QueryRow("INSERT into album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&lastInsertId)

	if err != nil {
		log.Fatalf("An error occurred while executing query: %v", err)
	}

	newAlbum.ID = lastInsertId

	c.Header("Location", fmt.Sprintf("/albums/%d", lastInsertId))
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID_(c *gin.Context) {
	param := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		panic(err)
	}

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getAlbumByID(c *gin.Context) {
	param := c.Param("id")

	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		panic(err)
	}

	db := c.MustGet("db").(*sql.DB)
	rows, err := db.Query("SELECT * FROM album WHERE id = $1", id)
	if err != nil {
		panic(err)
	}

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

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Use(Database(db))
	r.GET("/albums", getAlbums)
	r.GET("/albums/:id", getAlbumByID)
	r.POST("/albums", postAlbums)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
