package repositories

import "qlass-be/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByEmail(email string) (*entities.User, error)
	GetByID(id uint) (*entities.User, error)
	GetByUniID(uuid string) (*entities.User, error)
}