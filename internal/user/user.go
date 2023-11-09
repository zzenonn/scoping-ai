package TrainingNeedsUsers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotImplemented = errors.New("this function is not yet implemented")
)

// User representation
type User struct {
	ID           string  `json:"id" firestore:"id"`
	Name         *string `json:"name,omitempty" firestore:"name,omitempty"`
	EmailAddress *string `json:"email_address,omitempty" firestore:"email_address,omitempty"`
	Corporate    bool    `json:"corporate,omitempty" firestore:"corporate,omitempty"`
	Company      *string `json:"company,omitempty" firestore:"company,omitempty"`
}

// Implements the user repository interface design pattern
type UserRepository interface {
	GetUser(ctx context.Context, id string) (User, error)
	GetAllUsers(ctx context.Context, page int, pageSize int) ([]User, error)
	CreateUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, user User) (User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u *UserService) GetUser(ctx context.Context, id string) (User, error) {
	log.Debug("Retrieving user . . .")

	user, err := u.userRepository.GetUser(ctx, id)

	if err != nil {
		log.Error("Failed to retrieve user")
		return User{}, err
	}

	return user, nil
}

func (u *UserService) GetAllUsers(ctx context.Context, page int, pageSize int) ([]User, error) {
	log.Debug("Retrieving all users . . .")

	users, err := u.userRepository.GetAllUsers(ctx, page, pageSize)

	if err != nil {
		log.Error("Failed to retrieve all users")
		return nil, err
	}

	return users, nil
}

func (u *UserService) CreateUser(ctx context.Context, user User) (User, error) {
	log.Debug("Creating new user . . .")
	user.ID = uuid.New().String()

	createdUser, err := u.userRepository.CreateUser(ctx, user)

	if err != nil {
		log.Error("Failed to create user")
		return User{}, err
	}

	return createdUser, nil
}

func (u *UserService) UpdateUser(ctx context.Context, user User) (User, error) {
	log.Debug("Updating user . . .")

	updatedUser, err := u.userRepository.UpdateUser(ctx, user)

	if err != nil {
		log.Error("Failed to update user")
		return User{}, err
	}

	return updatedUser, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id string) error {
	log.Debug("Deleting user . . .")

	err := u.userRepository.DeleteUser(ctx, id)

	if err != nil {
		log.Error("Failed to delete user")
		return err
	}

	return nil
}
