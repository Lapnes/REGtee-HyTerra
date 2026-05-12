package auth

import (
	authErrors "backend/constant"
	httpAuth "backend/entity/http/auth"
	"backend/repository/user"
	hash "backend/util/cryptography"
	"backend/util/jwt/jws"
	"fmt"
)

type Service struct {
	repoUser user.Repository
}

func NewService(repoUser user.Repository) *Service {
	return &Service{repoUser: repoUser}
}

type Servicer interface {
	Login(login httpAuth.Login) (resp httpAuth.LoginResponse, err error)
}

func (svc Service) Login(login httpAuth.Login) (resp httpAuth.LoginResponse, err error) {
	userEntity, err := svc.repoUser.GetUserByName(login.Nama)
	if err != nil {
		err = authErrors.ErrUserNotFound
		return
	}

	if userEntity.Password != hash.HashPassword(login.Password) {
		err = authErrors.ErrWrongCredentials
		return
	}

	userClaim := jws.UserClaims{
		ID: int(userEntity.ID),
	}

	resp.AccessToken, _, err = userClaim.NewToken(true)
	if err != nil {
		err = fmt.Errorf("failed to generate token: %w", err)
		return
	}
	return
}
