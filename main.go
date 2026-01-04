package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Trust me, I write production code now B)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if err := db.PingContext(context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"db": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "up"})
	})

	r.GET("/items", getItems)
	r.GET("/items/:id", getItem)
	r.POST("/items", createItem)

	log.Println("Server listening on :8080")
	r.Run(":8080")
}

func getItems(c *gin.Context) {
	rows, err := db.QueryContext(context.Background(), `
		SELECT id, category_id, name, description, image_url,
		       quantity, owner, currently_borrowing, created_at
		FROM items
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var i Item
		err := rows.Scan(
			&i.ID,
			&i.CategoryID,
			&i.Name,
			&i.Description,
			&i.ImageURL,
			&i.Quantity,
			&i.Owner,
			&i.CurrentlyBorrowing,
			&i.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, i)
	}

	c.JSON(http.StatusOK, items)
}

func getItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var i Item
	err = db.QueryRowContext(context.Background(), `
		SELECT id, category_id, name, description, image_url,
		       quantity, owner, currently_borrowing, created_at
		FROM items
		WHERE id = $1
	`, id).Scan(
		&i.ID,
		&i.CategoryID,
		&i.Name,
		&i.Description,
		&i.ImageURL,
		&i.Quantity,
		&i.Owner,
		&i.CurrentlyBorrowing,
		&i.CreatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, i)
}

func createItem(c *gin.Context) {
	var i Item
	if err := c.BindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRowContext(context.Background(), `
		INSERT INTO items
		    (category_id, name, description, image_url, quantity, owner, currently_borrowing)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`,
		i.CategoryID,
		i.Name,
		i.Description,
		i.ImageURL,
		i.Quantity,
		i.Owner,
		i.CurrentlyBorrowing,
	).Scan(&i.ID, &i.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, i)
}
