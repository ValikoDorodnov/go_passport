package service

import (
	"context"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/dto"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
)

type UserService struct {
	repo       *repository.UserRepository
	hasher     *hasher.Hasher
	jwtService *JwtService
}

func NewUserService(repo *repository.UserRepository, hasher *hasher.Hasher, jwt *JwtService) *UserService {
	return &UserService{
		repo:       repo,
		hasher:     hasher,
		jwtService: jwt,
	}
}

func (us *UserService) LoginByEmail(ctx context.Context, requestDto *dto.LoginByEmailDto) (*entity.Jwt, error) {
	passwordHash := us.hasher.GetMD5Hash(requestDto.Pass)
	user, err := us.repo.FindUser(ctx, requestDto.Email, passwordHash)

	if err != nil {
		return nil, err
	}

	access := us.jwtService.IssueToken(user, us.jwtService.config.AccessTtl, "access")
	refresh := us.jwtService.IssueToken(user, us.jwtService.config.RefreshTtl, "refresh")

	return &entity.Jwt{AccessToken: access, RefreshToken: refresh}, err
}

func (us UserService) LoginByRefresh(requestDto *dto.LoginByRefreshDto) (*entity.Jwt, error) {
	user, err := us.jwtService.ParseToken(requestDto.RefreshToken)

	if err != nil {
		return nil, err
	}

	access := us.jwtService.IssueToken(user, us.jwtService.config.AccessTtl, "access")
	refresh := us.jwtService.IssueToken(user, us.jwtService.config.RefreshTtl, "refresh")

	return &entity.Jwt{AccessToken: access, RefreshToken: refresh}, err
}
