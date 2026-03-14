package databases

import (
	"errors"
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new instance of the repository
func NewPostgresUserRepository(db *gorm.DB) repositories.UserRepository {
	return &postgresUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *postgresUserRepository) Create(user *entities.User) error {
	// GORM handles the SQL Insert automatically
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetByEmail finds a user by their email address
func (r *postgresUserRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User

	if err := r.db.Preload("ProfileImgAttachment").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByID finds a user by their primary key ID
func (r *postgresUserRepository) GetByID(id uint) (*entities.User, error) {
	var user entities.User

	if err := r.db.Preload("ProfileImgAttachment").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByUniID finds a user by their university ID
func (r *postgresUserRepository) GetByUniID(universityID string) (*entities.User, error) {
	var user entities.User

	if err := r.db.Preload("ProfileImgAttachment").Where("university_id = ?", universityID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAll fetches all users from the database
func (r *postgresUserRepository) GetAll() ([]*entities.User, error) {
	var users []*entities.User

	if err := r.db.Preload("ProfileImgAttachment").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *postgresUserRepository) Update(user *entities.User) error {
	// GORM handles the SQL Update automatically
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresUserRepository) Delete(id uint) error {
	// GORM handles the SQL Delete automatically
	if err := r.db.Delete(&entities.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
