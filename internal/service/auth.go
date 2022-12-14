package service

import (
	"context"
	"sync"

	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/request"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/response"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/pkg/errors"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	hasher      *hasher.Hasher
	jwtService  *JwtService
}

func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
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

	err = s.sessionRepo.DeleteByPointer(ctx, user.CommonId, r.Fingerprint)
	if err != nil {
		return nil, err
	}

	access, refresh, err := s.issueTokens(user)
	if err != nil {
		return nil, err
	}

	err = s.sessionRepo.Create(ctx, user.CommonId, r.Fingerprint, refresh)
	if err != nil {
		return nil, err
	}

	return &response.JwtResponse{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, r *request.Refresh) (*response.JwtResponse, error) {
	session, err := s.sessionRepo.Find(ctx, r.RefreshToken)
	if err != nil {
		return nil, err
	}
	if session.Subject == "" {
		return nil, errors.New("no valid session")
	}

	err = s.sessionRepo.DeleteByPointer(ctx, session.Subject, session.Fingerprint)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserBySubject(ctx, session.Subject)
	if err != nil {
		return nil, err
	}

	access, refresh, err := s.issueTokens(user)
	if err != nil {
		return nil, err
	}

	err = s.sessionRepo.Create(ctx, user.CommonId, session.Fingerprint, refresh)
	if err != nil {
		return nil, err
	}

	return &response.JwtResponse{
		AccessToken:  access.Value,
		RefreshToken: refresh.Value,
	}, err
}

func (s *AuthService) Logout(ctx context.Context, r *request.Logout, token *entity.ParsedToken) error {
	if token == nil {
		return errors.New("access expired")
	}

	var err error
	if r.Fingerprint != "" {
		err = s.sessionRepo.DeleteByPointer(ctx, token.Subject, r.Fingerprint)
		if err != nil {
			return err
		}
	} else {
		s.sessionRepo.DeleteAll(ctx, token.Subject)
	}

	if err == nil {
		err = s.sessionRepo.AddTokenToBlackList(ctx, token)
	}

	return err
}

func (s *AuthService) issueTokens(user *entity.User) (access, refresh *entity.Token, err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		access, err = s.jwtService.IssueAccess(user)
	}()

	go func() {
		defer wg.Done()
		refresh = s.jwtService.IssueRefresh()
	}()
	wg.Wait()
	return
}
