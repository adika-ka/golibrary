package service

import (
	"context"
	"fmt"
	"golibrary/internal/entities"
	"golibrary/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user entities.User) (entities.User, error)
	FindUserByID(ctx context.Context, id int) (entities.User, error)
	ListUsers(ctx context.Context) ([]entities.User, error)
	UpdateUser(ctx context.Context, id int, user entities.User) (entities.User, error)
	DeleteUser(ctx context.Context, id int) error
	FindListUserWithBooks(ctx context.Context) ([]entities.User, error)
	ExistsUserByID(ctx context.Context, id int) (bool, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) CreateUser(ctx context.Context, user entities.User) (entities.User, error) {
	if err := ValidUser(ctx, user); err != nil {
		return user, err
	}

	return u.repo.Create(ctx, user)
}

func (u *userService) FindUserByID(ctx context.Context, id int) (entities.User, error) {
	if err := ValidID(ctx, id); err != nil {
		return entities.User{}, err
	}

	return u.repo.GetByID(ctx, id)
}

func (u *userService) ListUsers(ctx context.Context) ([]entities.User, error) {
	return u.repo.GetAllWithBooks(ctx)
}

func (u *userService) UpdateUser(ctx context.Context, id int, user entities.User) (entities.User, error) {
	if err := ValidID(ctx, id); err != nil {
		return user, err
	}

	if err := ValidUser(ctx, user); err != nil {
		return user, err
	}

	return u.repo.Update(ctx, id, user)
}

func (u *userService) DeleteUser(ctx context.Context, id int) error {
	if err := ValidID(ctx, id); err != nil {
		return err
	}

	return u.repo.Delete(ctx, id)
}

func (u *userService) FindListUserWithBooks(ctx context.Context) ([]entities.User, error) {
	return u.repo.GetAllWithBooks(ctx)
}

func (u *userService) ExistsUserByID(ctx context.Context, id int) (bool, error) {
	if err := ValidID(ctx, id); err != nil {
		return false, err
	}

	return u.repo.ExistsByID(ctx, id)
}

func ValidID(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("incorrect id %d", id)
	}

	return nil
}

func ValidUser(ctx context.Context, user entities.User) error {
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if user.Name == "" {
		return fmt.Errorf("user name cannot be empty")
	}

	return nil
}
