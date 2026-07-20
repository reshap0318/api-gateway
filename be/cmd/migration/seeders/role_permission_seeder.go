package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// SeedRolePermissions maps roles to their permissions.
func SeedRolePermissions(db *gorm.DB, roleIDs map[string]uint, permIDs map[string]uint) {
	fmt.Println("Seeding role permissions...")

	rolePerms := map[string][]string{
		"Super Admin": {},
		"Admin Gateway": {
			"service.index", "service.create", "service.edit", "service.delete", "service.health-check",
			"route.index", "route.create", "route.edit", "route.delete",
			"permission.index", "permission.create", "permission.edit", "permission.delete",
			"role.index", "role.create", "role.edit", "role.delete",
		},
		"Viewer": {
			"route.index",
			"service.index",
		},
	}

	count := 0
	for roleName, permNames := range rolePerms {
		roleID, roleOK := roleIDs[roleName]
		if !roleOK {
			log.Printf("Role %s not found, skipping", roleName)
			continue
		}

		// Super Admin gets ALL permissions
		if roleName == "Super Admin" {
			for permName := range permIDs {
				permID := permIDs[permName]

				var existing models.RoleHasPermission
				err := db.Where("role_id = ? AND permission_id = ?", roleID, permID).First(&existing).Error
				if err == nil {
					continue
				}

				rp := models.RoleHasPermission{
					RoleID:       roleID,
					PermissionID: permID,
				}

				if err := db.Create(&rp).Error; err != nil {
					log.Printf("Failed to create role_permission for %s-%s: %v", roleName, permName, err)
				} else {
					count++
				}
			}
			continue
		}

		for _, permName := range permNames {
			permID, permOK := permIDs[permName]
			if !permOK {
				log.Printf("Permission %s not found, skipping", permName)
				continue
			}

			var existing models.RoleHasPermission
			err := db.Where("role_id = ? AND permission_id = ?", roleID, permID).First(&existing).Error
			if err == nil {
				continue
			}

			rp := models.RoleHasPermission{
				RoleID:       roleID,
				PermissionID: permID,
			}

			if err := db.Create(&rp).Error; err != nil {
				log.Printf("Failed to create role_permission for %s-%s: %v", roleName, permName, err)
			} else {
				count++
			}
		}
	}

	fmt.Printf("✓ Seeded %d role permissions\n", count)
}
