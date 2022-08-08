package main

import (
	"database/sql"
	"log"
)

func getAlbums(db *sql.DB) []album {
	rows, err := db.Query("SELECT * FROM album")
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
	return albums
}

func createAlbum(db *sql.DB, newAlbum *album) int64 {
	lastInsertId := int64(0)
	err := db.QueryRow("INSERT into album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		newAlbum.Title, newAlbum.Artist, newAlbum.Price).Scan(&lastInsertId)

	if err != nil {
		log.Fatalf("An error occurred while executing query: %v", err)
	}

	return lastInsertId
}

func getAlbumByID(db *sql.DB, id int64) album {
	rows, err := db.Query("SELECT * FROM album WHERE id = $1", id)
	CheckError(err)

	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			panic(err)
		}

		return alb
	}
	return album{}
}
