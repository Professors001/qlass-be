package entities

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	Filename         string `json:"filename" gorm:"not null;type:varchar(255);comment:Original name e.g. Homework.pdf"`
	FileURL          string `json:"file_url" gorm:"not null;type:varchar(500);comment:MinIO URL"`
	FileType         string `json:"file_type" gorm:"type:varchar(50);comment:e.g. pdf, png, other"`
	FileSize         int    `json:"file_size" gorm:"comment:Size in bytes"`
	UploaderID       *int   `json:"uploader_id"`
	CourseMaterialID *int   `json:"course_material_id"`
	SubmissionID     *int   `json:"submission_id"`
}
