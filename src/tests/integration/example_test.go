package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// Example integration test - replace with your actual tests
// These tests require a running database or use test containers

func TestHealthEndpoint(t *testing.T) {
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a test Fiber app
	app := fiber.New()

	// Add health endpoint
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)

	// Perform request
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	// Assert status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// Example of a test with setup and teardown
func TestWithDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// TODO: Setup test database connection
	// db := setupTestDB(t)
	// defer teardownTestDB(db)

	t.Run("create entity", func(t *testing.T) {
		// Test creation
	})

	t.Run("read entity", func(t *testing.T) {
		// Test reading
	})

	t.Run("update entity", func(t *testing.T) {
		// Test updating
	})

	t.Run("delete entity", func(t *testing.T) {
		// Test deletion
	})
}

// Helper functions for database setup
// func setupTestDB(t *testing.T) *gorm.DB {
// 	// Connect to test database
// 	// Run migrations
// 	return db
// }

// func teardownTestDB(db *gorm.DB) {
// 	// Clean up test data
// 	// Close connection
// }
