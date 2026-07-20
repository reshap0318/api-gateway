package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

// SeedGatewayExample orchestrates the full "Data Master" example dataset in one call:
// example permissions (merged into the caller's permIDs map, since Go maps are reference
// types), the example Service, and its example Routes.
func SeedGatewayExample(db *gorm.DB, permIDs map[string]uint) {
	for name, id := range SeedMasterCategoryPermissions(db) {
		permIDs[name] = id
	}
	SeedGatewayExampleRoutes(db, SeedGatewayExampleService(db), permIDs)
}

// SeedMasterCategoryPermissions inserts example domain permissions for a fictional upstream
// "Data Master" service's Category resource. Permission names are scoped directly to the
// route/resource (e.g. "master-category.index") rather than reused from the built-in RBAC
// permissions (user.*, role.*, service.*, route.*), demonstrating that any permission name
// can be created ad-hoc via Permission Management and assigned to a Route, independent of
// what module actually owns it.
func SeedMasterCategoryPermissions(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding example permissions (master-category.*)...")

	permissions := []struct {
		Name        string
		Description string
	}{
		{"master-category.index", "View category list (contoh permission upstream Data Master)"},
		{"master-category.create", "Create new category (contoh permission upstream Data Master)"},
		{"master-category.update", "Update category (contoh permission upstream Data Master)"},
		{"master-category.delete", "Delete category (contoh permission upstream Data Master)"},
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

	fmt.Printf("✓ Seeded %d example permissions\n", len(resultMap))
	return resultMap
}

// SeedGatewayExampleService inserts one example upstream Service — "Data Master" —
// registered in the API Gateway.
func SeedGatewayExampleService(db *gorm.DB) uint {
	fmt.Println("Seeding example gateway service (Data Master)...")

	const name = "Data Master"

	var existing models.GatewayService
	err := db.Where("name = ?", name).First(&existing).Error
	if err == nil {
		fmt.Printf("  ⊘ Service %s already exists, skipping\n", name)
		return existing.ID
	}

	s := models.GatewayService{
		Name:         name,
		BaseURL:      "http://localhost:8081",
		BasePath:     "/master",
		Protocol:     "http",
		IsActive:     true,
		HealthStatus: "unknown",
	}

	if err := db.Create(&s).Error; err != nil {
		log.Printf("Failed to create service %s: %v", name, err)
		return 0
	}

	fmt.Printf("✓ Seeded gateway service: %s (ID %d)\n", name, s.ID)
	return s.ID
}

// SeedGatewayExampleRoutes inserts example Routes under Data Master, wired to the
// master-category.* permissions above. Demonstrates the Category resource's CRUD routes
// (create/read/update/delete) each guarded by its own dedicated permission.
func SeedGatewayExampleRoutes(db *gorm.DB, serviceID uint, permIDs map[string]uint) {
	if serviceID == 0 {
		log.Println("Data Master Service ID is 0, skipping example route seeding")
		return
	}

	fmt.Println("Seeding example gateway routes (Data Master)...")

	type routeSeed struct {
		Method      string
		PathPattern string
		MatchMode   string
		Permissions []string
	}

	// path_pattern is relative to the Service's base_path ("/master", see
	// SeedGatewayExampleService), so these resolve to e.g. "/master/categories",
	// "/master/categories/:id".
	routes := []routeSeed{
		{"GET", "/categories", "any", []string{"master-category.index"}},
		{"GET", "/categories/:id", "any", []string{"master-category.index"}},
		{"POST", "/categories", "any", []string{"master-category.create"}},
		{"PUT", "/categories/:id", "any", []string{"master-category.update"}},
		{"DELETE", "/categories/:id", "any", []string{"master-category.delete"}},
	}

	count := 0
	for _, r := range routes {
		var existing models.GatewayRoute
		err := db.Where("service_id = ? AND method = ? AND path_pattern = ?", serviceID, r.Method, r.PathPattern).First(&existing).Error
		if err == nil {
			fmt.Printf("  ⊘ Route %s %s already exists, skipping\n", r.Method, r.PathPattern)
			continue
		}

		route := models.GatewayRoute{
			ServiceID:           serviceID,
			Method:              r.Method,
			PathPattern:         r.PathPattern,
			PermissionMatchMode: r.MatchMode,
			IsActive:            true,
		}

		if err := db.Create(&route).Error; err != nil {
			log.Printf("Failed to create route %s %s: %v", r.Method, r.PathPattern, err)
			continue
		}

		if len(r.Permissions) > 0 {
			var perms []models.Permission
			for _, permName := range r.Permissions {
				permID, ok := permIDs[permName]
				if !ok {
					log.Printf("Permission %s not found, skipping for route %s %s", permName, r.Method, r.PathPattern)
					continue
				}
				perms = append(perms, models.Permission{ID: permID})
			}
			if len(perms) > 0 {
				if err := db.Model(&route).Association("Permissions").Append(perms); err != nil {
					log.Printf("Failed to assign permissions to route %s %s: %v", r.Method, r.PathPattern, err)
				}
			}
		}

		count++
	}

	fmt.Printf("✓ Seeded %d example gateway routes\n", count)
}
