package databases

import (
	"log"
	"qlass-be/domain/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	users := []entities.User{
		{
			UniversityID:  "admin",
			Email:         "admin@qlass.com",
			PasswordHash:  string(hashedPassword),
			FirstName:     "System",
			LastName:      "Admin",
			Role:          "admin",
			IsVerified:    true,
			IsActive:      true,
			ProfileImgURL: "https://ui-avatars.com/api/?name=System+Admin",
		},
		{
			UniversityID:  "teacher",
			Email:         "teacher@qlass.com",
			PasswordHash:  string(hashedPassword),
			FirstName:     "John",
			LastName:      "Teacher",
			Role:          "teacher",
			IsVerified:    true,
			IsActive:      true,
			ProfileImgURL: "https://ui-avatars.com/api/?name=John+Teacher",
		},
		{
			UniversityID:  "student",
			Email:         "student@qlass.com",
			PasswordHash:  string(hashedPassword),
			FirstName:     "Jane",
			LastName:      "Student",
			Role:          "student",
			IsVerified:    true,
			IsActive:      true,
			ProfileImgURL: "https://ui-avatars.com/api/?name=Jane+Student",
		},
	}

	for _, user := range users {
		var count int64
		if err := db.Model(&entities.User{}).Where("university_id = ?", user.UniversityID).Count(&count).Error; err != nil {
			log.Printf("Error checking user %s: %v", user.UniversityID, err)
			continue
		}

		if count == 0 {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", user.UniversityID, err)
			} else {
				log.Printf("User %s seeded successfully", user.UniversityID)
			}
		}
	}
}
