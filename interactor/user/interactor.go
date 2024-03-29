package user

import (
	"context"
	"time"

	"github.com/mohammadmokh/writino/contract"
	"github.com/mohammadmokh/writino/dto"
	"github.com/mohammadmokh/writino/entity"
	"github.com/mohammadmokh/writino/errors"
	"golang.org/x/crypto/bcrypt"
)

type PostInteractor interface {
	DeleteUserPosts(context.Context, dto.DeleteUserPostsReq) error
}

type CommentInteractor interface {
	DeleteUserComments(context.Context, dto.DeleteUserCommentsReq) error
}

type UserIntractor struct {
	store             contract.UserStore
	mail              contract.EmailService
	profilePic        contract.ProfilePicStore
	verificationCode  contract.VerificationCodeInteractor
	postInteractor    PostInteractor
	commentInteractor CommentInteractor
}

func (i UserIntractor) Register(ctx context.Context, req dto.RegisterReq) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := entity.User{

		Email:       req.Email,
		Password:    string(hashedPassword),
		DisplayName: req.Email,
		IsVerified:  false,
	}

	body, err := i.verificationCode.Create(ctx, user.Email)
	if err != nil {
		return err
	}

	err = i.store.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	err = i.mail.SendEmail(user.Email, "Verification Code", body)
	return err

}

func (i UserIntractor) CheckEmail(ctx context.Context, req dto.CheckEmailReq) (dto.CheckEmailRes, error) {

	user, err := i.store.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if err == errors.ErrNotFound {
			return dto.CheckEmailRes{IsUnique: true}, nil
		}
		return dto.CheckEmailRes{}, err
	}

	// if user registered but not verified we delete the user and free email Address
	if !user.IsVerified && (time.Since(user.CreatedAt).Minutes() > 5) {
		err := i.store.DeleteUser(ctx, user.Id)
		if err != nil {
			return dto.CheckEmailRes{}, err
		}
		return dto.CheckEmailRes{IsUnique: true}, nil
	}

	return dto.CheckEmailRes{IsUnique: false}, nil
}

func (i UserIntractor) Update(ctx context.Context, req dto.UpdateUserReq) error {

	user, err := i.store.FindUser(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}

	if req.ProfilePic != nil {
		user.ProfilePic = *req.ProfilePic
	}

	err = i.store.UpdateUser(ctx, user)
	return err
}

func (i UserIntractor) DeleteAccount(ctx context.Context, req dto.DeleteUserReq) error {

	err := i.postInteractor.DeleteUserPosts(ctx, dto.DeleteUserPostsReq{
		UserID: req.Id,
	})
	if err != nil {
		return err
	}

	err = i.commentInteractor.DeleteUserComments(ctx, dto.DeleteUserCommentsReq{
		UserID: req.Id,
	})
	if err != nil {
		return err
	}

	err = i.store.DeleteUser(ctx, req.Id)
	return err
}

func (i UserIntractor) Find(ctx context.Context, req dto.FindUserReq) (dto.FindUserRes, error) {

	user, err := i.store.FindUser(ctx, req.Id)
	if err != nil {
		return dto.FindUserRes{}, err
	}
	return dto.FindUserRes{
		ProfilePic:  user.ProfilePic,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		Email:       user.Email,
	}, nil
}

func (i UserIntractor) UpdatePassword(ctx context.Context, req dto.UpdatePasswordReq) error {

	user, err := i.store.FindUser(ctx, req.ID)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Old)) != nil {
		return errors.ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.New), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = i.store.UpdateUser(ctx, user)
	return err
}

func (i UserIntractor) VerifyUser(ctx context.Context, req dto.VerifyUserReq) error {

	code, err := i.verificationCode.Find(ctx, req.Email)
	if err != nil {
		return err
	}
	if code != req.VerificationCode {
		return errors.ErrInvalidCredentials
	}

	user, err := i.store.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	user.IsVerified = true

	err = i.store.UpdateUser(ctx, user)
	return err
}

func (i UserIntractor) UpdateProfilePic(ctx context.Context, req dto.UpdateProfilePicReq) (
	dto.UpdateProfilePicRes, error) {

	user, err := i.store.FindUser(ctx, req.ID)
	if err != nil {
		return dto.UpdateProfilePicRes{}, err
	}

	filename := user.Id + "." + req.Format
	err = i.profilePic.SaveImage(req.Image, filename)
	if err != nil {
		return dto.UpdateProfilePicRes{}, err
	}
	user.ProfilePic = filename

	err = i.store.UpdateUser(ctx, user)
	return dto.UpdateProfilePicRes{
		Link: filename,
	}, err
}
