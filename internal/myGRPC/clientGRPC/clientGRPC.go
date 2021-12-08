package clientGRPC

import (
	"context"
	"gin_training/internal/model"
	pb "gin_training/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type gRPCClient struct {
	client pb.BookServiceClient
}

func New(store pb.BookServiceClient) gRPCClient {
	return gRPCClient{
		client: store,
	}
}

func (gc gRPCClient) FindAll() []model.Book {
	ap, _ := gc.client.FindAll(context.Background(), &emptypb.Empty{})

	books := []model.Book{}

	for _, val := range ap.Allbooks {
		res, err := uuid.Parse(val.Id)
		if err != nil {
			return nil
		}

		books = append(books, model.Book{
			ID:     res,
			Title:  val.Title,
			Author: val.Author,
		})
	}

	return books
}
func (gc gRPCClient) Create(in model.Book) (model.Book, error) {
	b, err := gc.client.Create(context.Background(), &pb.BookObj{
		Title:  in.Title,
		Author: in.Author,
	})
	if err != nil {
		return model.Book{}, status.Error(codes.Internal, "couldn't create a book")
	}

	idStr, err := uuid.Parse(b.Id)
	if err != nil {
		return model.Book{}, status.Error(codes.Internal, "couldn't parse id")
	}

	return model.Book{
		ID:     idStr,
		Title:  b.Title,
		Author: b.Author,
	}, nil

}

func (gc gRPCClient) GetBook(id string) (model.Book, error) {

	b, err := gc.client.GetBook(context.Background(), &pb.BookID{
		ID: id,
	})
	if err != nil {
		return model.Book{}, status.Error(codes.Internal, "couldn't find a book")
	}

	uid, _ := uuid.Parse(b.Id)

	return model.Book{
		ID:     uid,
		Title:  b.Title,
		Author: b.Author,
	}, nil

}
func (gc gRPCClient) UpdateBook(id string, in model.UpdateBookInput) (model.Book, error) {

	b, err := gc.client.UpdateBook(context.Background(), &pb.NewBook{
		ID:   id,
		Book: &pb.BookObj{Title: in.Title, Author: in.Author},
	})
	if err != nil {
		return model.Book{}, status.Error(codes.Internal, "couldn't find a book")
	}

	uid, _ := uuid.Parse(b.Id)

	return model.Book{
		ID:     uid,
		Title:  b.Title,
		Author: b.Author,
	}, nil
}

func (gc gRPCClient) DeleteBook(id string) error {
	_, err := gc.client.DeleteBook(context.Background(), &pb.BookID{ID: id})
	if err != nil {
		status.Error(codes.Internal, "couldn't find a book")
	}
	return nil
}
