package main

import (
	"fmt"
	"log"
	"os"

	"vaxen/api/internal/config"
	"vaxen/api/internal/database"
	"vaxen/api/internal/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Usage:
//   go run cmd/seed/main.go admin <email> <password>
//
// Creates a platform admin user (not tied to any organization).
// Run this once after deploying to create your first admin.
func main() {
	if len(os.Args) < 4 || os.Args[1] != "admin" {
		fmt.Println("Usage: go run cmd/seed/main.go admin <email> <password>")
		fmt.Println("")
		fmt.Println("Creates a platform admin user who can approve access requests,")
		fmt.Println("manage settings, feature flags, and exchange rates.")
		os.Exit(1)
	}

	email := os.Args[2]
	password := os.Args[3]

	if len(password) < 8 {
		log.Fatal("Password must be at least 8 characters")
	}

	godotenv.Load()
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Check if already exists
	var existing models.User
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		if existing.Role == models.UserRoleAdmin {
			log.Fatalf("Admin user %s already exists", email)
		}
		// Promote existing user to admin
		database.DB.Model(&existing).Update("role", models.UserRoleAdmin)
		fmt.Printf("Promoted existing user %s to admin\n", email)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Admin users don't belong to an organization — create a placeholder org
	org := models.Organization{
		Name:               "Vaxen Platform",
		LegalName:          "Vaxen Platform",
		RegistrationNumber: "PLATFORM",
		Country:            "US",
		KYBStatus:          models.KybStatusApproved,
	}
	if err := database.DB.Create(&org).Error; err != nil {
		log.Fatalf("Failed to create platform org: %v", err)
	}

	admin := models.User{
		Email:          email,
		FirstName:      "Platform",
		LastName:       "Admin",
		PasswordHash:   string(hash),
		Role:           models.UserRoleAdmin,
		OrganizationID: org.ID,
		IsDirector:     false,
		IsActive:       true,
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("Admin user created successfully\n")
	fmt.Printf("  Email: %s\n", email)
	fmt.Printf("  Role:  admin\n")
	fmt.Printf("  ID:    %s\n", admin.ID)
	fmt.Println("")
	fmt.Println("You can now login at POST /api/v1/auth/login with these credentials.")
}
