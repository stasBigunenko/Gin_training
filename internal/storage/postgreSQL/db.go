package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"gin_training/internal/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

type PostgresDB struct {
	Pdb *sql.DB
	mu  sync.Mutex
}

func NewPDB(host string, port string, user string, psw string, dbname string, ssl string) (*PostgresDB, error) {
	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + psw + " dbname=" + dbname + " sslmode=" + ssl

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %w\n", err)
	}

	database := &PostgresDB{Pdb: db}

	database.Pdb.Exec("CREATE TABLE IF NOT EXISTS books (\n    id VARCHAR(40) PRIMARY KEY NOT NULL,\n    title VARCHAR(150) NOT NULL,\n    author VARCHAR(50) NOT NULL\n);")

	return database, nil
}

func (pdb *PostgresDB) FindAll() []model.Book {

	rows, err := pdb.Pdb.Query(
		`SELECT * FROM books`)
	if err != nil {
		return []model.Book{}
	}
	defer rows.Close()

	books := []model.Book{}

	for rows.Next() {
		b := model.Book{}
		var bb string
		err := rows.Scan(&bb, &b.Title, &b.Author)
		if err != nil {
			return []model.Book{}
		}
		b.ID, err = uuid.Parse(bb)
		if err != nil {
			return []model.Book{}
		}
		books = append(books, b)
	}

	return books
}

func (pdb *PostgresDB) Create(b model.Book) (model.Book, error) {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	id := uuid.New()
	idStr := id.String()

	log.Println(b)

	_, err := pdb.Pdb.Exec(
		"INSERT INTO books (id, title, author) VALUES ($1, $2, $3)", idStr, b.Title, b.Author)
	if err != nil {
		return model.Book{}, errors.New("couldn't create book in database")
	}

	b.ID = id

	return b, nil
}

func (pdb *PostgresDB) GetBook(id string) (model.Book, error) {

	var b model.Book

	err := pdb.Pdb.QueryRow(
		`SELECT title, author FROM books WHERE id=$1`, id).Scan(&b.Title, &b.Author)
	if err != nil {
		return model.Book{}, errors.New("couldn't find book")
	}

	b.ID, _ = uuid.Parse(id)

	return b, nil
}

func (pdb *PostgresDB) UpdateBook(id string, in model.UpdateBookInput) (model.Book, error) {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	var b model.Book

	err := pdb.Pdb.QueryRow(
		`SELECT title, author FROM books WHERE id=$1`, id).Scan(&b.Title, &b.Author)
	if err != nil {
		return model.Book{}, errors.New("couldn't find book")
	}

	if in.Title == "" {
		in.Title = b.Title
	}
	if in.Author == "" {
		in.Author = b.Author
	}

	_, err = pdb.Pdb.Exec(
		`UPDATE books SET title=$1, author=$2 WHERE id=$3`, in.Title, in.Author, id)
	if err != nil {
		return model.Book{}, errors.New("couldn't update post")
	}

	UID, err := uuid.Parse(id)

	b = model.Book{
		ID:     UID,
		Title:  in.Title,
		Author: in.Author,
	}

	return b, nil
}

func (pdb *PostgresDB) DeleteBook(id string) error {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	_, err := pdb.Pdb.Exec(
		`DELETE FROM books where id = $1`, id)
	if err != nil {
		return errors.New("couldn't delete book")
	}

	return nil
}
