package server

import (
	"context"
	"gin_training/internal/model"
	pb "gin_training/internal/proto"
	"gin_training/internal/storage/postgreSQL/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestStorageServer_GetBook(t *testing.T) {
	s := new(mocks.DB)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	b := model.Book{id, "title", "author"}
	s.On("GetBook", idStr).Return(b, nil)

	var tests = []struct {
		name    string
		stor    *mocks.DB
		param   *pb.BookID
		want    *pb.BookObj
		wantErr codes.Code
	}{
		{
			name:  "Get everything good",
			stor:  s,
			param: &pb.BookID{ID: idStr},
			want:  &pb.BookObj{Id: idStr, Title: "title", Author: "author"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStorage(tc.stor)
			got, err := u.GetBook(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_Create(t *testing.T) {
	s := new(mocks.DB)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	bb := model.Book{Title: "title", Author: "author"}
	b := model.Book{id, "title", "author"}
	s.On("Create", bb).Return(b, nil)

	var tests = []struct {
		name    string
		stor    *mocks.DB
		param   *pb.BookObj
		want    *pb.BookObj
		wantErr codes.Code
	}{
		{
			name:  "Create book everything good",
			stor:  s,
			param: &pb.BookObj{Title: "title", Author: "author"},
			want:  &pb.BookObj{Id: idStr, Title: "title", Author: "author"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStorage(tc.stor)
			got, err := u.Create(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_FindAll(t *testing.T) {
	s := new(mocks.DB)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	b := []model.Book{{id, "title", "author"}}
	s.On("FindAll").Return(b, nil)

	all := []*pb.BookObj{
		{Id: "00000000-0000-0000-0000-000000000000", Title: "title", Author: "author"},
	}

	aa := pb.AllBooks{Allbooks: all}

	var tests = []struct {
		name    string
		stor    *mocks.DB
		want    *pb.AllBooks
		wantErr codes.Code
	}{
		{
			name: "Create book everything good",
			stor: s,
			want: &aa,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStorage(tc.stor)
			got, err := u.FindAll(context.Background(), &emptypb.Empty{})
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_UpdateBook(t *testing.T) {
	s := new(mocks.DB)
	idStr := "00000000-0000-0000-0000-000000000000"
	id, _ := uuid.Parse(idStr)
	nb := pb.BookObj{Title: "title", Author: "author"}
	n := pb.NewBook{ID: idStr, Book: &nb}
	b := model.Book{id, "title", "author"}
	s.On("UpdateBook", mock.Anything, mock.Anything).Return(b, nil)

	var tests = []struct {
		name    string
		stor    *mocks.DB
		param   *pb.NewBook
		want    *pb.BookObj
		wantErr codes.Code
	}{
		{
			name:  "Create book everything good",
			stor:  s,
			param: &n,
			want:  &pb.BookObj{Id: idStr, Title: "title", Author: "author"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStorage(tc.stor)
			got, err := u.UpdateBook(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorageServer_DeleteBook(t *testing.T) {
	s := new(mocks.DB)
	idStr := "00000000-0000-0000-0000-000000000000"
	s.On("DeleteBook", idStr).Return(nil)

	var tests = []struct {
		name    string
		stor    *mocks.DB
		param   *pb.BookID
		want    *emptypb.Empty
		wantErr codes.Code
	}{
		{
			name:  "Get everything good",
			stor:  s,
			param: &pb.BookID{ID: idStr},
			want:  &emptypb.Empty{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := NewGRPCStorage(tc.stor)
			got, err := u.DeleteBook(context.Background(), tc.param)
			if err != nil && status.Code(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
