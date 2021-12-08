package clientGRPC

import (
	"gin_training/internal/model"
	mocks2 "gin_training/internal/myGRPC/clientGRPC/mocks"
	Gin_training "gin_training/internal/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestGRPCClient_GetBook(t *testing.T) {
	s := new(mocks2.BookServiceClient)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	b := model.Book{id, "title", "author"}
	bb := Gin_training.BookObj{Id: "00000000-0000-0000-0000-000000000000", Title: "title", Author: "author"}
	s.On("GetBook", mock.Anything, mock.Anything).Return(&bb, nil)

	var tests = []struct {
		name  string
		stor  *mocks2.BookServiceClient
		param string
		want  model.Book
	}{
		{
			name:  "Get everything good",
			stor:  s,
			param: idStr,
			want:  b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := New(s)
			got, err := u.GetBook(tc.param)
			if err != nil {
				t.Errorf("error = %v", err.Error())
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGRPCClient_Create(t *testing.T) {
	s := new(mocks2.BookServiceClient)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	b := model.Book{id, "title", "author"}
	bb := Gin_training.BookObj{Id: "00000000-0000-0000-0000-000000000000", Title: "title", Author: "author"}
	s.On("Create", mock.Anything, mock.Anything).Return(&bb, nil)

	var tests = []struct {
		name  string
		stor  *mocks2.BookServiceClient
		param model.Book
		want  model.Book
	}{
		{
			name:  "Get everything good",
			stor:  s,
			param: model.Book{Title: "title", Author: "author"},
			want:  b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := New(s)
			got, err := u.Create(tc.param)
			if err != nil {
				t.Errorf("error = %v", err.Error())
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGRPCClient_FindAll(t *testing.T) {
	s := new(mocks2.BookServiceClient)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	b := []model.Book{{id, "title", "author"}}
	all := []*Gin_training.BookObj{
		{Id: "00000000-0000-0000-0000-000000000000", Title: "title", Author: "author"},
	}
	aa := Gin_training.AllBooks{Allbooks: all}

	s.On("FindAll", mock.Anything, mock.Anything).Return(&aa, nil)

	var tests = []struct {
		name string
		stor *mocks2.BookServiceClient
		want []model.Book
	}{
		{
			name: "Get everything good",
			stor: s,
			want: b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := New(s)
			got := u.FindAll()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGRPCClient_UpdateBook(t *testing.T) {
	s := new(mocks2.BookServiceClient)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	b := model.Book{id, "title", "author"}
	bb := Gin_training.BookObj{Id: "00000000-0000-0000-0000-000000000000", Title: "title", Author: "author"}
	s.On("UpdateBook", mock.Anything, mock.Anything).Return(&bb, nil)

	var tests = []struct {
		name   string
		stor   *mocks2.BookServiceClient
		param1 string
		param2 model.UpdateBookInput
		want   model.Book
	}{
		{
			name:   "Get everything good",
			stor:   s,
			param1: "00000000-0000-0000-0000-000000000000",
			param2: model.UpdateBookInput{"title", "author"},
			want:   b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := New(s)
			got, err := u.UpdateBook(tc.param1, tc.param2)
			if err != nil {
				t.Errorf("error = %v", err.Error())
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGRPCClient_DeleteBook(t *testing.T) {
	s := new(mocks2.BookServiceClient)
	idStr := "00000000-0000-0000-0000-000000000000"
	//bb := Gin_training.BookID{ID: "00000000-0000-0000-0000-000000000000"}
	s.On("DeleteBook", mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

	var tests = []struct {
		name  string
		stor  *mocks2.BookServiceClient
		param string
		want  error
	}{
		{
			name:  "Get everything good",
			stor:  s,
			param: idStr,
			want:  nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := New(s)
			err := u.DeleteBook(tc.param)
			assert.NoError(t, err)
		})
	}
}
