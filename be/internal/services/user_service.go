package services

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// UserCreate creates a new user with optional roles.
func (s *Services) UserCreate(ctx context.Context, req dtos.UserCreateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserCreate", "Creating user: %s", req.Email)

	exists, err := s.repo.User.Exists(nil, map[string]interface{}{"email": req.Email})
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to check email: %v", err)
		return nil, err
	}
	if exists {
		s.Logger.LogEndWithError("UserCreate", "Email already exists: %s", req.Email)
		return nil, &helpers.FieldError{Field: "email", Message: "user already exists"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to hash password: %v", err)
		return nil, err
	}

	avatarPath := ""
	if req.Avatar != "" {
		avatarPath, err = helpers.MoveFile(req.Avatar, "storage/tmp", "storage/avatars")
		if err != nil {
			s.Logger.LogStep("UserCreate", "Failed to move avatar: %v", err)
			avatarPath = ""
		}
	}

	user := &models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
		Avatar:   avatarPath,
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.Create(tx, user)
		if err != nil {
			return nil, err
		}

		var roles []models.Role
		for _, roleID := range req.Roles {
			roles = append(roles, models.Role{ID: roleID})
		}
		if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
			return nil, fmt.Errorf("failed to assign roles: %w", err)
		}

		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("UserCreate", "Failed to create user: %v", err)
		return nil, err
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "User Created",
		Message: fmt.Sprintf("New user created: %s", req.Email),
		Data: map[string]interface{}{
			"id":    result.ID,
			"email": result.Email,
		},
	})

	s.Logger.LogEnd("UserCreate", "User created: %s (ID: %d)", dto.Email, dto.ID)
	return &dto, nil
}

// UserGetAll returns all users with roles (no pagination).
func (s *Services) UserGetAll(ctx context.Context) ([]dtos.UserDTO, error) {
	users, err := s.repo.User.FindAll(nil, "Roles")
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dtos.UserDTO, len(users))
	for i, u := range users {
		userDTOs[i] = dtos.ToUserDTO(&u)
	}
	return userDTOs, nil
}

// UserGetAllPaginated returns paginated users with roles.
func (s *Services) UserGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.UserDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}
	opts.Preloads = []string{"Roles"}

	result, err := s.repo.User.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dtos.UserDTO, len(result.Data))
	for i, u := range result.Data {
		userDTOs[i] = dtos.ToUserDTO(&u)
	}

	return &repositories.PagedResult[dtos.UserDTO]{
		Data:       userDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// UserGetByID returns a user by ID with roles.
func (s *Services) UserGetByID(ctx context.Context, id uint) (*dtos.UserDTO, error) {
	user, err := s.repo.User.FindByID(nil, id, "Roles")
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToUserDTO(user)
	return &dto, nil
}

// ProfileGet returns the authenticated user's profile.
func (s *Services) ProfileGet(ctx context.Context, userID uint) (*dtos.UserDTO, error) {
	user, err := s.repo.User.FindByID(nil, userID, "Roles")
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToUserDTO(user)
	return &dto, nil
}

// ProfileUpdate updates the authenticated user's profile.
func (s *Services) ProfileUpdate(ctx context.Context, userID uint, req dtos.ProfileUpdateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("ProfileUpdate", "Updating profile for user ID: %d", userID)

	existing, err := s.repo.User.FindByID(nil, userID)
	if err != nil {
		s.Logger.LogEndWithError("ProfileUpdate", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	updates := map[string]interface{}{
		"name": req.Name,
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Logger.LogEndWithError("ProfileUpdate", "Failed to hash password: %v", err)
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	oldAvatar := ""
	if req.Avatar != "" {
		avatarPath, err := helpers.MoveFile(req.Avatar, "storage/tmp", "storage/avatars")
		if err != nil {
			s.Logger.LogStep("ProfileUpdate", "Failed to move avatar: %v", err)
		} else {
			updates["avatar"] = avatarPath
			oldAvatar = existing.Avatar
		}
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.UpdateMap(tx, &models.User{ID: userID}, updates)
		if err != nil {
			return nil, err
		}

		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("ProfileUpdate", "Failed to update profile: %v", err)
		return nil, err
	}

	if oldAvatar != "" {
		helpers.DeleteFile(oldAvatar)
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Profile Updated",
		Message: "User profile has been updated",
		Data: map[string]interface{}{
			"id": result.ID,
		},
	})

	s.Access.Invalidate(userID)

	s.Logger.LogEnd("ProfileUpdate", "Profile updated for user: %s", dto.Email)
	return &dto, nil
}

// UserUpdate updates an existing user with optional roles.
func (s *Services) UserUpdate(ctx context.Context, id uint, req dtos.UserUpdateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserUpdate", "Updating user ID: %d", id)

	existing, err := s.repo.User.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("UserUpdate", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if existing.Email != req.Email {
		exists, err := s.repo.User.Exists(nil, map[string]interface{}{"email": req.Email})
		if err != nil {
			s.Logger.LogEndWithError("UserUpdate", "Failed to check email: %v", err)
			return nil, err
		}
		if exists {
			s.Logger.LogEndWithError("UserUpdate", "Email already exists: %s", req.Email)
			return nil, &helpers.FieldError{Field: "email", Message: "user already exists"}
		}
	}

	updates := map[string]interface{}{
		"name":  req.Name,
		"email": req.Email,
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.Logger.LogEndWithError("UserUpdate", "Failed to hash password: %v", err)
			return nil, err
		}
		updates["password"] = string(hashedPassword)
	}

	oldAvatar := ""
	if req.Avatar != "" {
		avatarPath, err := helpers.MoveFile(req.Avatar, "storage/tmp", "storage/avatars")
		if err != nil {
			s.Logger.LogStep("UserUpdate", "Failed to move avatar: %v", err)
		} else {
			updates["avatar"] = avatarPath
			oldAvatar = existing.Avatar
		}
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		result, err := s.repo.User.UpdateMap(tx, &models.User{ID: id}, updates)
		if err != nil {
			return nil, err
		}

		if err := tx.Model(&result).Association("Roles").Clear(); err != nil {
			return nil, err
		}

		var roles []models.Role
		for _, roleID := range req.Roles {
			roles = append(roles, models.Role{ID: roleID})
		}
		if err := tx.Model(&result).Association("Roles").Append(roles); err != nil {
			return nil, fmt.Errorf("failed to assign roles: %w", err)
		}

		reloaded, err := s.repo.User.FindByID(tx, result.ID, "Roles")
		if err != nil {
			return nil, err
		}

		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("UserUpdate", "Failed to update user: %v", err)
		return nil, err
	}

	if oldAvatar != "" {
		helpers.DeleteFile(oldAvatar)
	}

	result := res.(*models.User)
	dto := dtos.ToUserDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "User Updated",
		Message: fmt.Sprintf("User updated: %s", req.Email),
		Data: map[string]interface{}{
			"id":    result.ID,
			"email": result.Email,
		},
	})

	// Invalidate cached session so next request gets updated permissions
	s.Access.Invalidate(id)

	s.Logger.LogEnd("UserUpdate", "User updated: %s (ID: %d)", dto.Email, dto.ID)
	return &dto, nil
}

// UserUpdateStatus updates a user's administrative status (Active/Suspend) — FSD §2.25.
// Independent of the Lock mechanism (§2.26/§2.27): a suspended user is not necessarily
// locked, and a locked user is not necessarily suspended.
func (s *Services) UserUpdateStatus(ctx context.Context, id uint, req dtos.UserStatusRequest) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserUpdateStatus", "Updating status for user ID %d to %s", id, req.Status)

	if _, err := s.repo.User.FindByID(nil, id); err != nil {
		s.Logger.LogEndWithError("UserUpdateStatus", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	result, err := s.repo.User.UpdateMap(nil, &models.User{ID: id}, map[string]interface{}{"status": req.Status})
	if err != nil {
		s.Logger.LogEndWithError("UserUpdateStatus", "Failed: %v", err)
		return nil, err
	}

	reloaded, err := s.repo.User.FindByID(nil, result.ID, "Roles")
	if err != nil {
		return nil, err
	}

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "User Status Updated",
		Message: fmt.Sprintf("User %s status changed to %s", reloaded.Email, req.Status),
		Data: map[string]interface{}{
			"id":     reloaded.ID,
			"status": req.Status,
		},
	})

	// Access cache (permissions/roles) isn't affected by status, but invalidate anyway so
	// a freshly-suspended user's next authenticated request re-evaluates from DB immediately.
	s.Access.Invalidate(id)

	s.Logger.LogEnd("UserUpdateStatus", "User %d status updated to %s", id, req.Status)
	dto := dtos.ToUserDTO(reloaded)
	return &dto, nil
}

// UserUnlock manually clears a user's account lock before it would naturally expire — FSD §2.27.
func (s *Services) UserUnlock(ctx context.Context, id uint) (*dtos.UserDTO, error) {
	s.Logger.LogStart("UserUnlock", "Unlocking user ID %d", id)

	existing, err := s.repo.User.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("UserUnlock", "User not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if existing.LockedUntil == nil || !existing.LockedUntil.After(time.Now()) {
		return nil, &helpers.CustomError{Status: 400, Message: "Akun ini tidak sedang terkunci"}
	}

	result, err := s.repo.User.UpdateMap(nil, &models.User{ID: id}, map[string]interface{}{
		"locked_until":          nil,
		"failed_login_attempts": 0,
	})
	if err != nil {
		s.Logger.LogEndWithError("UserUnlock", "Failed: %v", err)
		return nil, err
	}

	reloaded, err := s.repo.User.FindByID(nil, result.ID, "Roles")
	if err != nil {
		return nil, err
	}

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "User Unlocked",
		Message: fmt.Sprintf("Account unlocked: %s", reloaded.Email),
		Data: map[string]interface{}{
			"id": reloaded.ID,
		},
	})

	s.Logger.LogEnd("UserUnlock", "User %d unlocked", id)
	dto := dtos.ToUserDTO(reloaded)
	return &dto, nil
}

// UserDelete soft deletes a user and its role associations.
func (s *Services) UserDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("UserDelete", "Deleting user ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		user := models.User{ID: id}
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		_, err := s.repo.User.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("UserDelete", "Failed to delete user: %v", err)
		return err
	}

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "warning",
		Title:   "User Deleted",
		Message: fmt.Sprintf("User deleted: ID %d", id),
		Data: map[string]interface{}{
			"id": id,
		},
	})

	// Invalidate cached session
	s.Access.Invalidate(id)

	s.Logger.LogEnd("UserDelete", "User deleted: ID: %d", id)
	return nil
}
