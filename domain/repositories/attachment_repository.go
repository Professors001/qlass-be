package repositories

import "qlass-be/domain/entities"

type AttachmentRepository interface {
	Create(attachment *entities.Attachment) error
	GetByID(id uint) (*entities.Attachment, error)
	GetByCourseMaterialID(courseMaterialID uint) ([]*entities.Attachment, error)
	GetBySubmissionID(submissionID uint) ([]*entities.Attachment, error)
	Update(attachment *entities.Attachment) error
	Delete(id uint) error
}
