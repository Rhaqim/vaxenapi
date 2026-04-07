package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"vaxen/api/internal/config"
	"vaxen/api/internal/database"
	"vaxen/api/internal/models"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

// Usage:
//
//	go run cmd/seed/main.go admin <email> <password>
//	go run cmd/seed/main.go exchange-rates --all
//	go run cmd/seed/main.go exchange-rates USD EUR
//	go run cmd/seed/main.go exchange-rates USD EUR,GBP,JPY,NGN
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
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

	switch os.Args[1] {
	case "admin":
		seedAdmin()
	case "exchange-rates":
		seedExchangeRates()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Vaxen Seed Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/seed/main.go admin <email> <password>")
	fmt.Println("  go run cmd/seed/main.go exchange-rates --all")
	fmt.Println("  go run cmd/seed/main.go exchange-rates <base> <targets>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  admin             Create a platform admin user")
	fmt.Println("  exchange-rates    Seed exchange rates from public API data")
	fmt.Println("")
	fmt.Println("Exchange rate examples:")
	fmt.Println("  exchange-rates --all                    Seed all common fiat + crypto pairs")
	fmt.Println("  exchange-rates USD EUR                  Seed USD/EUR pair")
	fmt.Println("  exchange-rates USD EUR,GBP,JPY,NGN      Seed USD to multiple targets")
}

// ==================== Admin Seed ====================

func seedAdmin() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run cmd/seed/main.go admin <email> <password>")
		os.Exit(1)
	}

	email := os.Args[2]
	password := os.Args[3]

	if len(password) < 8 {
		log.Fatal("Password must be at least 8 characters")
	}

	var existing models.User
	if err := database.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		if existing.Role == models.UserRoleAdmin {
			log.Fatalf("Admin user %s already exists", email)
		}
		database.DB.Model(&existing).Update("role", models.UserRoleAdmin)
		fmt.Printf("Promoted existing user %s to admin\n", email)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

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

	fmt.Printf("Admin user created: %s (ID: %s)\n", email, admin.ID)
}

// ==================== Exchange Rate Seed ====================

// Default currency sets for --all
var (
	defaultFiatBases   = []string{"USD", "EUR", "GBP"}
	defaultFiatTargets = []string{"USD", "EUR", "GBP", "JPY", "CHF", "CAD", "AUD", "NGN", "KES", "ZAR", "BRL", "MXN", "INR", "CNY"}
)

func seedExchangeRates() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  exchange-rates --all")
		fmt.Println("  exchange-rates <base> <targets>")
		os.Exit(1)
	}

	if os.Args[2] == "--all" {
		seedAllRates()
		return
	}

	base := strings.ToUpper(os.Args[2])
	if len(os.Args) < 4 {
		fmt.Println("Usage: exchange-rates <base> <target1,target2,...>")
		os.Exit(1)
	}

	targets := strings.Split(strings.ToUpper(os.Args[3]), ",")
	fetchAndSeedRates(base, targets)
}

func seedAllRates() {
	fmt.Println("Seeding all common exchange rates...")

	for _, base := range defaultFiatBases {
		targets := make([]string, 0)
		for _, t := range defaultFiatTargets {
			if t != base {
				targets = append(targets, t)
			}
		}
		fetchAndSeedRates(base, targets)
	}

	fmt.Println("Done seeding exchange rates.")
}

func fetchAndSeedRates(base string, targets []string) {
	fmt.Printf("Fetching rates for %s -> %s...\n", base, strings.Join(targets, ","))

	rates, err := fetchRatesFromPublicAPI(base, targets)
	if err != nil {
		log.Printf("Warning: failed to fetch rates for %s: %v (using fallback 1.0)", base, err)
		// Insert with rate 1.0 so the pair exists, admin can update later
		for _, target := range targets {
			upsertRate(base, target, decimal.NewFromInt(1), decimal.NewFromFloat(0.005))
		}
		return
	}

	for _, target := range targets {
		rate, ok := rates[target]
		if !ok {
			log.Printf("  %s/%s: not found in API response, skipping", base, target)
			continue
		}
		spread := decimal.NewFromFloat(0.005) // 0.5% default spread
		upsertRate(base, target, rate, spread)
	}
}

func upsertRate(from, to string, rate, spread decimal.Decimal) {
	var existing models.ExchangeRate
	err := database.DB.Where("from_currency = ? AND to_currency = ?", from, to).First(&existing).Error
	if err == nil {
		database.DB.Model(&existing).Updates(map[string]any{
			"rate":       rate,
			"spread":     spread,
			"is_active":  true,
			"updated_by": "seed-tool",
		})
		fmt.Printf("  %s/%s = %s (updated)\n", from, to, rate.StringFixed(6))
	} else {
		newRate := models.ExchangeRate{
			FromCurrency: from,
			ToCurrency:   to,
			Rate:         rate,
			Spread:       spread,
			IsActive:     true,
			UpdatedBy:    "seed-tool",
		}
		database.DB.Create(&newRate)
		fmt.Printf("  %s/%s = %s (created)\n", from, to, rate.StringFixed(6))
	}
}

// fetchRatesFromPublicAPI fetches exchange rates from the free Frankfurt API.
// https://api.frankfurter.app — no API key required, ECB data.
func fetchRatesFromPublicAPI(base string, targets []string) (map[string]decimal.Decimal, error) {
	url := fmt.Sprintf("https://api.frankfurter.app/latest?from=%s&to=%s",
		base, strings.Join(targets, ","))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Base  string                 `json:"base"`
		Rates map[string]float64     `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	rates := make(map[string]decimal.Decimal, len(result.Rates))
	for currency, rate := range result.Rates {
		rates[currency] = decimal.NewFromFloat(rate)
	}

	return rates, nil
}
