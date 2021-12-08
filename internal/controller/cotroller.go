package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gin_training/internal/model"
	storage "gin_training/internal/storage/postgreSQL"
)

type Controller struct {
	database storage.DB
}

func NewController(db storage.DB) *Controller {
	return &Controller{
		database: db,
	}
}

// Handlers
func (cr *Controller) Routes() *gin.Engine {
	r := gin.Default()
	r.GET("/books", cr.AllBooks)
	r.POST("/create", cr.CreateBook)
	r.GET("/books/:id", cr.FindBook)
	r.PATCH("/books/:id", cr.UpdateBook)
	r.DELETE("/books/:id", cr.DeleteBook)
	return r
}

// GET /books
// Get all books from db
func (cr *Controller) AllBooks(c *gin.Context) {
	var books []model.Book

	books = cr.database.FindAll()

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// POST /create
// Create a book
func (cr *Controller) CreateBook(c *gin.Context) {
	var input model.CreateBookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := model.Book{
		Title:  input.Title,
		Author: input.Author,
	}

	res, err := cr.database.Create(book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// GET /book/:id
// Find the book by id
func (cr *Controller) FindBook(c *gin.Context) {

	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	res, err := cr.database.GetBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// PUT /book/:id
// Update information about the book by id
func (cr *Controller) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var input model.UpdateBookInput
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := cr.database.UpdateBook(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

// DELETE /book/:id
// Delete book from db
func (cr *Controller) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = cr.database.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "book have been deleted"})
}
