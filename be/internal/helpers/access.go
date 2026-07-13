package helpers

import (
	"context"
	"fmt"
	"sync"

	"github.com/reshap0318/api-gateway/internal/database"
	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

type userAccessData struct {
	permissions map[string]bool
	roles       map[string]bool
}

// Access handles permission and role checking with 2-tier caching.
type Access struct {
	redis *database.RedisCache
	db    *gorm.DB
	mu    sync.RWMutex
	cache map[uint]*userAccessData
}

// NewAccess creates a new Access instance.
func NewAccess(redis *database.RedisCache, db *gorm.DB) *Access {
	return &Access{
		redis: redis,
		db:    db,
		cache: make(map[uint]*userAccessData),
	}
}

// getUserAccess retrieves user permissions and roles using 2-tier cache.
// L1: Local in-memory cache → L2: Database
func (a *Access) getUserAccess(userID uint) (*userAccessData, bool) {
	// L1: Check local cache
	a.mu.RLock()
	data, ok := a.cache[userID]
	a.mu.RUnlock()
	if ok {
		return data, true
	}

	// L2: Fallback to DB
	user, err := a.findUserWithRolesPermissions(userID)
	if err != nil {
		return nil, false
	}

	data = a.buildAccessDataFromUser(user)

	a.mu.Lock()
	a.cache[userID] = data
	a.mu.Unlock()

	return data, true
}

// findUserWithRolesPermissions fetches user with roles and permissions from DB.
func (a *Access) findUserWithRolesPermissions(userID uint) (*models.User, error) {
	var user models.User
	err := a.db.Preload("Roles.Permissions").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *Access) buildAccessDataFromUser(user *models.User) *userAccessData {
	data := &userAccessData{
		permissions: make(map[string]bool),
		roles:       make(map[string]bool),
	}
	for _, r := range user.Roles {
		data.roles[r.Name] = true
		for _, p := range r.Permissions {
			data.permissions[p.Name] = true
		}
	}
	return data
}

// HasPermission checks if the caller has ANY of the specified permissions.
func (a *Access) HasPermission(ctx context.Context, permissions ...string) bool {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(userID)
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if data.permissions[perm] {
			return true
		}
	}

	return false
}

// HasAllPermissions checks if the caller has ALL of the specified permissions.
// Used for GatewayRoute permission_match_mode = "all".
func (a *Access) HasAllPermissions(ctx context.Context, permissions ...string) bool {
	if len(permissions) == 0 {
		return true
	}

	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(userID)
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if !data.permissions[perm] {
			return false
		}
	}

	return true
}

// HasRole checks if the caller has the specified role.
func (a *Access) HasRole(ctx context.Context, role string) bool {
	userID := GetCallerID(ctx)
	if userID == 0 {
		return false
	}

	data, ok := a.getUserAccess(userID)
	if !ok {
		return false
	}

	return data.roles[role]
}

// Invalidate clears the cached access data for a user.
func (a *Access) Invalidate(userID uint) {
	a.mu.Lock()
	delete(a.cache, userID)
	a.mu.Unlock()

	if a.redis != nil && a.redis.IsCacheAvailable() {
		key := fmt.Sprintf("session:%d", userID)
		a.redis.Delete(key)
	}
}

// InvalidateAll clears cached access data for every user. Used when a Role/Permission
// change affects more users than we'd want to look up individually (e.g. editing a Role's
// permissions affects every user holding that Role).
func (a *Access) InvalidateAll() {
	a.mu.Lock()
	a.cache = make(map[uint]*userAccessData)
	a.mu.Unlock()
}
