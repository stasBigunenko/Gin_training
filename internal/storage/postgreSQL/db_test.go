package storage

import (
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gin_training/internal/model"
)

func TestPostgresDB_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO books (id, title, author) VALUES ($1, $2, $3)`).
		WithArgs(sqlmock.AnyArg(), "title", "author").
		WillReturnResult(sqlmock.NewResult(0, 1))

	pdb := &PostgresDB{Pdb: db}

	err = pdb.Pdb.Ping()
	if err != nil {
		log.Fatal(err)
	}

	b := model.Book{Title: "title", Author: "author"}

	res, err := pdb.Create(b)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, b.Title, res.Title)
	assert.Equal(t, b.Author, res.Author)
}

func TestPostgresDB_GetBook(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT title, author FROM books WHERE id=$1`).
		WithArgs("00000000-0000-0000-0000-000000000000").
		WillReturnRows(
			mock.
				NewRows([]string{"id", "title", "author"}).
				AddRow("00000000-0000-0000-0000-000000000000", "title", "author"),
		)

	postgreSQL := &PostgresDB{Pdb: db}

	res, err := postgreSQL.GetBook("00000000-0000-0000-0000-000000000000")
	if err != nil {
		errors.New("error in the database")
		return
	}

	uid, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		errors.New("error with parsing uuid")
		return
	}

	exp := model.Book{uid, "title", "author"}

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, exp, res)
}

func TestPostgresDB_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT * FROM books`).
		WillReturnRows(
			mock.
				NewRows([]string{"id", "title", "author"}).
				AddRow(
					"00000000-0000-0000-0000-000000000000", "title", "author"),
		)

	postgreSQL := &PostgresDB{Pdb: db}

	res := postgreSQL.FindAll()

	uid, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		errors.New("error with parsing uuid")
		return
	}

	exp := []model.Book{
		{uid, "title", "author"},
	}

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, exp, res)
}

func TestPostgresDB_UpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	uid, err := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if err != nil {
		errors.New("error with parsing uuid")
		return
	}

	exp := model.Book{uid, "title2", "author2"}

	in := model.UpdateBookInput{Title: "title2", Author: "author2"}

	mock.ExpectQuery(`SELECT title, author FROM books WHERE id=$1`).
		WithArgs("00000000-0000-0000-0000-000000000000").
		WillReturnRows(mock.
			NewRows([]string{"id", "title", "author"}).
			AddRow("00000000-0000-0000-0000-000000000000", "title", "author"),
		)
	mock.ExpectExec(`UPDATE books SET title=$1, author=$2 WHERE id=$3`).
		WithArgs("title2", "author2", "00000000-0000-0000-0000-000000000000").
		WillReturnResult(sqlmock.NewResult(0, 1))

	postgreSQL := &PostgresDB{Pdb: db}

	res, err := postgreSQL.UpdateBook("00000000-0000-0000-0000-000000000000", in)
	if err != nil {
		errors.New("error in the database")
		return
	}

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, exp, res)
}

func TestPostgresDB_DeleteBook(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection: %s", err, mock)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM books where id = $1`).
		WithArgs("00000000-0000-0000-0000-000000000000").
		WillReturnResult(sqlmock.NewResult(0, 1))

	postgreSQL := &PostgresDB{Pdb: db}

	err = postgreSQL.DeleteBook("00000000-0000-0000-0000-000000000000")

	require.NoError(t, err)
}
