package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gin_training/internal/model"
	storage "gin_training/internal/storage/postgreSQL"
	"gin_training/internal/storage/postgreSQL/mocks"
)

func TestController_FindBook(t *testing.T) {
	db := new(mocks.DB)

	uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

	b := model.Book{uid, "title", "author"}

	db.On("GetBook", "00000000-0000-0000-0000-000000000000").Return(b, nil)

	tests := []struct {
		name       string
		db         storage.DB
		method     string
		url        string
		wantStatus int
		exp        model.Book
	}{
		{
			name:       "Everything ok",
			db:         db,
			method:     "GET",
			url:        "/books/00000000-0000-0000-0000-000000000000",
			wantStatus: http.StatusOK,
			exp:        b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := NewController(tc.db)

			rr := httptest.NewRecorder()

			testRouter := h.Routes()

			testRouter.Use(func(c *gin.Context) {
				c.Set("id", "00000000-0000-0000-0000-000000000000")
			})

			req, err := http.NewRequest(tc.method, tc.url, nil)
			assert.NoError(t, err)

			testRouter.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)

			respBody, err := json.Marshal(gin.H{
				"data": b,
			})

			assert.Equal(t, rr.Body.Bytes(), respBody)
		})
	}
}

func TestController_CreateBook(t *testing.T) {
	db := new(mocks.DB)

	uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

	b := model.Book{uid, "title", "author"}

	db.On("Create", mock.Anything).Return(b, nil)

	tests := []struct {
		name       string
		db         storage.DB
		method     string
		url        string
		wantStatus int
		body       string
		exp        model.Book
	}{
		{
			name:       "Everything ok",
			db:         db,
			method:     "POST",
			url:        "/create",
			wantStatus: http.StatusOK,
			body:       `{"Title":"title","Author":"author"}`,
			exp:        b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := NewController(tc.db)

			rr := httptest.NewRecorder()

			testRouter := h.Routes()

			testRouter.Use(func(c *gin.Context) {
				c.BindJSON(model.CreateBookInput{Title: "title", Author: "author"})
			})

			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			assert.NoError(t, err)

			testRouter.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)

			respBody, err := json.Marshal(gin.H{
				"data": b,
			})

			assert.Equal(t, rr.Body.Bytes(), respBody)
		})
	}
}

func TestController_AllBooks(t *testing.T) {
	db := new(mocks.DB)

	uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

	b := []model.Book{{
		uid, "title", "author"},
	}

	db.On("FindAll").Return(b, nil)

	tests := []struct {
		name       string
		db         storage.DB
		method     string
		url        string
		wantStatus int
		exp        []model.Book
	}{
		{
			name:       "Everything ok",
			db:         db,
			method:     "GET",
			url:        "/books",
			wantStatus: http.StatusOK,
			exp:        b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := NewController(tc.db)

			rr := httptest.NewRecorder()

			testRouter := h.Routes()

			req, err := http.NewRequest(tc.method, tc.url, nil)
			assert.NoError(t, err)

			testRouter.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)

			respBody, err := json.Marshal(gin.H{
				"data": b,
			})

			assert.Equal(t, rr.Body.Bytes(), respBody)
		})
	}
}

func TestController_UpdateBook(t *testing.T) {
	db := new(mocks.DB)

	uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

	b := model.Book{uid, "title", "author"}

	db.On("UpdateBook", mock.Anything, mock.Anything).Return(b, nil)

	tests := []struct {
		name       string
		db         storage.DB
		method     string
		url        string
		wantStatus int
		body       string
		exp        model.Book
	}{
		{
			name:       "Everything ok",
			db:         db,
			method:     "PATCH",
			url:        "/books/00000000-0000-0000-0000-000000000000",
			wantStatus: http.StatusOK,
			body:       `{"Title":"title","Author":"author"}`,
			exp:        b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := NewController(tc.db)

			rr := httptest.NewRecorder()

			testRouter := h.Routes()

			testRouter.Use(func(c *gin.Context) {
				c.Set("id", "00000000-0000-0000-0000-000000000000")
				c.BindJSON(model.CreateBookInput{Title: "title", Author: "author"})
			})

			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			assert.NoError(t, err)

			testRouter.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)

			respBody, err := json.Marshal(gin.H{
				"data": b,
			})

			assert.Equal(t, rr.Body.Bytes(), respBody)
		})
	}
}

func TestController_DeleteBook(t *testing.T) {
	db := new(mocks.DB)

	uid, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")

	b := model.Book{uid, "title", "author"}

	db.On("DeleteBook", "00000000-0000-0000-0000-000000000000").Return(nil)

	tests := []struct {
		name       string
		db         storage.DB
		method     string
		url        string
		wantStatus int
		exp        model.Book
	}{
		{
			name:       "Everything ok",
			db:         db,
			method:     "DELETE",
			url:        "/books/00000000-0000-0000-0000-000000000000",
			wantStatus: http.StatusOK,
			exp:        b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			h := NewController(tc.db)

			rr := httptest.NewRecorder()

			testRouter := h.Routes()

			testRouter.Use(func(c *gin.Context) {
				c.Set("id", "00000000-0000-0000-0000-000000000000")
			})

			req, err := http.NewRequest(tc.method, tc.url, nil)
			assert.NoError(t, err)

			testRouter.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantStatus, rr.Code)
		})
	}
}
