package main

// album represents data about a record album.
type album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
