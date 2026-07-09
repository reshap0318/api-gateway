package repositories

import "gorm.io/gorm"

type Repositories struct {
	TxManager       *TransactionManager
	User            *UserRepository
	PasswordReset   *PasswordResetRepository
	Permission      *PermissionRepository
	Role            *RoleRepository
	RoleHasPerm     *RoleHasPermissionRepository
	UserRole        *UserRoleRepository
	Notification    *NotificationRepository
	GatewayService  *GatewayServiceRepository
	GatewayRoute    *GatewayRouteRepository
	GatewayAuditLog *GatewayAuditLogRepository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	txManager := NewTransactionManager(db)
	userRepo := NewUserRepository(db)
	passwordResetRepo := NewPasswordResetRepository(db)
	permissionRepo := NewPermissionRepository(db)
	roleRepo := NewRoleRepository(db)
	roleHasPermRepo := NewRoleHasPermissionRepository(db)
	userRoleRepo := NewUserRoleRepository(db)
	notificationRepo := NewNotificationRepository(db)
	gatewayServiceRepo := NewGatewayServiceRepository(db)
	gatewayRouteRepo := NewGatewayRouteRepository(db)
	gatewayAuditLogRepo := NewGatewayAuditLogRepository(db)

	return &Repositories{
		TxManager:       txManager,
		User:            userRepo,
		PasswordReset:   passwordResetRepo,
		Permission:      permissionRepo,
		Role:            roleRepo,
		RoleHasPerm:     roleHasPermRepo,
		UserRole:        userRoleRepo,
		Notification:    notificationRepo,
		GatewayService:  gatewayServiceRepo,
		GatewayRoute:    gatewayRouteRepo,
		GatewayAuditLog: gatewayAuditLogRepo,
	}, nil
}
