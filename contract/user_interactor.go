package contract

import (
	"context"

	"github.com/mohammadmokh/writino/dto"
)

type UserInteractor interface {
	CheckEmail(context.Context, dto.CheckEmailReq) (dto.CheckEmailRes, error)
	Register(context.Context, dto.RegisterReq) error
	Update(context.Context, dto.UpdateUserReq) error
	UpdatePassword(context.Context, dto.UpdatePasswordReq) error
	Find(context.Context, dto.FindUserReq) (dto.FindUserRes, error)
	DeleteAccount(context.Context, dto.DeleteUserReq) error
	VerifyUser(context.Context, dto.VerifyUserReq) error
	UpdateProfilePic(context.Context, dto.UpdateProfilePicReq) (dto.UpdateProfilePicRes, error)
}
