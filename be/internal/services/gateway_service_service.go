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

// GatewayServiceCreate creates a new upstream service registration.
func (s *Services) GatewayServiceCreate(ctx context.Context, req dtos.GatewayServiceRequest) (*dtos.GatewayServiceDTO, error) {
	s.Logger.LogStart("GatewayServiceCreate", "Creating gateway service: %s", req.Name)

	exists, err := s.repo.GatewayService.Exists(nil, map[string]interface{}{"name": req.Name})
	if err != nil {
		s.Logger.LogEndWithError("GatewayServiceCreate", "Failed to check name uniqueness: %v", err)
		return nil, err
	}
	if exists {
		return nil, &helpers.FieldError{Field: "name", Message: "Nama service sudah digunakan"}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	service := &models.GatewayService{
		Name:               req.Name,
		BaseURL:            req.BaseURL,
		Protocol:           req.Protocol,
		IsActive:           isActive,
		RateLimitPerMinute: req.RateLimitPerMinute,
		HealthStatus:       "unknown",
	}

	var result *models.GatewayService
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.GatewayService.Create(tx, service)
		return err
	})
	if err != nil {
		s.Logger.LogEndWithError("GatewayServiceCreate", "Failed: %v", err)
		return nil, err
	}

	s.RefreshRouteCache("GatewayServiceCreate")

	dtoForAudit := dtos.ToGatewayServiceDTO(result)
	s.recordAuditLog(ctx, "service", result.ID, "create", dtoForAudit)

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Gateway Service Created",
		Message: fmt.Sprintf("New gateway service registered: %s", result.Name),
		Data: map[string]interface{}{
			"id":   result.ID,
			"name": result.Name,
		},
	})

	s.Logger.LogEnd("GatewayServiceCreate", "Gateway service created: %s (ID: %d)", result.Name, result.ID)
	dto := dtos.ToGatewayServiceDTO(result)
	return &dto, nil
}

// GatewayServiceGetAllPaginated returns paginated gateway services.
func (s *Services) GatewayServiceGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions, search, protocol, isActive, healthStatus string) (*repositories.PagedResult[dtos.GatewayServiceDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	opts.Preloads = []string{"Routes"}

	if search != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic: "OR",
			Conditions: []repositories.QueryCondition{
				{Column: "name", Operator: "LIKE", Value: "%" + search + "%"},
			},
		})
	}
	if protocol != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "protocol", Operator: "=", Value: protocol}},
		})
	}
	if isActive != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "is_active", Operator: "=", Value: isActive == "true"}},
		})
	}
	if healthStatus != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "health_status", Operator: "=", Value: healthStatus}},
		})
	}

	result, err := s.repo.GatewayService.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.GatewayServiceDTO]{
		Data:       dtos.ToGatewayServiceDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// GatewayServiceGetByID returns a gateway service by ID with its routes.
func (s *Services) GatewayServiceGetByID(ctx context.Context, id uint) (*dtos.GatewayServiceDTO, error) {
	service, err := s.repo.GatewayService.FindByID(nil, id, "Routes")
	if err != nil {
		return nil, err
	}
	dto := dtos.ToGatewayServiceDTO(service)
	return &dto, nil
}

// GatewayServiceUpdate updates an existing gateway service.
func (s *Services) GatewayServiceUpdate(ctx context.Context, id uint, req dtos.GatewayServiceRequest) (*dtos.GatewayServiceDTO, error) {
	s.Logger.LogStart("GatewayServiceUpdate", "Updating gateway service ID: %d", id)

	existing, err := s.repo.GatewayService.FindByID(nil, id)
	if err != nil {
		s.Logger.LogEndWithError("GatewayServiceUpdate", "Not found: %v", err)
		return nil, err
	}

	if req.Name != "" && req.Name != existing.Name {
		exists, err := s.repo.GatewayService.Exists(nil, map[string]interface{}{"name": req.Name})
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, &helpers.FieldError{Field: "name", Message: "Nama service sudah digunakan"}
		}
	}

	isActive := existing.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Struct-based Update() drops zero-value fields (false, nil) from the UPDATE
	// statement, so is_active=false and a cleared rate_limit_per_minute would
	// silently fail to persist — UpdateMap sends every key explicitly instead.
	update := map[string]interface{}{
		"name":                  req.Name,
		"base_url":              req.BaseURL,
		"protocol":              req.Protocol,
		"rate_limit_per_minute": req.RateLimitPerMinute,
		"is_active":             isActive,
	}

	var result *models.GatewayService
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.GatewayService.UpdateMap(tx, &models.GatewayService{ID: id}, update)
		return err
	})
	if err != nil {
		s.Logger.LogEndWithError("GatewayServiceUpdate", "Failed: %v", err)
		return nil, err
	}

	s.RefreshRouteCache("GatewayServiceUpdate")

	s.recordAuditLog(ctx, "service", id, "update", map[string]interface{}{
		"before": dtos.ToGatewayServiceDTO(existing),
		"after":  dtos.ToGatewayServiceDTO(result),
	})

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Gateway Service Updated",
		Message: fmt.Sprintf("Gateway service updated: %s", result.Name),
		Data: map[string]interface{}{
			"id":   result.ID,
			"name": result.Name,
		},
	})

	s.Logger.LogEnd("GatewayServiceUpdate", "Gateway service updated: ID %d", id)
	dto := dtos.ToGatewayServiceDTO(result)
	return &dto, nil
}

// GatewayServiceDelete soft-deletes a gateway service, cascading to its routes.
func (s *Services) GatewayServiceDelete(ctx context.Context, id uint, cascade bool) error {
	s.Logger.LogStart("GatewayServiceDelete", "Deleting gateway service ID: %d", id)

	service, err := s.repo.GatewayService.FindByID(nil, id, "Routes")
	if err != nil {
		s.Logger.LogEndWithError("GatewayServiceDelete", "Not found: %v", err)
		return err
	}

	activeRoutes := 0
	for _, r := range service.Routes {
		if r.IsActive {
			activeRoutes++
		}
	}
	if activeRoutes > 0 && !cascade {
		return &helpers.CustomError{
			Status:  409,
			Message: fmt.Sprintf("Service ini masih punya %d route aktif, tetap hapus? (kirim ulang dengan cascade=true)", activeRoutes),
		}
	}

	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		if cascade {
			for _, r := range service.Routes {
				if _, err := s.repo.GatewayRoute.Delete(tx, r.ID); err != nil {
					return err
				}
			}
		}
		_, err := s.repo.GatewayService.Delete(tx, id)
		return err
	})
	if err != nil {
		s.Logger.LogEndWithError("GatewayServiceDelete", "Failed: %v", err)
		return err
	}

	s.RefreshRouteCache("GatewayServiceDelete")

	s.recordAuditLog(ctx, "service", id, "delete", dtos.ToGatewayServiceDTO(service))

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "warning",
		Title:   "Gateway Service Deleted",
		Message: fmt.Sprintf("Gateway service deleted: %s", service.Name),
		Data: map[string]interface{}{
			"id":   id,
			"name": service.Name,
		},
	})

	s.Logger.LogEnd("GatewayServiceDelete", "Gateway service deleted: ID %d", id)
	return nil
}
