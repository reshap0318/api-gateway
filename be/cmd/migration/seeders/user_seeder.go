package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// SeedUsers inserts default user data.
func SeedUsers(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding users...")

	defaultUsers := []struct {
		Email    string
		Password string
		Name     string
	}{
		{"suAdmin@app.com", "@dmin#123", "Super Admin"},
		{"admin@app.com", "Admin#123", "Admin"},
		{"editor@app.com", "Editor#123", "Editor"},
		{"viewer@app.com", "Viewer#123", "Viewer"},
	}

	resultMap := make(map[string]uint)

	for _, userData := range defaultUsers {
		var existing models.User
		result := db.Where("email = ?", userData.Email).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			hashedPassword, err := helpers.HashString(userData.Password)
			if err != nil {
				log.Printf("Failed to hash password for %s: %v", userData.Email, err)
				continue
			}

			user := models.User{
				Email:    userData.Email,
				Password: hashedPassword,
				Name:     userData.Name,
			}

			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", userData.Email, err)
			} else {
				resultMap[user.Email] = user.ID
			}
		} else if result.Error != nil {
			log.Printf("Failed to check user %s: %v", userData.Email, result.Error)
		} else {
			resultMap[userData.Email] = existing.ID
			fmt.Printf("  ⊘ User %s already exists, skipping\n", userData.Email)
		}
	}

	fmt.Printf("✓ Seeded %d users\n", len(resultMap))
	return resultMap
}
