package blog

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	blogpb "github.com/hpierre74/go-grpc-firestore/proto"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



// ReadBlog gets a BlogItem from the database
func (db *DB) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*BlogItem, error) {
	var data map[string]interface{}

	id := req.GetId()

	snapshot, err := db.Collection.Doc(id).Get(ctx)
	handleInternalError(err)

	// Assign the result to the BlogItem struct
	snapshot.DataTo(&data)
	blog := BlogItem{}
	mapstructure.Decode(&data, &blog)

	return &blog, nil
}

func (db *DB) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*BlogItem, error) {
	uuid, err := uuid.NewRandom()
	handleInternalError(err)

	id := uuid.String()

	blog := req.GetBlog()
	data := BlogItem{
		Id:       id,
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	_, err = db.Collection.Doc(id).Create(ctx, data)
	handleInternalError(err)


	// return the blog in a CreateBlogRes type
	return &data, nil
}

func (db *DB) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*BlogItem, error) {
	blog := req.GetBlog()

	id := blog.GetId()

	update := BlogItem{
		Id:       id,
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	_, err := db.Collection.Doc(id).Set(ctx, update)
	handleInternalError(err)


	return &update, nil
}

func (db *DB) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) error {
	id := req.GetId()

	_, err := db.Collection.Doc(id).Delete(ctx)
	if err != nil {
		return status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete blog with id %s: %v", req.GetId(), err))
	}

	return nil
}


func (db *DB) ListBlogs(stream *StreamServer) error {
	var data map[string]interface{}
	ctx := context.Background()
	iter := db.Collection.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return status.Errorf(codes.NotFound, fmt.Sprintf("Could not iterate on documents:  %v", err))
		}

		// Assign the result to the BlogItem struct
		doc.DataTo(&data)
		blog := BlogItem{}
		mapstructure.Decode(&data, &blog)

		stream.Handler(&blog)
	}

	return nil
}
