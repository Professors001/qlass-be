package entities

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	Filename   string  `json:"filename" gorm:"not null;type:varchar(255);comment:Original name e.g. Homework.pdf"`
	ObjectKey  string  `json:"-" gorm:"column:object_key;not null;type:varchar(500);comment:MinIO Object Key"`
	FileType   string  `json:"file_type" gorm:"type:varchar(50);comment:e.g. pdf, png, other"`
	FileSize   int     `json:"file_size" gorm:"comment:Size in bytes"`
	UploaderID uint    `json:"uploader_id"`
	Uploader   User    `json:"uploader" gorm:"foreignKey:UploaderID"`
	OwnerID    *uint   `json:"owner_id" gorm:"index;comment:Polymorphic ID"`
	OwnerType  *string `json:"owner_type" gorm:"index;type:varchar(255);comment:Polymorphic Type"`
}
