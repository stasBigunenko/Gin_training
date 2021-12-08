package server

import (
	"context"
	"fmt"
	"gin_training/internal/model"
	pb "gin_training/internal/proto"
	storage "gin_training/internal/storage/postgreSQL"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type StorageServer struct {
	pb.UnimplementedBookServiceServer

	Storage storage.DB
}

func NewGRPCStorage(store storage.DB) *StorageServer {
	return &StorageServer{
		Storage: store,
	}
}

func (s *StorageServer) FindAll(_ context.Context, in *emptypb.Empty) (*pb.AllBooks, error) {
	books := s.Storage.FindAll()

	pbBooks := []*pb.BookObj{}

	for _, val := range books {
		res := val.ID.String()
		pbBooks = append(pbBooks, &pb.BookObj{
			Id:     res,
			Title:  val.Title,
			Author: val.Author,
		})
	}

	return &pb.AllBooks{
		Allbooks: pbBooks,
	}, nil
}
func (s *StorageServer) Create(_ context.Context, in *pb.BookObj) (*pb.BookObj, error) {
	b := model.Book{
		Title:  in.Title,
		Author: in.Author,
	}

	book, err := s.Storage.Create(b)
	if err != nil {
		fmt.Printf(error.Error(err))
		return nil, status.Error(codes.Internal, "internal storage problem")
	}

	res := book.ID.String()

	return &pb.BookObj{
		Id:     res,
		Title:  book.Title,
		Author: book.Author,
	}, nil
}
func (s *StorageServer) GetBook(_ context.Context, in *pb.BookID) (*pb.BookObj, error) {
	book, err := s.Storage.GetBook(in.ID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get post")
	}

	res := book.ID.String()

	return &pb.BookObj{
		Id:     res,
		Title:  book.Title,
		Author: book.Author,
	}, nil
}

func (s *StorageServer) UpdateBook(_ context.Context, in *pb.NewBook) (*pb.BookObj, error) {

	book := model.UpdateBookInput{
		Title:  in.Book.Title,
		Author: in.Book.Author,
	}

	res, err := s.Storage.UpdateBook(in.ID, book)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update post")
	}

	resId := res.ID.String()

	return &pb.BookObj{
		Id:     resId,
		Title:  res.Title,
		Author: res.Author,
	}, nil
}
func (s *StorageServer) DeleteBook(_ context.Context, in *pb.BookID) (*emptypb.Empty, error) {

	err := s.Storage.DeleteBook(in.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to delete post")
	}
	return &emptypb.Empty{}, nil
}
