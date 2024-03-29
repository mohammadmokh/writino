package contract

import (
	"context"

	"github.com/mohammadmokh/writino/dto"
)

type PostInteractor interface {
	CreatePost(context.Context, dto.CreatePostReq) (dto.CreatePostRes, error)
	UpdatePost(context.Context, dto.UpdatePostReq) error
	FindPostByID(context.Context, dto.FindPostByIDReq) (dto.FindPostRes, error)
	DeletePost(context.Context, dto.DeletePostReq) error
	SearchPost(context.Context, dto.SearchPostReq) (dto.SearchPostRes, error)
	FindUsersPosts(context.Context, dto.FindUsersPostsReq) (dto.SearchPostRes, error)
	FindAll(context.Context, dto.SearchPostReq) (dto.SearchPostRes, error)
	LikePost(context.Context, dto.LikePostReq) error
	DeleteUserPosts(context.Context, dto.DeleteUserPostsReq) error
}
