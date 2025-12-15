package repository

import (
	"errors"
	"qlass-be/internal/domain"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new instance of the repository
func NewPostgresUserRepository(db *gorm.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *postgresUserRepository) Create(user *domain.User) error {
	// GORM handles the SQL Insert automatically
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByEmail finds a user by their email address
func (r *postgresUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	
	// SELECT * FROM users WHERE email = ? LIMIT 1
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

// GetByID finds a user by their primary key ID
func (r *postgresUserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	
	// SELECT * FROM users WHERE id = ?
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}