package main

import "time"

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Item struct {
	ID                 int       `json:"id"`
	CategoryID         int       `json:"category_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ImageURL           string    `json:"image_url"`
	Quantity           int       `json:"quantity"`
	Owner              string    `json:"owner"`
	CurrentlyBorrowing string    `json:"currently_borrowing"`
	CreatedAt          time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	ItemID    int       `json:"item_id"`
	Author    string    `json:"author"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
