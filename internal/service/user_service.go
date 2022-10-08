package service

import (
	"context"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/dto"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
)

type UserService struct {
	repo   *repository.UserRepository
	hasher *hasher.Hasher
}

func NewUserService(repo *repository.UserRepository, hasher *hasher.Hasher) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
	}
}

func (us *UserService) AuthByEmail(ctx context.Context, requestDto *dto.RequestDto) (*entity.User, error) {
	passwordHash := us.hasher.GetMD5Hash(requestDto.Pass)
	return us.repo.FindUser(ctx, requestDto.Email, passwordHash)
}
