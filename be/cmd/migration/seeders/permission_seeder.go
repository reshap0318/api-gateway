package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// SeedPermissions inserts default permission data.
func SeedPermissions(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding permissions...")

	permissions := []struct {
		Name        string
		Description string
	}{
		// User
		{"user.index", "View users list"},
		{"user.create", "Create new user"},
		{"user.edit", "Update user"},
		{"user.delete", "Delete user"},
		// Role
		{"role.index", "View roles list"},
		{"role.index-su", "View super admin role (ID 1)"},
		{"role.create", "Create new role"},
		{"role.edit", "Update role"},
		{"role.delete", "Delete role"},
		// Permission
		{"permission.index", "View permissions list"},
		{"permission.create", "Create new permission"},
		{"permission.edit", "Update permission"},
		{"permission.delete", "Delete permission"},
		// Gateway Service
		{"service.index", "View gateway services list"},
		{"service.create", "Register new gateway service"},
		{"service.edit", "Update gateway service"},
		{"service.delete", "Delete gateway service"},
		{"service.health-check", "Trigger manual health check on a service"},
		// Gateway Route
		{"route.index", "View gateway routes list"},
		{"route.create", "Create new gateway route"},
		{"route.edit", "Update gateway route"},
		{"route.delete", "Delete gateway route"},
		// Gateway Audit Trail
		{"audit.index", "View gateway audit log"},
	}

	resultMap := make(map[string]uint)

	for _, perm := range permissions {
		var existing models.Permission
		err := db.Where("name = ?", perm.Name).First(&existing).Error
		if err == nil {
			resultMap[perm.Name] = existing.ID
			fmt.Printf("  ⊘ Permission %s already exists, skipping\n", perm.Name)
			continue
		}

		p := models.Permission{
			Name:        perm.Name,
			Description: strPtr(perm.Description),
		}

		if err := db.Create(&p).Error; err != nil {
			log.Printf("Failed to create permission %s: %v", perm.Name, err)
		} else {
			resultMap[perm.Name] = p.ID
		}
	}

	fmt.Printf("✓ Seeded %d permissions\n", len(resultMap))
	return resultMap
}

func strPtr(s string) *string {
	return &s
}
