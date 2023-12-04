package service

import (
	"main/03-mockDemo/consts"
	"main/03-mockDemo/repo"
	"main/03-mockDemo/service/dto"

	"github.com/pkg/errors"
)

type UserService struct {
	userRepo repo.UserRepoInterface
	userDto  dto.UserDto
}

func NewUserService(userRepo repo.UserRepoInterface, userDto dto.UserDto) *UserService {
	return &UserService{userRepo: userRepo, userDto: userDto}
}

func (s *UserService) AddUser() (bool, error) {
	user, err := s.userRepo.SelectOne(s.userDto.Id)
	if err != nil {
		return false, errors.Wrap(err, consts.SystemErrorPrefix)
	}

	// user exist.
	if len(user.Id) != 0 {
		return false, errors.New(consts.UserExist)
	}

	// do create.

	return true, nil
}
