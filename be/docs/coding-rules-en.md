# Go Project

## Project Overview

This is a **Go project** for REST API with a modern, layered architecture. The project is designed to accelerate API development with standard features already implemented.

### Tech Stack

- **Framework**: Gin (HTTP server)
- **Database**: MySQL (default) / PostgreSQL (optional)
- **ORM**: GORM
- **Cache**: Redis (optional)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Environment**: godotenv

### Architecture

This project uses **Dependency Injection** pattern with a layered structure:

```
Routes → Handlers → Services → Repositories → Database/Redis
              ↑
         Middleware (JWT Auth, CORS)
```

#### Layer Components:

| Layer | Location | Description |
|-------|----------|-------------|
| **Routes** | `internal/routes/` | Defines endpoints and maps to handlers |
| **Handlers** | `internal/handlers/` | HTTP handlers, request/response handling |
| **Services** | `internal/services/` | Business logic, JWT operations |
| **Repositories** | `internal/repositories/` | Data access layer with generic repository pattern |
| **Models** | `internal/models/` | GORM models for database entities |
| **DTOs** | `internal/dtos/` | Data Transfer Objects for request/response |
| **Middleware** | `internal/middleware/` | JWT auth, CORS, and other middleware |
| **DI Container** | `internal/di/` | Dependency Injection container |

## Project Structure

```
go-project/
├── cmd/
│   ├── api/              # Main application entry point
│   └── migration/        # Database migration scripts
├── docs/                 # API documentation
├── internal/             # Private application code
│   ├── clients/          # External API clients (email, etc.)
│   ├── database/         # Database connection & Redis
│   ├── di/               # Dependency Injection
│   ├── dtos/             # Data Transfer Objects
│   ├── handlers/         # HTTP handlers (single struct: Handlers)
│   ├── helpers/          # Helper functions
│   ├── middleware/       # Gin middleware
│   ├── models/           # GORM models (mandatory TableName)
│   ├── repositories/     # Data access layer
│   │   ├── 00_generic.go      # ⚠️ DO NOT MODIFY
│   │   ├── 00_transaction.go  # ⚠️ DO NOT MODIFY
│   │   └── 00_repository.go   # ⚠️ Register new repos here
│   ├── routes/           # Route definitions
│   └── services/         # Business logic (single struct: Services)
├── storage/              # Application storage (keys, logs, uploads)
├── .env                  # Environment variables (gitignore)
├── .env.example          # Environment template
├── go.mod                # Go module definition
└── main                  # Compiled binary
```

---

## ⚠️ Critical Rules (MANDATORY)

> 📖 **COMPLETE**: For step-by-step CRUD implementation flow, naming conventions, and anti-patterns, open **[docs/development-guide-en.md](docs/development-guide-en.md)**.

### 1. DO NOT Modify Core Repository Files

The following files **MUST NOT be modified** in any way:

- `internal/repositories/00_generic.go` — Generic repository pattern (core)
- `internal/repositories/00_transaction.go` — Transaction manager (core)

These files are the **foundation** of the repository architecture. If you need custom logic, create a new repository file (e.g., `user_repository.go`, `product_repository.go`).

### 2. Models MUST Define TableName

Every model **MUST** have a `TableName()` method to explicitly define the table name:

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
    // ... fields
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MANDATORY
func (User) TableName() string {
    return "users"
}
```

### 3. Services & Handlers MUST Use Single Struct Pattern

**All service methods MUST be methods of the `Services` struct**, and **all handler methods MUST be methods of the `Handlers` struct**.

**All functions MUST use feature name as prefix.**

> 📖 **COMPLETE**: See complete example in **[docs/development-guide-en.md](docs/development-guide-en.md)** — Service & Handler section.

```go
// ✅ CORRECT — Service with feature prefix (FeatureName + Action)
func (s *Services) PermissionCreate(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) { ... }
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error) { ... }

// ✅ CORRECT — Handler with feature prefix
func (h *Handlers) PermissionCreate(c *gin.Context) { ... }
func (h *Handlers) PermissionGetAll(c *gin.Context) { ... }

// ❌ WRONG — Action first instead of feature first
func (s *Services) CreatePermission(...) { }
func (s *Services) GetAllPermissions(...) { }
```

---

## 📝 Logging Rules (MANDATORY)

> ⚠️ **Logger available in Services**: `s.Logger`

### **When to Log?**

| Operation | Required? | Operation | Required? |
|-----------|-----------|-----------|-----------|
| **CREATE** | ✅ YES | **GET Simple** | ❌ NO |
| **UPDATE** | ✅ YES | **GET Complex** | ✅ YES |
| **DELETE** | ✅ YES | **Auth** | ✅ YES |

### **How to Use**

```go
s.Logger.LogStart("FuncName", "Message: %s", value)      // Start
s.Logger.LogStep("FuncName", "Step: %s", value)          // Step
s.Logger.LogStepWithPrefix("Func", "[OK]", "Done")       // Step with prefix
s.Logger.LogEnd("FuncName", "Success: %s", value)        // End
s.Logger.LogEndWithError("Func", "Error: %v", err)       // End + error
s.Logger.LogError("FuncName", "Error: %v", err)          // Error
s.Logger.LogWarn("FuncName", "Warning: %s", value)       // Warning
s.Logger.LogInfo("FuncName", "Info: %s", value)          // Info
```

### **Example — CREATE (MUST LOG)**

```go
func (s *Services) UserCreate(ctx context.Context, email string) error {
    s.Logger.LogStart("UserCreate", "Creating user: %s", email)

    user := &models.User{Email: email}
    if err := s.repo.User.Create(s.repo.User.DB, user); err != nil {
        s.Logger.LogEndWithError("UserCreate", "Failed: %v", err)
        return err
    }

    s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", email, user.ID)
    return nil
}
```

### **Example — GET Simple (NO LOG)**

```go
func (s *Services) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
    return s.repo.User.FindByID(nil, id)
}
```

**Output**: `storage/logs/YYYY-MM-DD.log` | Auto rotation & cleanup (30 days)

---

## 📚 Available Helpers & Clients

### **Helpers** (`internal/helpers/`)

> ⚠️ **IMPORTANT**: Always use existing helpers before creating new ones. If an operation can be generalized, create a new helper!

#### **Environment Helpers** (`env_helper.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GetEnv(key, default)` | Get environment variable with default | `helpers.GetEnv("APP_PORT", "8080")` |
| `GetEnvInt(key, default)` | Get environment variable as int | `helpers.GetEnvInt("JWT_EXPIRATION", 24)` |

#### **Error Helpers** (`error_helper.go`)
| Error | Description |
|-------|-------------|
| `ErrNotFound` | Record not found |
| `ErrInvalidToken` | Invalid JWT token |
| `ErrExpiredToken` | Token has expired |
| `ErrInvalidCredential` | Invalid email or password |
| `ErrUserExists` | User already exists |
| `ErrInvalidEmail` | Invalid email address |
| `ErrTokenExpired` | Reset token has expired |
| `ErrTokenUsed` | Reset token already used |
| `ErrTokenInvalid` | Invalid reset token |
| `ErrForbidden` | Forbidden access |

#### **Crypto Helpers** (`crypto_helper.go`)
| Function | Description | Example |
|----------|-------------|---------|
| `GenerateRandomString(length)` | Generate cryptographically secure random string | `helpers.GenerateRandomString(32)` |
| `HashString(str)` | Hash string using bcrypt | `helpers.HashString("mytoken")` |
| `VerifyString(str, hash)` | Verify string against hash | `helpers.VerifyString("token", hash)` |

#### **Response Helpers** (`response_helper.go`)
| Function | Status | Description | Example |
|----------|--------|-------------|---------|
| `OK(c, msg, data)` | 200 | Success response | `helpers.OK(c, "Success", user)` |
| `Created(c, msg, data)` | 201 | Created response | `helpers.Created(c, "Created", user)` |
| `BadRequest(c, msg)` | 400 | Bad request | `helpers.BadRequest(c, "Invalid input")` |
| `Unauthorized(c, msg)` | 401 | Unauthorized | `helpers.Unauthorized(c, "Invalid token")` |
| `Forbidden(c, msg)` | 403 | Forbidden | `helpers.Forbidden(c, "Access denied")` |
| `NotFound(c, msg)` | 404 | Not found | `helpers.NotFound(c, "User not found")` |
| `InternalServerError(c, msg)` | 500 | Server error | `helpers.InternalServerError(c, "Error")` |
| `ValidationErrorWithMap(c, errorsMap)` | 422 | Validation errors (map) | `helpers.ValidationErrorWithMap(c, map)` |
| `ValidationErrorWithField(c, field, msg)` | 422 | Single field error | `helpers.ValidationErrorWithField(c, "email", "Exists")` |

**Response format:**
```json
{
  "code": 200,
  "message": "Success message",
  "data": {}
}
```

---

### **Clients** (`internal/pkg/`)

#### **Email Client** (`internal/pkg/email/`)

> 📧 Email client is injected into Services via DI Container. Access via `s.EmailClient`.

**Struct:** `EmailClient`

**Methods:**

| Method | Description | Example |
|--------|-------------|---------|
| `IsConfigured()` | Check if SMTP is configured | `s.EmailClient.IsConfigured()` |
| `SendEmail(req)` | Send email with custom request | `s.EmailClient.SendEmail(req)` |
| `SendResetPasswordEmail(to, token, resetURL)` | Send reset password email | `s.EmailClient.SendResetPasswordEmail(email, token, url)` |

**DTOs:**

```go
// EmailRequest - for custom email
type EmailRequest struct {
    To      []string  // Recipients
    Subject string    // Email subject
    Body    string    // HTML body
    CC      []string  // CC recipients
    BCC     []string  // BCC recipients
}
```

**Example Usage:**

```go
// Send reset password email (recommended)
err := s.EmailClient.SendResetPasswordEmail(
    user.Email,
    token,
    "https://myapp.com/reset-password?token="+token,
)

// Send custom email
req := email.EmailRequest{
    To:      []string{"user@example.com"},
    Subject: "Welcome!",
    Body:    "<h1>Welcome to our app!</h1>",
}
err := s.EmailClient.SendEmail(req)
```

**Templates** (`internal/pkg/email/templates/`):

| Template | Description |
|----------|-------------|
| `ResetPasswordEmail(token, resetURL, appName)` | HTML template for reset password email |

---

### **Redis Cache** (`internal/database/redis_cache.go`)

> 🔴 Redis client is injected into Services via DI Container. Access via `s.RedisClient`.

**IMPORTANT**: Always check `s.RedisClient.IsCacheAvailable()` before accessing Redis. Redis errors **MUST NOT** cause operations to fail.

#### Usage

```go
// Check first if Redis is active
if s.RedisClient.IsCacheAvailable() {
    // SET - store data
    err := s.RedisClient.SetJSON("session:1", userDTO, time.Hour*24)
    if err != nil {
        s.Logger.LogWarn("FuncName", "Redis SET failed: %v", err) // fallback, don't return error
    }

    // GET - retrieve data
    var cached dtos.UserDTO
    if err := s.RedisClient.GetJSON("session:1", &cached); err == nil {
        // Cache hit - use cached data
    } else {
        // Cache miss/error - fallback to DB
    }

    // DELETE
    s.RedisClient.Delete("session:1")
}
```

#### Recommended Key Pattern
| Pattern | Example | Description |
|---------|--------|-----------|
| `session:{userID}` | `session:1` | User session cache |

#### Main Methods
| Method | Description |
|--------|-----------|
| `IsCacheAvailable()` | Check if Redis is active |
| `SetJSON(key, value, ttl)` | Store JSON with TTL |
| `GetJSON(key, &dest)` | Get & unmarshal JSON |
| `Delete(keys...)` | Delete key |

---

### **Access — Permission & Role Checking** (`internal/helpers/access.go`)

> 🔐 Access helper is injected into Services via DI Container. Access via `s.Access`.
> Also available in routes as `container.Access` for middleware.

**Struct:** `Access`

**3-Tier Caching:**
| Layer | Source | Behavior |
|-------|--------|----------|
| **L1** | Local in-memory (`sync.RWMutex` map) | Fastest, checked first |
| **L2** | Redis (`session:{userID}`) | Checked if L1 miss |
| **L3** | Database (user.Roles.Permissions) | Fallback if L1+L2 miss, then caches result to L1+L2 |

**Methods:**

| Method | Description | Example |
|--------|-------------|---------|
| `HasPermission(ctx, permissions...) bool` | Check if user has **ANY** of the specified permissions | `s.Access.HasPermission(ctx, "user.delete", "user.admin")` |
| `HasRole(ctx, role string) bool` | Check if user has the specified role | `s.Access.HasRole(ctx, "admin")` |
| `Invalidate(userID uint)` | Clear cached access data for a user | `s.Access.Invalidate(userID)` |

**Behavior:**
- Redis unavailable or cache miss → fallback to DB → cache result
- If user not found (all tiers) → returns `false`
- `Invalidate()` clears both local cache and Redis session

#### Usage in Services

```go
// Single permission check
func (s *Services) UserDelete(ctx context.Context, id uint) error {
    if !s.Access.HasPermission(ctx, "user.delete") {
        return helpers.ErrForbidden
    }
    // ... delete logic
}

// Multiple permissions - returns true if user has ANY of them
func (s *Services) UserAdminAction(ctx context.Context, req dtos.AdminRequest) error {
    if !s.Access.HasPermission(ctx, "user.delete", "user.update", "user.admin") {
        return helpers.ErrForbidden
    }
    // ...
}

// Role check
func (s *Services) SuperAdminOnly(ctx context.Context) error {
    if !s.Access.HasRole(ctx, "superadmin") {
        return helpers.ErrForbidden
    }
    // ...
}
```

#### Usage in Routes (Middleware)

```go
// internal/routes/user_route.go
func RegisterUserRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
    users := r.Group("/users")
    {
        users.POST("", middleware.RequirePermission(acc, "user.create"), handlers.UserCreate)
        users.GET("", middleware.RequirePermission(acc, "user.index"), handlers.UserGetAll)
        users.GET("/:id", middleware.RequirePermission(acc, "user.index"), handlers.UserGetByID)
        users.PUT("/:id", middleware.RequirePermission(acc, "user.edit"), handlers.UserUpdate)
        users.DELETE("/:id", middleware.RequirePermission(acc, "user.delete"), handlers.UserDelete)
    }
}
```

#### Cache Invalidation

Call `s.Access.Invalidate(userID)` when user's roles/permissions change:

```go
func (s *Services) UserUpdate(ctx context.Context, id uint, req dtos.UserUpdateRequest) (*dtos.UserDTO, string, error) {
    // ... update logic with role changes ...

    // Invalidate cached session so next request gets updated permissions
    s.Access.Invalidate(id)

    return &dto, oldAvatar, nil
}
```

#### Permission Naming Convention

Use dot separator: `{resource}.{action}`

| Permission | Description |
|------------|-------------|
| `user.index` | View users list |
| `user.create` | Create new user |
| `user.edit` | Update user |
| `user.delete` | Delete user |
| `role.index` | View roles |
| `role.create` | Create role |
| `role.edit` | Update role |
| `role.delete` | Delete role |
| `permission.index` | View permissions |
| `permission.create` | Create permission |
| `permission.edit` | Update permission |
| `permission.delete` | Delete permission |

---

### **Creating New Helpers**

> 💡 **Rule of Thumb**: If a function can be used in more than 1 place, make it a helper!

**Criteria for creating new helpers:**
1. ✅ General operation (not specific business logic)
2. ✅ Reusable across multiple services
3. ✅ No dependencies on repo, db, etc
4. ✅ Pure function (input → output, no side effects)

**Example of creating a new helper:**

```go
// internal/helpers/your_helper.go
package helpers

// FormatCurrency formats number as Indonesian Rupiah
func FormatCurrency(amount int64) string {
    return fmt.Sprintf("Rp %s", strings.ReplaceAll(
        strconv.FormatInt(amount, 10),
        "000",
        ".000",
    ))
}
```

**❌ DO NOT create helpers if:**
- ❌ Only used in 1 service
- ❌ Contains specific business logic
- ❌ Requires dependency injection (db, repo, etc)
- ❌ Has side effects (write to db, send email, etc)

---

## 🏗️ Adding New Features

### 1. Add New Model

```go
// internal/models/your_model.go
type YourModel struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    // ... fields
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// MANDATORY: Define table name
func (YourModel) TableName() string {
    return "your_table_name"
}
```

### 2. Add DTOs

**File:** `internal/dtos/your_feature_dto.go`

**DTO Strategy:**
- **Simple features** (e.g., Permission) → Merge Create & Update into **one** `{Feature}Request` struct
- **Complex features** (e.g., User with roles, avatar, etc.) → Separate `{Feature}CreateRequest` and `{Feature}UpdateRequest` structs

```go
// ✅ Simple feature — merged request
type PermissionRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description *string `json:"description" validate:"omitempty,max=255"`
}

// ✅ Complex feature — separate requests
type UserCreateRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Roles    []uint `json:"roles"`
}

type UserUpdateRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2"`
    Roles []uint `json:"roles"`
}
```

**MANDATORY:** Converter functions with feature prefix:
```go
func ToYourFeatureDTO(m *models.YourModel) YourFeatureDTO { ... }
func ToYourFeatureDTOList(models []models.YourModel) []YourFeatureDTO { ... }
```

### 3. Add New Repository

**IMPORTANT**: Use **GenericRepository** for all standard CRUD operations (Create, Read, Update, Delete, Find, Count, etc.). Add custom methods **only** for operations that cannot be handled by the generic repository, such as:
- Queries with complex JOINs (not regular preloads)
- Queries with specific subqueries
- Special batch operations
- Queries with specific aggregations

```go
// internal/repositories/your_repository.go
type YourRepository struct {
    *GenericRepository[YourModel]
}

func NewYourRepository(db *gorm.DB) *YourRepository {
    return &YourRepository{
        GenericRepository: NewGenericRepository(db, &YourModel{}),
    }
}

// ✅ EXAMPLE OF CORRECT CUSTOM METHOD (complex JOIN)
func (r *YourRepository) FindWithComplexJoin(id uint) (*YourModel, error) {
    var result YourModel
    query := r.DB.Joins("LEFT JOIN other_table ON other_table.your_id = your_models.id").
        Where("your_models.id = ?", id).
        First(&result)

    if query.Error != nil {
        return nil, query.Error
    }
    return &result, nil
}

// ❌ DO NOT DO THIS (use generic repo directly)
// func (r *YourRepository) FindByID(id uint) (*YourModel, error) {
//     return r.FindByID(r.DB, id) // This already exists in GenericRepository!
// }
```

### 4. Register Repository

After creating a repository, you **MUST** register it in `internal/repositories/00_repository.go`:

```go
// internal/repositories/00_repository.go
type Repositories struct {
    TxManager *TransactionManager
    User      *UserRepository
    YourModel *YourRepository  // ← Add this
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
    txManager := NewTransactionManager(db)
    userRepo := NewUserRepository(db)
    yourRepo := NewYourRepository(db)  // ← Initialize

    return &Repositories{
        TxManager: txManager,
        User:      userRepo,
        YourModel: yourRepo,  // ← Register
    }, nil
}
```

### 5. Add Service Method

> 📖 **COMPLETE**: See complete write/read operations example in **[docs/development-guide-en.md](docs/development-guide-en.md)** — Service section.

```go
// internal/services/your_feature_service.go
// MANDATORY: Feature name prefix (FeatureName + Action), TxManager for write, nil for read
func (s *Services) YourFeatureCreate(ctx context.Context, req dtos.YourFeatureRequest) (*dtos.YourFeatureDTO, error) { ... }
func (s *Services) YourFeatureGetByID(ctx context.Context, id uint) (*dtos.YourFeatureDTO, error) { ... }
```

The repository is automatically available in services via `s.repo.YourModel`.

#### Write Operations — Use `s.repo.TxManager.WithinTransaction()`

```go
func (s *Services) YourFeatureCreate(ctx context.Context, req dtos.YourFeatureRequest) (*dtos.YourFeatureDTO, error) {
    s.Logger.LogStart("YourFeatureCreate", "Creating: %s", req.Name)

    entity := &models.YourModel{Name: req.Name}

    var result *models.YourModel
    if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
        var err error
        result, err = s.repo.YourModel.Create(tx, entity)
        return err
    }); err != nil {
        s.Logger.LogEndWithError("YourFeatureCreate", "Failed: %v", err)
        return nil, err
    }

    dto := dtos.ToYourFeatureDTO(result)
    s.Logger.LogEnd("YourFeatureCreate", "Created: %s (ID: %d)", dto.Name, dto.ID)
    return &dto, nil
}
```

#### Write with Result Reload — Use `s.repo.TxManager.WithinTransactionWithResult()`

Use this when you need to return the entity after additional operations inside the transaction (e.g., assigning roles, reloading associations):

```go
func (s *Services) UserCreate(ctx context.Context, req dtos.UserCreateRequest) (*dtos.UserDTO, error) {
    s.Logger.LogStart("UserCreate", "Creating user: %s", req.Email)

    user := &models.User{Email: req.Email, Name: req.Name, Password: hashedPassword}

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
        s.Logger.LogEndWithError("UserCreate", "Failed: %v", err)
        return nil, err
    }

    result := res.(*models.User)
    dto := dtos.ToUserDTO(result)
    s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", dto.Email, dto.ID)
    return &dto, nil
}
```

#### Read Operations — Use `nil` parameter (NOT `s.repo.Feature.DB`)

```go
func (s *Services) YourFeatureGetAll(ctx context.Context) ([]dtos.YourFeatureDTO, error) {
    entities, err := s.repo.YourModel.FindAll(nil)
    if err != nil {
        return nil, err
    }
    return dtos.ToYourFeatureDTOList(entities), nil
}
```

#### Pagination — Use `FindAllWithOpts`

Use `FindAllWithOpts` for paginated, sorted, and searchable queries:

```go
func (s *Services) YourFeatureGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[models.YourModel], error) {
    if opts == nil {
        opts = &repositories.QueryOptions{}
    }
    if opts.SortBy == "" {
        opts.SortBy = "id"
    }
    if opts.Order == "" {
        opts.Order = "ASC"
    }

    return s.repo.YourModel.FindAllWithOpts(nil, opts)
}
```

`QueryOptions` fields:
- `Page`, `PageSize` — pagination
- `SortBy`, `Order` — sorting (`"ASC"` or `"DESC"`)
- `Search`, `SearchFields` — LIKE search across multiple fields
- `Preloads` — relations to preload (e.g., `[]string{"Roles", "Roles.Permissions"}`)

### 6. Add Handler Method

> 📖 **COMPLETE**: See complete example in **[docs/development-guide-en.md](docs/development-guide-en.md)** — Handler section.

```go
// internal/handlers/your_feature_handler.go
// MANDATORY: Feature name prefix — {Feature}{Action}
func (h *Handlers) YourFeatureCreate(c *gin.Context) { ... }
func (h *Handlers) YourFeatureGetAll(c *gin.Context) { ... }
```

### 7. Add Routes

```go
// internal/routes/your_routes.go
func RegisterYourRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
    your := r.Group("/your")
    {
        your.GET("/", handlers.YourHandler)
    }
}
```

### 8. Register Routes

Add to `cmd/api/main.go` in the appropriate route group (protected or public).

---

## 🧪 Testing

To add testing, create `_test.go` files in the same folder as the file being tested:

```bash
# Examples
internal/handlers/auth_handler_test.go
internal/services/auth_service_test.go
```

---

## 🔒 Security Notes

- ⚠️ **JWT_SECRET**: Change to a strong secret key in production
- ⚠️ **Database Password**: Do not commit credentials to git
- ⚠️ **GIN_MODE**: Set to `release` in production
- ⚠️ **CORS**: Configure `ALLOWED_ORIGINS` according to allowed domains
- ⚠️ **TRUSTED_PROXIES**: Set if using reverse proxy (nginx, load balancer)

---

> 📖 **For development workflow, build/run commands, and CRUD implementation guide**, see **[docs/development-guide-en.md](docs/development-guide-en.md)**.
