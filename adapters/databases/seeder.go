package databases

import (
	"log"
	"qlass-be/domain/entities"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) {
	var count int64
	// Check if an admin already exists
	if err := db.Model(&entities.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		log.Printf("Error checking for admin user: %v", err)
		return
	}

	if count > 0 {
		return
	}

	log.Println("Seeding admin user...")

	// Create Admin User
	password := "admin1234"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	admin := entities.User{
		UniversityID:  "ADMIN001",
		Email:         "admin@qlass.com",
		PasswordHash:  string(hashedPassword),
		FirstName:     "System",
		LastName:      "Admin",
		Role:          "admin",
		IsVerified:    true,
		IsActive:      true,
		ProfileImgURL: "https://ui-avatars.com/api/?name=System+Admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
	} else {
		log.Println("Admin user seeded successfully (Email: admin@qlass.com, Password: admin1234)")
	}
}
