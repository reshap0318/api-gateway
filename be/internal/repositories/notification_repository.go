package repositories

import (
	"time"

	"github.com/reshap0318/api-gateway/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	*GenericRepository[models.Notification]
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		GenericRepository: NewGenericRepository(db, &models.Notification{}),
	}
}

func (r *NotificationRepository) FindByUserIDWithOpts(userID uint, opts *QueryOptions) (*PagedResult[models.Notification], error) {
	db := r.getDB(nil)
	query := db.Where("user_id = ?", userID)

	if opts != nil {
		if opts.Page <= 0 {
			opts.Page = 1
		}
		if opts.PageSize <= 0 {
			opts.PageSize = 10
		}
		if opts.SortBy == "" {
			opts.SortBy = "created_at"
		}
		if opts.Order == "" {
			opts.Order = "DESC"
		}
	}

	return r.FindByUserIDWithOptsCustom(query, opts)
}

func (r *NotificationRepository) FindByUserIDAndFilters(userID uint, opts *QueryOptions, isRead *bool, notifType string) (*PagedResult[models.Notification], error) {
	db := r.getDB(nil)
	query := db.Where("user_id = ?", userID)

	if isRead != nil {
		if *isRead {
			query = query.Where("read_at IS NOT NULL")
		} else {
			query = query.Where("read_at IS NULL")
		}
	}

	if notifType != "" {
		query = query.Where("type = ?", notifType)
	}

	return r.FindByUserIDWithOptsCustom(query, opts)
}

func (r *NotificationRepository) FindByUserIDWithOptsCustom(query *gorm.DB, opts *QueryOptions) (*PagedResult[models.Notification], error) {
	var instance *models.Notification
	query = r.applyOptions(query, opts)

	var total int64
	if err := query.Model(&instance).Count(&total).Error; err != nil {
		return nil, err
	}

	page := 1
	pageSize := 10
	if opts != nil {
		if opts.Page > 0 {
			page = opts.Page
		}
		if opts.PageSize > 0 {
			pageSize = opts.PageSize
		}
	}

	if pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Limit(pageSize).Offset(offset)
	}

	datas := []models.Notification{}
	if err := query.Find(&datas).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &PagedResult[models.Notification]{
		Data:       datas,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *NotificationRepository) MarkAsRead(notificationID, userID uint) error {
	now := time.Now()
	return r.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", notificationID, userID).
		Update("read_at", now).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return r.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL AND deleted_at IS NULL", userID).
		Update("read_at", now).Error
}

func (r *NotificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}
