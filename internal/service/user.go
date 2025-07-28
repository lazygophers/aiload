package service

import (
	"aiload/internal/model"
	"aiload/internal/repository"
	"encoding/json"
	"strconv"

	"go.etcd.io/bbolt"
)

type UserService struct {
	repo  *repository.UserRepository
	cache *bbolt.DB
}

func NewUserService(repo *repository.UserRepository, cache *bbolt.DB) *UserService {
	return &UserService{repo: repo, cache: cache}
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUser(id uint) (*model.User, error) {
	var user *model.User
	var err error

	// Try to get user from cache first
	err = s.cache.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		val := b.Get([]byte(strconv.Itoa(int(id))))
		if val != nil {
			return json.Unmarshal(val, &user)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// If not in cache, get from DB
	user, err = s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Put user into cache
	err = s.cache.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		val, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(strconv.Itoa(int(id))), val)
	})

	return user, err
}

func (s *UserService) UpdateUser(user *model.User) error {
	err := s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	// Invalidate cache
	return s.cache.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.Delete([]byte(strconv.Itoa(int(user.ID))))
	})
}

func (s *UserService) DeleteUser(id uint) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	return s.cache.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.Delete([]byte(strconv.Itoa(int(id))))
	})
}
