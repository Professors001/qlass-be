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
			UniversityID: "admin",
			Email:        "admin@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "System",
			LastName:     "Admin",
			Role:         "admin",
			IsVerified:   true,
			IsActive:     true,
		},
		{
			UniversityID: "fsci001",
			Email:        "fsci001@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "John",
			LastName:     "Teacher",
			Role:         "teacher",
			IsVerified:   true,
			IsActive:     true,
		},
		{
			UniversityID: "6910450001",
			Email:        "6910450001@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Alice",
			LastName:     "Johnson",
			Role:         "student",
			IsVerified:   true,
			IsActive:     true,
		},
		{
			UniversityID: "6910450002",
			Email:        "6910450002@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Bob",
			LastName:     "Smith",
			Role:         "student",
			IsVerified:   true,
			IsActive:     true,
		},
		{
			UniversityID: "6910450003",
			Email:        "6910450003@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Carol",
			LastName:     "Williams",
			Role:         "student",
			IsVerified:   true,
			IsActive:     true,
		},
		{
			UniversityID: "6910450004",
			Email:        "6910450004@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "David",
			LastName:     "Brown",
			Role:         "student",
			IsVerified:   true,
			IsActive:     true,
		},
		{
			UniversityID: "6910450005",
			Email:        "6910450005@qlass.com",
			PasswordHash: string(hashedPassword),
			FirstName:    "Emma",
			LastName:     "Davis",
			Role:         "student",
			IsVerified:   true,
			IsActive:     true,
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
