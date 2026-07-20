package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// PermissionCreate creates a new permission.
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
		if err != nil {
			return err
		}

		return err
	}); err != nil {
		s.Logger.LogEndWithError("PermissionCreate", "Failed to create permission: %v", err)
		return nil, err
	}

	dto := dtos.ToPermissionDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Permission Created",
		Message: fmt.Sprintf("New permission created: %s", req.Name),
		Data: map[string]interface{}{
			"id":   result.ID,
			"name": result.Name,
		},
	})

	s.Logger.LogEnd("PermissionCreate", "Permission created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// PermissionGetAll returns permissions. opts.PageSize negative returns all records
// unpaginated; unset/positive paginates (default page size 10) — either way the response
// carries pagination metadata.
func (s *Services) PermissionGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.PermissionDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "ASC"
	}

	if !s.Access.HasPermission(ctx, "role.index-su") {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic: "AND",
			Conditions: []repositories.QueryCondition{
				{Column: "id", Operator: ">", Value: 23},
			},
		})
	}

	result, err := s.repo.Permission.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	dtoList := dtos.ToPermissionDTOList(result.Data)

	return &repositories.PagedResult[dtos.PermissionDTO]{
		Data:       dtoList,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// PermissionGetByID returns a permission by ID.
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
	if id <= 23 && !s.Access.HasPermission(ctx, "role.index-su") {
		return nil, helpers.ErrForbidden
	}

	permission, err := s.repo.Permission.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToPermissionDTO(permission)
	return &dto, nil
}

// PermissionUpdate updates an existing permission.
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
		if err != nil {
			return err
		}

		return err
	}); err != nil {
		s.Logger.LogEndWithError("PermissionUpdate", "Failed to update permission: %v", err)
		return nil, err
	}

	// Access cache keys permissions by Name — a rename invalidates every user's cached set.
	if req.Name != "" {
		s.Access.InvalidateAll()
	}

	dto := dtos.ToPermissionDTO(result)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Permission Updated",
		Message: fmt.Sprintf("Permission updated: %s", result.Name),
		Data: map[string]interface{}{
			"id":   result.ID,
			"name": result.Name,
		},
	})

	s.Logger.LogEnd("PermissionUpdate", "Permission updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// PermissionDelete soft deletes a permission.
func (s *Services) PermissionDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("PermissionDelete", "Deleting permission ID: %d", id)

	var permission *models.Permission
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		permission, err = s.repo.Permission.FindByID(tx, id)
		if err != nil {
			return err
		}
		_, err = s.repo.Permission.Delete(tx, id)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		s.Logger.LogEndWithError("PermissionDelete", "Failed to delete permission: %v", err)
		return err
	}

	// Permission is gone — every user who had it cached still "has" it until cleared.
	s.Access.InvalidateAll()

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "warning",
		Title:   "Permission Deleted",
		Message: fmt.Sprintf("Permission deleted: %s", permission.Name),
		Data: map[string]interface{}{
			"id": permission.ID,
		},
	})

	s.Logger.LogEnd("PermissionDelete", "Permission deleted: ID: %d", id)
	return nil
}
