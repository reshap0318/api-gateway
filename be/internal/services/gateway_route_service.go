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

// GatewayRouteCreate creates a new route under a service, with permission assignment.
func (s *Services) GatewayRouteCreate(ctx context.Context, req dtos.GatewayRouteRequest) (*dtos.GatewayRouteDTO, error) {
	s.Logger.LogStart("GatewayRouteCreate", "Creating route: %s %s (service %d)", req.Method, req.PathPattern, req.Service)

	if _, err := s.repo.GatewayService.FindByID(nil, req.Service); err != nil {
		return nil, &helpers.FieldError{Field: "service", Message: "Service tidak ditemukan"}
	}

	if err := helpers.ValidateRoutePathPattern(req.PathPattern); err != nil {
		return nil, err
	}

	exists, err := s.repo.GatewayRoute.Exists(nil, map[string]interface{}{
		"service_id":   req.Service,
		"method":       req.Method,
		"path_pattern": req.PathPattern,
	})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &helpers.FieldError{Field: "path_pattern", Message: "Route dengan method dan path ini sudah terdaftar untuk service ini"}
	}

	matchMode := req.PermissionMatchMode
	if matchMode == "" {
		matchMode = "any"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	route := &models.GatewayRoute{
		ServiceID:           req.Service,
		Method:              req.Method,
		PathPattern:         req.PathPattern,
		PermissionMatchMode: matchMode,
		RateLimitPerMinute:  req.RateLimitPerMinute,
		IsActive:            isActive,
	}

	var result *models.GatewayRoute
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.GatewayRoute.Create(tx, route)
		if err != nil {
			return nil, err
		}

		if len(req.Permissions) > 0 {
			var perms []models.Permission
			for _, permID := range req.Permissions {
				perms = append(perms, models.Permission{ID: permID})
			}
			if err := tx.Model(result).Association("Permissions").Append(perms); err != nil {
				return nil, err
			}
		}

		reloaded, err := s.repo.GatewayRoute.FindByID(tx, result.ID, "Permissions", "Service")
		if err != nil {
			return nil, err
		}
		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("GatewayRouteCreate", "Failed: %v", err)
		return nil, err
	}
	result = res.(*models.GatewayRoute)

	s.RefreshRouteCache("GatewayRouteCreate")

	s.recordAuditLog(ctx, "route", result.ID, "create", dtos.ToGatewayRouteDTO(result))

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Gateway Route Created",
		Message: fmt.Sprintf("New route created: %s %s", result.Method, result.PathPattern),
		Data: map[string]interface{}{
			"id":   result.ID,
			"path": result.PathPattern,
		},
	})

	s.Logger.LogEnd("GatewayRouteCreate", "Route created: ID %d", result.ID)
	dto := dtos.ToGatewayRouteDTO(result)
	return &dto, nil
}

// GatewayRouteGetAllPaginated returns paginated routes with filters.
func (s *Services) GatewayRouteGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions, serviceID, method, isActive string) (*repositories.PagedResult[dtos.GatewayRouteDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	opts.Preloads = []string{"Permissions", "Service"}

	if serviceID != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "service_id", Operator: "=", Value: serviceID}},
		})
	}
	if method != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "method", Operator: "=", Value: method}},
		})
	}
	if isActive != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "is_active", Operator: "=", Value: isActive == "true"}},
		})
	}

	result, err := s.repo.GatewayRoute.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.GatewayRouteDTO]{
		Data:       dtos.ToGatewayRouteDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// GatewayRouteGetByID returns a route by ID.
func (s *Services) GatewayRouteGetByID(ctx context.Context, id uint) (*dtos.GatewayRouteDTO, error) {
	route, err := s.repo.GatewayRoute.FindByID(nil, id, "Permissions", "Service")
	if err != nil {
		return nil, err
	}
	dto := dtos.ToGatewayRouteDTO(route)
	return &dto, nil
}

// GatewayRouteUpdate updates an existing route, replacing its permission assignment (full replace).
func (s *Services) GatewayRouteUpdate(ctx context.Context, id uint, req dtos.GatewayRouteRequest) (*dtos.GatewayRouteDTO, error) {
	s.Logger.LogStart("GatewayRouteUpdate", "Updating route ID: %d", id)

	existing, err := s.repo.GatewayRoute.FindByID(nil, id, "Permissions", "Service")
	if err != nil {
		s.Logger.LogEndWithError("GatewayRouteUpdate", "Not found: %v", err)
		return nil, err
	}
	beforeDTO := dtos.ToGatewayRouteDTO(existing)

	if req.PathPattern != "" && req.PathPattern != existing.PathPattern {
		if err := helpers.ValidateRoutePathPattern(req.PathPattern); err != nil {
			return nil, err
		}
	}

	matchMode := req.PermissionMatchMode
	if matchMode == "" {
		matchMode = existing.PermissionMatchMode
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Struct-based Update() drops zero-value fields (false, nil) from the UPDATE
	// statement, so is_active=false and a cleared rate_limit_per_minute would
	// silently fail to persist — UpdateMap sends every key explicitly instead.
	update := map[string]interface{}{
		"method":                req.Method,
		"path_pattern":          req.PathPattern,
		"permission_match_mode": matchMode,
		"rate_limit_per_minute": req.RateLimitPerMinute,
		"is_active":             isActive,
	}

	var result *models.GatewayRoute
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.GatewayRoute.UpdateMap(tx, &models.GatewayRoute{ID: id}, update)
		if err != nil {
			return nil, err
		}

		if err := tx.Model(result).Association("Permissions").Clear(); err != nil {
			return nil, err
		}
		if len(req.Permissions) > 0 {
			var perms []models.Permission
			for _, permID := range req.Permissions {
				perms = append(perms, models.Permission{ID: permID})
			}
			if err := tx.Model(result).Association("Permissions").Append(perms); err != nil {
				return nil, err
			}
		}

		reloaded, err := s.repo.GatewayRoute.FindByID(tx, result.ID, "Permissions", "Service")
		if err != nil {
			return nil, err
		}
		return reloaded, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("GatewayRouteUpdate", "Failed: %v", err)
		return nil, err
	}
	result = res.(*models.GatewayRoute)

	s.RefreshRouteCache("GatewayRouteUpdate")

	s.recordAuditLog(ctx, "route", id, "update", map[string]interface{}{
		"before": beforeDTO,
		"after":  dtos.ToGatewayRouteDTO(result),
	})

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Gateway Route Updated",
		Message: fmt.Sprintf("Route updated: %s %s", result.Method, result.PathPattern),
		Data: map[string]interface{}{
			"id":   result.ID,
			"path": result.PathPattern,
		},
	})

	s.Logger.LogEnd("GatewayRouteUpdate", "Route updated: ID %d", id)
	dto := dtos.ToGatewayRouteDTO(result)
	return &dto, nil
}

// GatewayRouteDelete soft-deletes a route.
func (s *Services) GatewayRouteDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("GatewayRouteDelete", "Deleting route ID: %d", id)

	route, err := s.repo.GatewayRoute.FindByID(nil, id, "Permissions", "Service")
	if err != nil {
		s.Logger.LogEndWithError("GatewayRouteDelete", "Not found: %v", err)
		return err
	}

	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GatewayRoute{ID: id}).Association("Permissions").Clear(); err != nil {
			return err
		}
		_, err := s.repo.GatewayRoute.Delete(tx, id)
		return err
	})
	if err != nil {
		s.Logger.LogEndWithError("GatewayRouteDelete", "Failed: %v", err)
		return err
	}

	s.RefreshRouteCache("GatewayRouteDelete")

	s.recordAuditLog(ctx, "route", id, "delete", dtos.ToGatewayRouteDTO(route))

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "warning",
		Title:   "Gateway Route Deleted",
		Message: fmt.Sprintf("Route deleted: %s %s", route.Method, route.PathPattern),
		Data: map[string]interface{}{
			"id": id,
		},
	})

	s.Logger.LogEnd("GatewayRouteDelete", "Route deleted: ID %d", id)
	return nil
}
