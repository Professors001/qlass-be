package databases

import (
	"qlass-be/domain/entities"
	"qlass-be/domain/repositories"

	"gorm.io/gorm"
)

type postgresClassRepository struct {
	db *gorm.DB
}

// NewPostgresClassRepository creates a new instance of the repository
func NewPostgresClassRepository(db *gorm.DB) repositories.ClassRepository {
	return &postgresClassRepository{db: db}
}

func (r *postgresClassRepository) Create(class *entities.Class) error {
	// GORM handles the SQL Insert automatically
	if err := r.db.Create(class).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresClassRepository) GetByID(id uint) (*entities.Class, error) {
	var class entities.Class

	// SELECT * FROM classes WHERE id = ?
	if err := r.db.Preload("Owner").Where("id = ?", id).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *postgresClassRepository) GetByInviteCode(code string) (*entities.Class, error) {
	var class entities.Class

	// SELECT * FROM classes WHERE invite_code = ?
	if err := r.db.Preload("Owner").Where("invite_code = ?", code).First(&class).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *postgresClassRepository) GetByUserID(userID uint) ([]entities.Class, error) {
	var classes []entities.Class

	// Subquery: SELECT class_id FROM class_enrollments WHERE user_id = ?
	subQuery := r.db.Model(&entities.ClassEnrollment{}).Select("class_id").Where("user_id = ?", userID)

	// Main Query: OwnerID = ? OR ID IN (subQuery)
	if err := r.db.Preload("Owner").Where("owner_id = ?", userID).Or("id IN (?)", subQuery).Order("created_at desc").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}
