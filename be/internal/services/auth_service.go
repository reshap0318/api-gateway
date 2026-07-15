package services

import (
	"context"
	"fmt"
	"time"

	"github.com/reshap0318/api-gateway/internal/dtos"
	"github.com/reshap0318/api-gateway/internal/helpers"
	"github.com/reshap0318/api-gateway/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthValidateToken validates a JWT token and returns the claims.
func (s *Services) AuthValidateToken(tokenString string) (*helpers.JWTClaims, error) {
	claims, err := helpers.ValidateToken(tokenString, s.JWKSManager.GetPublicKey())
	if err != nil {
		return nil, helpers.ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, helpers.ErrExpiredToken
	}

	// Check token blacklist — reject if jti has been revoked (logout)
	if s.RedisClient.IsCacheAvailable() {
		blacklistKey := fmt.Sprintf("blacklist:jti:%s", claims.ID)
		if blacklisted, err := s.RedisClient.Exists(blacklistKey); err == nil && blacklisted {
			return nil, helpers.ErrInvalidToken
		}
	}

	// Try to get cached session from Redis with fallback to DB
	if s.RedisClient.IsCacheAvailable() {
		sessionKey := fmt.Sprintf("session:%d", claims.UserID)
		var cachedUserDTO dtos.UserDTO
		if err := s.RedisClient.GetJSON(sessionKey, &cachedUserDTO); err == nil {
			return claims, nil
		} else {
			s.Logger.LogWarn("AuthValidateToken", "Cache miss/error for session:%d, falling back to DB: %v", claims.UserID, err)
		}
	}

	// Fallback: validate user from database
	_, err = s.repo.User.FindByID(nil, claims.UserID)
	if err != nil {
		s.Logger.LogStep("AuthValidateToken", "User not found in DB: %d", claims.UserID)
		return nil, helpers.ErrInvalidCredential
	}

	return claims, nil
}

// AuthLogin authenticates a user and returns tokens.
func (s *Services) AuthLogin(ctx context.Context, email, password string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("AuthLogin", "User login attempt: %s", email)

	user, err := s.repo.User.FindByEmail(nil, email)
	if err != nil {
		s.Logger.LogStep("AuthLogin", "User not found: %s", email)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - user not found")
		return nil, helpers.ErrInvalidCredential
	}
	s.Logger.LogStep("AuthLogin", "User found: %s", email)

	// Status check (FSD §2.25) — administrative suspend, checked before credential validation.
	if user.Status == "suspended" {
		s.Logger.LogStep("AuthLogin", "Account suspended: %s", email)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - account suspended")
		return nil, &helpers.CustomError{Status: 401, Message: "Akun tidak aktif, hubungi admin"}
	}

	// Lock check (FSD §2.26) — automatic security lock, independent of Status.
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		s.Logger.LogStep("AuthLogin", "Account locked: %s", email)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - account locked")
		return nil, &helpers.CustomError{Status: 401, Message: "Akun terkunci sementara, coba lagi nanti"}
	}

	if !s.checkPassword(password, user.Password) {
		s.Logger.LogStep("AuthLogin", "Invalid password for user: %s", email)
		s.registerFailedLogin(user)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - invalid password")
		return nil, helpers.ErrInvalidCredential
	}
	s.Logger.LogStep("AuthLogin", "Password validated successfully")

	if user.FailedLoginAttempts > 0 {
		if _, err := s.repo.User.UpdateMap(nil, &models.User{ID: user.ID}, map[string]interface{}{"failed_login_attempts": 0}); err != nil {
			s.Logger.LogWarn("AuthLogin", "Failed to reset failed_login_attempts for user %d: %v", user.ID, err)
		}
	}

	token, err := s.generateTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthLogin", "Failed to generate token: %v", err)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthLogin", "Access token generated")

	refreshToken, err := s.generateRefreshTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthLogin", "Failed to generate refresh token: %v", err)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - refresh token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthLogin", "Refresh token generated")

	// Reload user with roles for response
	userWithRoles, err := s.repo.User.FindByID(nil, user.ID, "Roles.Permissions")
	if err != nil {
		s.Logger.LogWarn("AuthLogin", "Failed to load user roles: %v", err)
		userWithRoles = user
	}

	// Cache UserDTO to Redis with fallback
	if s.RedisClient.IsCacheAvailable() {
		userDTO := dtos.ToUserDTO(userWithRoles)
		sessionKey := fmt.Sprintf("session:%d", userWithRoles.ID)
		if err := s.RedisClient.SetJSON(sessionKey, userDTO, s.cfg.Expiration); err != nil {
			s.Logger.LogWarn("AuthLogin", "Failed to cache session to Redis: %v", err)
		} else {
			s.Logger.LogStep("AuthLogin", "Session cached to Redis")
		}
	}

	s.Logger.LogEnd("AuthLogin", "Login successful for user: %s", email)
	return &dtos.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         dtos.ToUserDTO(userWithRoles),
	}, nil
}

// AuthRefreshToken refreshes the access token using a refresh token.
func (s *Services) AuthRefreshToken(ctx context.Context, refreshToken string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("AuthRefreshToken", "Token refresh attempt")

	claims, err := s.AuthValidateToken(refreshToken)
	if err != nil {
		s.Logger.LogStep("AuthRefreshToken", "Invalid refresh token: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - invalid token")
		return nil, err
	}
	s.Logger.LogStep("AuthRefreshToken", "Refresh token validated")

	user, err := s.repo.User.FindByID(nil, claims.UserID)
	if err != nil {
		s.Logger.LogStep("AuthRefreshToken", "User not found: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - user not found")
		return nil, helpers.ErrInvalidCredential
	}
	s.Logger.LogStep("AuthRefreshToken", "User found: %s", user.Email)

	token, err := s.generateTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthRefreshToken", "Failed to generate token: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthRefreshToken", "Access token regenerated")

	newRefreshToken, err := s.generateRefreshTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthRefreshToken", "Failed to generate refresh token: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - refresh token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthRefreshToken", "Refresh token regenerated")

	s.Logger.LogEnd("AuthRefreshToken", "Token refreshed successfully for user: %s", user.Email)
	return &dtos.LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User:         dtos.ToUserDTO(user),
	}, nil
}

// AuthLogout blacklists the token's jti in Redis so it cannot be reused.
func (s *Services) AuthLogout(ctx context.Context, tokenString string) error {
	s.Logger.LogStart("AuthLogout", "Logout request")

	claims, err := helpers.ValidateToken(tokenString, s.JWKSManager.GetPublicKey())
	if err != nil {
		s.Logger.LogEndWithError("AuthLogout", "Invalid token: %v", err)
		return helpers.ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		// Token already expired — nothing to blacklist
		s.Logger.LogEnd("AuthLogout", "Token already expired, skipping blacklist")
		return nil
	}

	if s.RedisClient.IsCacheAvailable() {
		blacklistKey := fmt.Sprintf("blacklist:jti:%s", claims.ID)
		if err := s.RedisClient.Set(blacklistKey, "1", ttl); err != nil {
			s.Logger.LogWarn("AuthLogout", "Failed to blacklist jti %s: %v", claims.ID, err)
		} else {
			s.Logger.LogStep("AuthLogout", "jti blacklisted: %s (TTL: %s)", claims.ID, ttl)
		}
	} else {
		s.Logger.LogWarn("AuthLogout", "Redis unavailable — token not blacklisted")
	}

	s.Logger.LogEnd("AuthLogout", "Logout successful for user: %d", claims.UserID)
	return nil
}

// Helper functions

func (s *Services) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// registerFailedLogin increments the failed-login counter and locks the account once the
// configured threshold is reached (FSD §2.26). Errors are logged only — a failure here must
// never surface as a different error to the caller than the plain "invalid credential".
func (s *Services) registerFailedLogin(user *models.User) {
	maxAttempts := helpers.GetEnvInt("AUTH_MAX_FAILED_LOGIN_ATTEMPTS", 5)
	lockDurationMin := helpers.GetEnvInt("AUTH_LOCK_DURATION_MINUTES", 15)

	attempts := user.FailedLoginAttempts + 1
	update := map[string]interface{}{"failed_login_attempts": attempts}

	if attempts >= maxAttempts {
		lockedUntil := time.Now().Add(time.Duration(lockDurationMin) * time.Minute)
		update["locked_until"] = lockedUntil
		s.Logger.LogWarn("registerFailedLogin", "User %d locked until %s after %d failed attempts", user.ID, lockedUntil, attempts)
	}

	if _, err := s.repo.User.UpdateMap(nil, &models.User{ID: user.ID}, update); err != nil {
		s.Logger.LogWarn("registerFailedLogin", "Failed to update failed_login_attempts for user %d: %v", user.ID, err)
	}
}

// getUserRolesAndPermissions fetches role names and permission names for a user
func (s *Services) getUserRolesAndPermissions(userID uint) (roles []string, permissions []string) {
	user, err := s.repo.User.FindByID(nil, userID, "Roles.Permissions")
	if err != nil {
		return []string{}, []string{}
	}

	roleSet := make(map[string]bool)
	permSet := make(map[string]bool)

	for _, r := range user.Roles {
		roleSet[r.Name] = true
		for _, p := range r.Permissions {
			permSet[p.Name] = true
		}
	}

	for name := range roleSet {
		roles = append(roles, name)
	}
	for name := range permSet {
		permissions = append(permissions, name)
	}

	if roles == nil {
		roles = []string{}
	}
	if permissions == nil {
		permissions = []string{}
	}

	return roles, permissions
}

func (s *Services) generateTokenWithClaims(user *models.User) (string, error) {
	roles, permissions := s.getUserRolesAndPermissions(user.ID)
	return helpers.GenerateToken(
		user.ID,
		user.Email,
		user.Name,
		roles,
		permissions,
		s.JWKSManager.GetPrivateKey(),
		s.JWKSManager.GetKeyID(),
		helpers.GetEnvInt("JWT_EXPIRATION", 24),
	)
}

func (s *Services) generateRefreshTokenWithClaims(user *models.User) (string, error) {
	roles, permissions := s.getUserRolesAndPermissions(user.ID)
	return helpers.GenerateRefreshToken(
		user.ID,
		user.Email,
		user.Name,
		roles,
		permissions,
		s.JWKSManager.GetPrivateKey(),
		s.JWKSManager.GetKeyID(),
		helpers.GetEnvInt("JWT_REFRESH_EXPIRATION", 168),
	)
}

// AuthForgetPassword generates a reset token and sends it via email.
func (s *Services) AuthForgetPassword(ctx context.Context, email string) error {
	s.Logger.LogStart("AuthForgetPassword", "Reset password request for: %s", email)

	user, err := s.repo.User.FindByEmail(nil, email)
	if err != nil {
		s.Logger.LogStep("AuthForgetPassword", "User not found: %s", email)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - user not found")
		return helpers.ErrNotFound
	}

	s.Logger.LogStep("AuthForgetPassword", "User found: %s", email)

	frontendURL := helpers.GetEnv("APP_FE_URL", "http://localhost:3000")

	token, err := helpers.GenerateRandomString(32)
	if err != nil {
		s.Logger.LogError("AuthForgetPassword", "Failed to generate token: %v", err)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - token generation error")
		return err
	}

	s.Logger.LogStep("AuthForgetPassword", "Reset token generated")

	hashedToken, err := helpers.HashString(token)
	if err != nil {
		s.Logger.LogError("AuthForgetPassword", "Failed to hash token: %v", err)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - token hashing error")
		return err
	}

	s.Logger.LogStep("AuthForgetPassword", "Token hashed successfully")

	expiresAt := time.Now().Add(1 * time.Hour)

	passwordReset := &models.PasswordReset{
		Email:     user.Email,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		if _, err := s.repo.PasswordReset.Create(tx, passwordReset); err != nil {
			return err
		}
		return nil
	}); err != nil {
		s.Logger.LogError("AuthForgetPassword", "Failed to save reset token: %v", err)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - database error")
		return err
	}

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "info",
		Title:   "Password Reset Requested",
		Message: fmt.Sprintf("Password reset requested for %s", user.Email),
		Data: map[string]interface{}{
			"email": user.Email,
		},
	})

	s.Logger.LogStep("AuthForgetPassword", "Reset token saved to database")

	resetURL := frontendURL + "/reset-password?token=" + token

	// Async email sending dengan logger
	go func() {
		if err := s.EmailClient.SendResetPasswordEmail(user.Email, token, resetURL); err != nil {
			s.Logger.LogError("AuthForgetPassword", "Failed to send reset email to %s: %v", user.Email, err)
		} else {
			s.Logger.LogStep("AuthForgetPassword", "Reset email sent successfully to %s", user.Email)
		}
	}()

	s.Logger.LogEnd("AuthForgetPassword", "Reset password request processed successfully")

	return nil
}

// AuthResetPassword validates token and resets user password.
func (s *Services) AuthResetPassword(ctx context.Context, token, newPassword string) error {
	s.Logger.LogStart("AuthResetPassword", "Password reset attempt")

	// Hash token to find in database
	hashedToken, err := helpers.HashString(token)
	if err != nil {
		return helpers.ErrTokenInvalid
	}

	var resetEmail string
	_, err = s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		// Find reset record by token (hashed)
		reset, err := s.repo.PasswordReset.FindByToken(tx, hashedToken)
		if err != nil {
			return nil, helpers.ErrTokenInvalid
		}

		// Check if token is expired
		if reset.ExpiresAt.Before(time.Now()) {
			return nil, helpers.ErrTokenExpired
		}

		// Check if token already used
		if reset.Used {
			return nil, helpers.ErrTokenUsed
		}

		// Find user by email
		if _, err := s.repo.User.FindByEmail(tx, reset.Email); err != nil {
			return nil, helpers.ErrNotFound
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		// Update user password using generic Update
		if _, err := s.repo.User.Update(tx, &models.User{Email: reset.Email}, &models.User{Password: string(hashedPassword)}); err != nil {
			return nil, err
		}

		// Invalidate the token using generic Update
		if _, err := s.repo.PasswordReset.Update(tx, &models.PasswordReset{Token: reset.Token}, &models.PasswordReset{Used: true}); err != nil {
			return nil, err
		}

		resetEmail = reset.Email
		return nil, nil
	})

	if err != nil {
		s.Logger.LogEndWithError("AuthResetPassword", "Password reset failed: %v", err)
		return err
	}

	_ = s.NotificationCreate(ctx, &NotificationCreateParams{
		Type:    "success",
		Title:   "Password Reset Successful",
		Message: "Your password has been reset successfully",
		Data: map[string]interface{}{
			"email": resetEmail,
		},
	})

	s.Logger.LogEnd("AuthResetPassword", "Password reset successful")
	return nil
}
