package services

import (
	"context"
	"encoding/json"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
	"github.com/reshap0318/api-gateway/internal/repositories"
)

// recordAuditLog persists one immutable audit trail entry for a Service/Route config
// change (FSD §2.23). Called AFTER the CUD transaction has committed — a failure here must
// never roll back or fail the operation that already succeeded, so errors are only logged.
//
// changes payload convention: create/delete → full entity snapshot; update →
// {"before": ..., "after": ...} (only meaningful when the caller wants full-record diffing;
// callers here pass DTOs for readability).
func (s *Services) recordAuditLog(ctx context.Context, entityType string, entityID uint, action string, changes interface{}) {
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		s.Logger.LogWarn("recordAuditLog", "Failed to marshal changes for %s %d: %v", entityType, entityID, err)
		changesJSON = []byte("{}")
	}

	entry := &models.GatewayAuditLog{
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		ActorUserID: helpers.GetCallerID(ctx),
		Changes:     string(changesJSON),
	}

	if _, err := s.repo.GatewayAuditLog.Create(nil, entry); err != nil {
		s.Logger.LogWarn("recordAuditLog", "Failed to record audit log for %s %d: %v", entityType, entityID, err)
	}
}

// GatewayAuditLogGetAllPaginated returns paginated audit log entries with filters (FSD §2.24).
func (s *Services) GatewayAuditLogGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions, entityType, entity, actor, from, to string) (*repositories.PagedResult[dtos.GatewayAuditLogDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}
	opts.Preloads = []string{"Actor"}

	if entityType != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "entity_type", Operator: "=", Value: entityType}},
		})
	}
	if entity != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "entity_id", Operator: "=", Value: entity}},
		})
	}
	if actor != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "actor_user_id", Operator: "=", Value: actor}},
		})
	}
	if from != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "created_at", Operator: ">=", Value: from}},
		})
	}
	if to != "" {
		opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
			Logic:      "AND",
			Conditions: []repositories.QueryCondition{{Column: "created_at", Operator: "<=", Value: to}},
		})
	}

	result, err := s.repo.GatewayAuditLog.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.GatewayAuditLogDTO]{
		Data:       dtos.ToGatewayAuditLogDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// GatewayAuditLogGetByID returns a single audit log entry by ID.
func (s *Services) GatewayAuditLogGetByID(ctx context.Context, id uint) (*dtos.GatewayAuditLogDTO, error) {
	entry, err := s.repo.GatewayAuditLog.FindByID(nil, id, "Actor")
	if err != nil {
		return nil, err
	}
	dto := dtos.ToGatewayAuditLogDTO(entry)
	return &dto, nil
}
