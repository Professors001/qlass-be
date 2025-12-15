package usecases

import (
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/usecases/repositories"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Register(user *entities.User, password string) error
	GetUserByID(id uint) (*entities.User, error)
	GetUserByUUID(uuid string) (*entities.User, error)
}

type userUseCase struct {
	userRepo repositories.UserRepository 
}

func NewUserUseCase(repo repositories.UserRepository) UserUseCase {
	return &userUseCase{userRepo: repo}
}

func (u *userUseCase) Register(user *entities.User, password string) error {
	// 1. Check if email exists (Optional, repo will error anyway, but good for UI)
	if existingUser, _ := u.userRepo.GetByEmail(user.Email); existingUser != nil {
		return errors.New("email already in use")
	}

	// 2. Hash Password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPwd)

	// 3. Generate UUID and Metadata
	user.UUID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	user.Role = "student" // Default role

	// 4. Save to DB
	return u.userRepo.Create(user)
}

func (u *userUseCase) GetUserByID(id uint) (*entities.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userUseCase) GetUserByUUID(uuid string) (*entities.User, error) {
	return u.userRepo.GetByUUID(uuid)
}