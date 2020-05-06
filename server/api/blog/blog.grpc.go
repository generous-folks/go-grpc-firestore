package blog

import (
	"context"
	"fmt"
	blogpb "github.com/hpierre74/go-grpc-firestore/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (s *ServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {
	blog, err := s.DB.ReadBlog(ctx, req)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	// Cast to ReadBlogRes type
	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id:       blog.Id,
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}
	return response, nil
}

func (s *ServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	_, err := s.DB.CreateBlog(ctx, req)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	// return the blog in a CreateBlogRes type
	return &blogpb.CreateBlogRes{Blog: req.GetBlog()}, nil
}


func (s *ServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	blog, err := s.DB.UpdateBlog(ctx, req)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return &blogpb.UpdateBlogRes{
		Blog: &blogpb.Blog{
			Id:       blog.Id,
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

func (s *ServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	err := s.DB.DeleteBlog(ctx, req)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return &blogpb.DeleteBlogRes{
		Success: true,
	}, nil
}

func (s *ServiceServer) ListBlogs(req *blogpb.ListBlogsReq, stream blogpb.BlogService_ListBlogsServer) error {
	streamServer := StreamServer{
		Handler: func() StreamHandler {
			return func(blog *BlogItem) {
				stream.Send(&blogpb.ListBlogsRes{
					Blog: &blogpb.Blog{
						Id:       blog.Id,
						AuthorId: blog.AuthorID,
						Content:  blog.Content,
						Title:    blog.Title,
					},
				})
			}
		}(),
	}

	err := s.DB.ListBlogs(&streamServer)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return nil
}

