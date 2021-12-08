package storage

import "gin_training/internal/model"

type DB interface {
	FindAll() []model.Book
	Create(model.Book) (model.Book, error)
	GetBook(string) (model.Book, error)
	UpdateBook(string, model.UpdateBookInput) (model.Book, error)
	DeleteBook(string) error
}
