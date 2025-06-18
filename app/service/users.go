package service

import (
	"app/models"
	"strconv"
)

func (s *Service) Self() (*models.User, error) {
	user := new(models.User)

	if err := s.request(user, get, urlUserSelf, nil); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUsers() ([]models.User, error) {
	users := make([]models.User, 0)

	if err := s.request(&users, get, urlUsers, nil); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) NewUser(user *models.User) error {
	id := new(models.Id)

	if err := s.request(id, post, urlUsers, user); err != nil {
		return err
	}

	user.Id = id.Id

	return nil
}

func (s *Service) UpdateUser(user models.User) error {
	if err := s.request(nil, put, urlUsers, user); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteUser(userId int) error {
	if err := s.request(nil, del, urlUsers, nil,
		[2]string{"user_id", strconv.Itoa(userId)},
	); err != nil {
		return err
	}
	return nil
}
