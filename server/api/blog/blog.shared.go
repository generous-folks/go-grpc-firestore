package blog

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceServer struct {
	DB *DB
}

type DB struct {
	Collection *firestore.CollectionRef
}

type StreamHandler func(*BlogItem)
type StreamServer struct {
	Handler StreamHandler
}

type BlogItem struct {
	Id       string `firestore:"Id"`
	AuthorID string `firestore:"AuthorID"`
	Content  string `firestore:"Content"`
	Title    string `firestore:"Title"`
}

func handleInternalError(err error) error {
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	return nil
}