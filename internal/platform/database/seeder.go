package database

import (
	"context"
	userDto "go_platform_template/internal/domain/user/dto"
	model "go_platform_template/internal/domain/user/model"
	userRepo "go_platform_template/internal/domain/user/repo"
	userService "go_platform_template/internal/domain/user/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SeedAdminUser seeds the system admin user
func SeedAdminUser(db *gorm.DB, logger *zap.SugaredLogger) {
	uRepo := userRepo.NewUserRepo(db)
	uService := userService.NewUserService(uRepo, nil)

	const (
		adminEmail    = "admin@example.com"
		adminUsername = "admin"
		adminPassword = "adminPassword"
	)

	// 1. Check if any admin exists
	var count int64
	if err := db.Model(&model.User{}).
		Where("user_type = ?", model.UserTypeAdmin).
		Count(&count).Error; err != nil {
		logger.Errorf("Failed to check existing admins: %v", err)
		return
	}

	if count > 0 {
		logger.Info("Admin user already exists. Skipping seeding.")
		return
	}

	// 2. Create admin DTO (plain password)
	adminReq := &userDto.UserCreateRequest{
		FirstName: "System",
		LastName:  "Administrator",
		Email:     adminEmail,
		Username:  adminUsername,
		Password:  adminPassword, // <-- plain text (Register will hash it)
		UserType:  string(model.UserTypeAdmin),
	}

	// 3. Register the admin with background context
	_, err := uService.Register(context.Background(), adminReq)
	if err != nil {
		logger.Errorf("Failed to seed admin user: %v", err)
		return
	}

	logger.Info("Admin user seeded successfully")
}
