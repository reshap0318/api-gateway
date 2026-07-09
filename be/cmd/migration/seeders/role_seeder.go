package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// SeedRoles inserts default role data.
func SeedRoles(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding roles...")

	roles := []struct {
		Name        string
		Description string
	}{
		{"Super Admin", "Full access to all features"},
		{"Admin", "Manage users, roles, and permissions"},
		{"Editor", "Manage users and view data"},
		{"Viewer", "Read-only access"},
	}

	resultMap := make(map[string]uint)

	for _, roleData := range roles {
		var existing models.Role
		err := db.Where("name = ?", roleData.Name).First(&existing).Error
		if err == nil {
			resultMap[roleData.Name] = existing.ID
			fmt.Printf("  ⊘ Role %s already exists, skipping\n", roleData.Name)
			continue
		}

		role := models.Role{
			Name:        roleData.Name,
			Description: strPtr(roleData.Description),
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("Failed to create role %s: %v", roleData.Name, err)
		} else {
			resultMap[roleData.Name] = role.ID
		}
	}

	fmt.Printf("✓ Seeded %d roles\n", len(resultMap))
	return resultMap
}
