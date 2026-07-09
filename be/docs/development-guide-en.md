# Development Guide

This guide provides operational instructions and workflows for contributing to the Go Boilerplate project.

## 🚀 Building and Running

### Prerequisites
- Go 1.25.0+
- MySQL 8+ / PostgreSQL 14+ (one of them)
- Redis (optional)

### Setup Environment
1. Copy environment file:
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` as needed.

### Commands
```bash
# Install dependencies
go mod download

# Run development
go run cmd/api/main.go

# Build binary
go build -o main cmd/api/main.go

# Run compiled binary
./main

# Run migration (folder available at cmd/migration/)
go run cmd/migration/*.go
```

## 🔄 Development Workflow (MANDATORY)

Every new feature or bug fix **MUST** follow this flow:

### Phase 1: Planning
1. **Create `plan.md`** in root directory:
   - Feature/bug description
   - Technical analysis (files to create/modify)
   - Implementation steps
   - Impact on existing code
2. **User Approval**: Wait for user confirmation.
3. **Create GitHub Issue** with format: `[BE] [Type] Title`
   - `[Feat]`, `[Bug]`, `[Fix]`, `[Refactor]`, `[Chore]`

### Phase 2: Implementation
1. **Create Branch**: `git checkout -b feat/feature-name`
2. **Implement**: Based on approved plan and issue.
3. **Create Pull Request**: Reference issue in description (`Closes #<id>`).
4. **User Approval**: Summary of changes $\rightarrow$ approval $\rightarrow$ merge $\rightarrow$ close issue.

---

## 🛠️ CRUD Implementation Flow

### 📋 Step-by-Step Order
Every new CRUD feature MUST follow this order:
1. Model          $\rightarrow$ `models/feature.go`
2. DTOs           $\rightarrow$ `dtos/feature_dto.go`
3. Repository     $\rightarrow$ `repositories/feature_repository.go`
4. Register Repo  $\rightarrow$ `repositories/00_repository.go`
5. Service        $\rightarrow$ `services/feature_service.go`
6. Handler        $\rightarrow$ `handlers/feature_handler.go`
7. Routes         $\rightarrow$ `routes/feature_route.go`
8. Register Routes $\rightarrow$ `cmd/api/main.go`

### 📁 Detailed Implementation

#### 1. Model

**File:** `internal/models/feature.go`

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Permission struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
    Description *string        `gorm:"size:255" json:"description"`    // nullable
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                // soft delete
}

// ⚠️ MANDATORY
func (Permission) TableName() string {
    return "permissions"
}
```

**Key points:**
- ✅ `TableName()` is **MANDATORY**
- ✅ `DeletedAt gorm.DeletedAt` for soft delete
- ✅ Use pointer (`*string`) for nullable fields
- ✅ Consistent JSON tags on all fields

---

#### 2. DTOs

**File:** `internal/dtos/feature_dto.go`

```go
package dtos

import "github.com/reshap0318/go-boilerplate/internal/models"

// ⚠️ MANDATORY: Prefix with feature name
// ⚠️ MANDATORY: For simple features, merge Create & Update into 1 struct
// ⚠️ NOTE: For complex features (e.g., User), separate structs are allowed (UserCreateRequest, UserUpdateRequest)
type PermissionRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description *string `json:"description" validate:"omitempty,max=255"`
}

// ⚠️ MANDATORY: Prefix with feature name
type PermissionDTO struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Description *string `json:"description"`
}

// ⚠️ MANDATORY: Converter with feature prefix
func ToPermissionDTO(p *models.Permission) PermissionDTO {
    return PermissionDTO{
        ID:          p.ID,
        Name:        p.Name,
        Description: p.Description,
    }
}

func ToPermissionDTOList(permissions []models.Permission) []PermissionDTO {
    result := make([]PermissionDTO, len(permissions))
    for i, p := range permissions {
        result[i] = ToPermissionDTO(&p)
    }
    return result
}
```

**Key points:**
- ✅ **ONE** request struct for both Create & Update (`PermissionRequest`)
- ✅ DTO struct for response (`PermissionDTO`)
- ✅ Converter functions `To{Feature}DTO()` and `To{Feature}DTOList()`
- ✅ Request fields: `required` for Create, `omitempty` for Update
- ✅ Import `models` package

---

#### 3. Repository

**File:** `internal/repositories/feature_repository.go`

```go
package repositories

import (
    "github.com/reshap0318/go-boilerplate/internal/models"
    "gorm.io/gorm"
)

// ⚠️ MANDATORY: Extend GenericRepository
type PermissionRepository struct {
    *GenericRepository[models.Permission]
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
    return &PermissionRepository{
        GenericRepository: NewGenericRepository(db, &models.Permission{}),
    }
}
```

**Key points:**
- ✅ **ONLY** extend `GenericRepository[Model]`
- ✅ Custom methods **ONLY** for complex queries (JOINs, subqueries, aggregations)
- ❌ DO NOT duplicate methods already in GenericRepository (Create, FindByID, FindAll, Update, Delete, etc.)
- ❌ DO NOT modify `00_generic.go` or `00_transaction.go`

---

#### 4. Register Repository

**File:** `internal/repositories/00_repository.go`

```go
type Repositories struct {
    TxManager     *TransactionManager
    User          *UserRepository
    PasswordReset *PasswordResetRepository
    Permission    *PermissionRepository  // ← Add this
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
    txManager := NewTransactionManager(db)
    userRepo := NewUserRepository(db)
    passwordResetRepo := NewPasswordResetRepository(db)
    permissionRepo := NewPermissionRepository(db)  // ← Initialize

    return &Repositories{
        TxManager:     txManager,
        User:          userRepo,
        PasswordReset: passwordResetRepo,
        Permission:    permissionRepo,  // ← Register
    }, nil
}
```

---

#### 5. Service

**File:** `internal/services/feature_service.go`

**Write Operations (Create/Update/Delete) — MANDATORY Transaction**

```go
// ⚠️ MANDATORY: Feature name FIRST ({Feature}{Action})
// ⚠️ MANDATORY: Use TxManager for write operations
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
    s.Logger.LogStart("PermissionCreate", "Creating permission: %s", req.Name)

    permission := &models.Permission{
        Name:        req.Name,
        Description: req.Description,
    }

    var result *models.Permission
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.Permission.Create(tx, permission)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionCreate", "Failed to create permission: %v", err)
        return nil, err
    }

    dto := dtos.ToPermissionDTO(result)
    s.Logger.LogEnd("PermissionCreate", "Permission created: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}
```

**Read Operations (Get/Find) — MANDATORY use `nil`**

```go
// ⚠️ MANDATORY: Feature name FIRST ({Feature}{Action})
// ⚠️ MANDATORY: Use nil for read operations (NOT s.repo.Permission.DB)
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error) {
    permissions, err := s.repo.Permission.FindAll(nil)  // ← nil, NOT .DB
    if err != nil {
        return nil, err
    }
    return dtos.ToPermissionDTOList(permissions), nil
}

// ⚠️ Simple GET does NOT need logging
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
    permission, err := s.repo.Permission.FindByID(nil, id)  // ← nil, NOT .DB
    if err != nil {
        return nil, helpers.ErrNotFound
    }
    return dtos.ToPermissionDTO(permission), nil
}
```

**Update with Transaction**

```go
func (s *Services) PermissionUpdate(ctx context.Context, id uint, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
    s.Logger.LogStart("PermissionUpdate", "Updating permission ID: %d", id)

    permission := &models.Permission{
        ID: id,
    }
    if req.Name != "" {
        permission.Name = req.Name
    }
    if req.Description != nil {
        permission.Description = req.Description
    }

    var result *models.Permission
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.Permission.Update(tx, &models.Permission{ID: id}, permission)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionUpdate", "Failed to update permission: %v", err)
        return nil, err
    }

    dto := dtos.ToPermissionDTO(result)
    s.Logger.LogEnd("PermissionUpdate", "Permission updated: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}
```

**Delete with Transaction**

```go
func (s *Services) PermissionDelete(ctx context.Context, id uint) error {
    s.Logger.LogStart("PermissionDelete", "Deleting permission ID: %d", id)

    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        _, err := s.repo.Permission.Delete(tx, id)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("PermissionDelete", "Failed to delete permission: %v", err)
        return err
    }

    s.Logger.LogEnd("PermissionDelete", "Permission deleted: ID: %d", id)
    return nil
}
```

**Service key points:**
- ✅ **Write** (Create/Update/Delete) $\rightarrow$ `s.repo.TxManager.WithinTransaction()`
- ✅ **Write with result reload** $\rightarrow$ `s.repo.TxManager.WithinTransactionWithResult()` (for complex operations like role assignment)
- ✅ **Read** (Get/Find) $\rightarrow$ `nil` parameter (NOT `s.repo.Feature.DB`)
- ✅ **Function name** $\rightarrow$ Feature FIRST (`PermissionCreate`, `PermissionGetAll`, etc.)
- ✅ **Logging** $\rightarrow$ Create/Update/Delete MUST log, Simple GET does NOT need logging
- ✅ **Error handling** $\rightarrow$ use `helpers.ErrNotFound` for record not found
- ✅ **Single struct** $\rightarrow$ all methods on `(s *Services)`, do NOT create separate structs
- ✅ **Update/Delete** $\rightarrow$ NO FindByID check before operation (generic repo handles not found)

**Advanced: WithinTransactionWithResult**

Use `WithinTransactionWithResult` when you need to return the created/updated entity after additional operations inside the transaction (e.g., assigning roles, reloading associations):

```go
func (s *Services) UserCreate(ctx context.Context, req dtos.UserCreateRequest) (*dtos.UserDTO, error) {
    s.Logger.LogStart("UserCreate", "Creating user: %s", req.Email)

    user := &models.User{
        Email:    req.Email,
        Name:     req.Name,
        Password: hashedPassword,
    }

    res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
        result, err := s.repo.User.Create(tx, user)
        if err != nil {
            return nil, err
        }

        // Assign roles inside transaction
        var roles []models.Role
        for _, roleID := range req.Roles {
            roles = append(roles, models.Role{ID: roleID})
        }
        if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
            return nil, err
        }

        // Reload with associations
        return s.repo.User.FindByID(tx, result.ID, "Roles")
    })
    if err != nil {
        s.Logger.LogEndWithError("UserCreate", "Failed to create user: %v", err)
        return nil, err
    }

    result := res.(*models.User)
    dto := dtos.ToUserDTO(result)
    s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", dto.Email, dto.ID)
    return &dto, nil
}
```

**Advanced: Pagination with FindAllWithOpts**

Use `FindAllWithOpts` for paginated, sorted, and searchable queries:

```go
func (s *Services) PermissionGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[models.Permission], error) {
    if opts == nil {
        opts = &repositories.QueryOptions{}
    }
    if opts.SortBy == "" {
        opts.SortBy = "id"
    }
    if opts.Order == "" {
        opts.Order = "ASC"
    }

    return s.repo.Permission.FindAllWithOpts(nil, opts)
}
```

QueryOptions fields:
- `Page`, `PageSize` — pagination
- `SortBy`, `Order` — sorting (`"ASC"` or `"DESC"`)
- `Search`, `SearchFields` — LIKE search across multiple fields
- `Preloads` — relations to preload (e.g., `[]string{"Roles", "Roles.Permissions"}`)

---

#### 6. Handler

**File:** `internal/handlers/feature_handler.go`

```go
package handlers

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "github.com/reshap0318/go-boilerplate/internal/dtos"
    "github.com/reshap0318/go-boilerplate/internal/helpers"
)

// ⚠️ MANDATORY: Feature name FIRST — {Feature}{Action}
func (h *Handlers) PermissionCreate(c *gin.Context) {
    var req dtos.PermissionRequest
    if err := c.BindJSON(&req); err != nil {
        helpers.BadRequest(c, "Invalid JSON payload")
        return
    }

    if err := h.Validate.Struct(req); err != nil {
        helpers.ValidationErrorWithMap(c, h.getErrorsMap(err))
        return
    }

    dto, err := h.svcs.PermissionCreate(c.Request.Context(), req)
    if err != nil {
        helpers.InternalServerError(c, "Failed to create permission")
        return
    }

    helpers.Created(c, "Permission created successfully", dto)
}

func (h *Handlers) PermissionGetAll(c *gin.Context) {
    dtos, err := h.svcs.PermissionGetAll(c.Request.Context())
    if err != nil {
        helpers.InternalServerError(c, "Failed to fetch permissions")
        return
    }

    helpers.OK(c, "Permissions fetched successfully", dtos)
}

func (h *Handlers) PermissionGetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        helpers.BadRequest(c, "Invalid permission ID")
        return
    }

    dto, err := h.svcs.PermissionGetByID(c.Request.Context(), uint(id))
    if err != nil {
        helpers.NotFound(c, "Permission not found")
        return
    }

    helpers.OK(c, "Permission fetched successfully", dto)
}

func (h *Handlers) PermissionUpdate(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        helpers.BadRequest(c, "Invalid permission ID")
        return
    }

    var req dtos.PermissionRequest
    if err := c.BindJSON(&req); err != nil {
        helpers.BadRequest(c, "Invalid JSON payload")
        return
    }

    if err := h.Validate.Struct(req); err != nil {
        helpers.ValidationErrorWithMap(c, h.getErrorsMap(err))
        return
    }

    dto, err := h.svcs.PermissionUpdate(c.Request.Context(), uint(id), req)
    if err != nil {
        helpers.NotFound(c, "Permission not found")
        return
    }

    helpers.OK(c, "Permission updated successfully", dto)
}

func (h *Handlers) PermissionDelete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        helpers.BadRequest(c, "Invalid permission ID")
        return
    }

    err = h.svcs.PermissionDelete(c.Request.Context(), uint(id))
    if err != nil {
        helpers.NotFound(c, "Permission not found")
        return
    }

    helpers.OK(c, "Permission deleted successfully", nil)
}
```

**Handler key points:**
- ✅ **Function name** $\rightarrow$ `{Feature}{Action}` (`PermissionCreate`, `PermissionGetAll`, etc.)
- ✅ **Parse ID** $\rightarrow$ `strconv.ParseUint` with error handling
- ✅ **Bind JSON** $\rightarrow$ `c.BindJSON()` + `h.Validate.Struct()` with validation
- ✅ **Response helper** $\rightarrow$ `helpers.Created`, `helpers.OK`, `helpers.BadRequest`, etc.
- ✅ **Single struct** $\rightarrow$ all methods on `(h *Handlers)`, do NOT create separate structs
- ✅ **Service calls** $\rightarrow$ MUST match service function name (`h.svcs.PermissionCreate`, NOT `h.svcs.CreatePermission`)

---

#### 7. Routes

**File:** `internal/routes/feature_route.go`

```go
package routes

import (
    "github.com/gin-gonic/gin"

    "github.com/reshap0318/go-boilerplate/internal/handlers"
    "github.com/reshap0318/go-boilerplate/internal/helpers"
    "github.com/reshap0318/go-boilerplate/internal/middleware"
)

func RegisterPermissionRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
    permissions := r.Group("/permissions")
    {
        permissions.POST("", middleware.RequirePermission(acc, "permission.create"), handlers.PermissionCreate)
        permissions.GET("", middleware.RequirePermission(acc, "permission.index"), handlers.PermissionGetAll)
        permissions.GET("/:id", middleware.RequirePermission(acc, "permission.index"), handlers.PermissionGetByID)
        permissions.PUT("/:id", middleware.RequirePermission(acc, "permission.edit"), handlers.PermissionUpdate)
        permissions.DELETE("/:id", middleware.RequirePermission(acc, "permission.delete"), handlers.PermissionDelete)
    }
}
```

---

#### 8. Register Routes

**File:** `cmd/api/main.go`

Add to the `protected` group (after JWT middleware):

```go
protected := apiGroup.Group("")
protected.Use(middleware.JWTAuth(container.Services))
{
    routes.RegisterAuthProtectedRoutes(protected, container.Handlers)
    routes.RegisterPermissionRoutes(protected, container.Handlers, container.Access)  // ← Add this
}
```

---

## 🔐 Scoped Permission Implementation

This pattern is used when a role must only access data within their own unit/satker.

### Example: User Feature with Satker Restriction

**Permission List:**
| Permission | Description |
|------------|-------------|
| `user.index` | View all users |
| `user.index-satker` | View users in own satker only |
| `user.create` | Create user |
| `user.edit` | Edit user |
| `user.delete` | Delete user |

**Route Registration:**
All GET endpoints use the same middleware permission. The actual data filtering happens in the Service.

```go
func RegisterUserRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
    users := r.Group("/users")
    {
        users.GET("", middleware.RequirePermission(acc, "user.index", "user.index-satker"), handlers.UserGetAll)
        users.GET("/:id", middleware.RequirePermission(acc, "user.index", "user.index-satker"), handlers.UserGetByID)
        users.POST("", middleware.RequirePermission(acc, "user.create"), handlers.UserCreate)
        users.PUT("/:id", middleware.RequirePermission(acc, "user.edit"), handlers.UserUpdate)
        users.DELETE("/:id", middleware.RequirePermission(acc, "user.delete"), handlers.UserDelete)
    }
}
```

**Service Implementation:**
The Service checks permissions and applies the satker filter if necessary.

```go
func (s *Services) UserGetAll(ctx context.Context) ([]dtos.UserDTO, error) {
    // 1. Check if user has full access
    if s.Access.HasPermission(ctx, "user.index") {
        users, err := s.repo.User.FindAll(nil, "Roles")
        if err != nil {
            return nil, err
        }
        return dtos.ToUserDTOList(users), nil
    }

    // 2. Check if user has satker-scoped access
    if s.Access.HasPermission(ctx, "user.index-satker") {
        callerID := helpers.GetCallerID(ctx)
        
        // Get caller's satker ID (from DB or cached context)
        caller, err := s.repo.User.FindByID(nil, callerID, "Satker")
        if err != nil {
            return nil, helpers.ErrForbidden
        }

        // Filter by satker ID using custom repository method or generic FindByFieldMap
        users, err := s.repo.User.FindByFieldMap(nil, map[string]interface{}{
            "satker_id": caller.SatkerID,
        }, "Roles")
        if err != nil {
            return nil, err
        }
        return dtos.ToUserDTOList(users), nil
    }

    // 3. No permission
    return nil, helpers.ErrForbidden
}
```

**Key Points:**
- Routes allow BOTH permissions (`user.index` OR `user.index-satker`).
- Service decides the data scope based on which permission the user actually has.
- Use `helpers.GetCallerID(ctx)` to identify the current user.
- Satker filtering logic belongs in the Service layer.

---

## 🗄️ Database Schema (Reference)

### Users Table
| Column | Type | Constraints |
|--------|------|-------------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT |
| email | string(255) | UNIQUE, NOT NULL |
| password | string(255) | NOT NULL (hashed) |
| name | string(255) | - |
| created_at | timestamp | - |
| updated_at | timestamp | - |
| deleted_at | timestamp | SOFT DELETE index |

---

## 🏗️ Services Struct (Dependencies)

```go
type Services struct {
    repo         *repositories.Repositories  // Access repos: s.repo.User, s.repo.Permission, etc.
    RedisClient  *database.RedisCache        // Redis cache client
    EmailClient  *email.EmailClient          // Email client
    JWKSManager  *services.JWKSManager       // JWKS manager
    Access       *helpers.Access             // Permission/role checker
    Logger       *helpers.Logger             // Logger
    cfg          *JWTConfig                  // JWT config
}
```

---

## ✅ Pre-Push Checklist
- [ ] Model has `TableName()` method
- [ ] DTO variables use feature prefix
- [ ] Request DTO: merged for simple features OR separate (`{Feature}CreateRequest` + `{Feature}UpdateRequest`) for complex features
- [ ] Service functions use feature prefix (`{Feature}{Action}`)
- [ ] Handler functions use feature prefix (`{Feature}{Action}`)
- [ ] Write operations use `TxManager.WithinTransaction()` or `WithinTransactionWithResult()`
- [ ] Read operations use `nil` (NOT `.DB`)
- [ ] Logging on Create/Update/Delete
- [ ] Repository registered in `00_repository.go`
- [ ] Routes registered in `cmd/api/main.go`
- [ ] Build success (`go build ./...`)
- [ ] Vet clean (`go vet ./...`)
