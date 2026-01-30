package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/yourusername/online-library/internal/config"
	"github.com/yourusername/online-library/internal/infrastructure/db/postgres"
	"github.com/yourusername/online-library/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Connect to database
	conn, err := postgres.NewConnection(context.Background(), cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer conn.Close()

	// Get admin credentials from environment
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@amarpathagar.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	adminFullName := os.Getenv("ADMIN_FULL_NAME")
	if adminFullName == "" {
		adminFullName = "System Administrator"
	}

	// Check if admin user already exists
	userRepo := repository.NewUserRepository(conn.DB, logger)
	existingUser, err := userRepo.FindByUsername(context.Background(), adminUsername)
	if err == nil && existingUser != nil {
		fmt.Printf("✓ Admin user '%s' already exists\n", adminUsername)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("failed to hash password:", err)
	}

	// Create admin user
	query := `
		INSERT INTO users (
			id, username, email, password_hash, full_name, role, 
			success_score, books_shared, books_received, 
			created_at, updated_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, 'admin', 
			100, 0, 0, 
			NOW(), NOW()
		)
	`

	_, err = conn.DB.ExecContext(
		context.Background(),
		query,
		adminUsername,
		adminEmail,
		string(hashedPassword),
		adminFullName,
	)

	if err != nil {
		log.Fatal("failed to create admin user:", err)
	}

	fmt.Println("✓ Admin user created successfully!")
	fmt.Printf("  Username: %s\n", adminUsername)
	fmt.Printf("  Email: %s\n", adminEmail)
	fmt.Printf("  Password: %s\n", adminPassword)
	fmt.Println("\n⚠️  IMPORTANT: Change the admin password after first login!")
}
