package service

import (
	"context"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/request"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/response"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.RefreshSessionRepository
	hasher      *hasher.Hasher
	jwtService  *JwtService
}

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.RefreshSessionRepository,
	hasher *hasher.Hasher,
	jwtService *JwtService,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		hasher:      hasher,
		jwtService:  jwtService,
	}
}

func (s *AuthService) SignIn(ctx context.Context, r *request.LoginByEmail) (*response.JwtResponse, error) {

	passwordHash := s.hasher.GetMD5Hash(r.Pass)
	user, err := s.userRepo.FindUserByCredentials(ctx, r.Email, passwordHash)
	if err != nil {
		return nil, err
	}

	err = s.sessionRepo.Delete(ctx, user.CommonId, r.Platform)
	if err != nil {
		return nil, err
	}

	access, refresh := s.issueTokens(user)
	err = s.sessionRepo.Create(ctx, user.CommonId, r.Platform, refresh)

	if err != nil {
		return nil, err
	}

	return &response.JwtResponse{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
	}, err
}

func (s AuthService) RefreshTokens(ctx context.Context, r *request.Refresh) (*response.JwtResponse, error) {
	session, err := s.sessionRepo.FindByRefresh(ctx, r.RefreshToken)

	if err != nil {
		return nil, err
	}
	err = s.sessionRepo.Delete(ctx, session.Subject, session.Platform)
	if err != nil {
		return nil, err
	}

	if session.ExpiresIn <= time.Now().Unix() {
		return nil, errors.Wrap(err, "session expired")
	}

	user, err := s.userRepo.FindUserById(ctx, session.Subject)
	if err != nil {
		return nil, err
	}

	access, refresh := s.issueTokens(user)
	err = s.sessionRepo.Create(ctx, user.CommonId, session.Platform, refresh)

	if err != nil {
		return nil, err
	}

	return &response.JwtResponse{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
	}, err
}

func (s AuthService) issueTokens(user *entity.User) (access, refresh *entity.Token) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		access = s.jwtService.IssueAccess(user)
	}()

	go func() {
		defer wg.Done()
		refresh = s.jwtService.IssueRefresh()
	}()
	wg.Wait()
	return
}
