// handlers.go
package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/saurabhag23/receipt-processor/internal/models"
	"github.com/saurabhag23/receipt-processor/internal/utils" // Import JWT helper for authentication
)

var (
	receipts = make(map[string]*models.ProcessedReceipt) // In-memory store for processed receipts
	mu       sync.RWMutex                                 // Mutex for thread-safe access to receipts map
)

// ProcessReceipt handles the POST request to process a receipt.
// It validates the receipt, calculates points, generates a unique ID,
// and stores it in memory.
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Verify JWT token from Authorization header for secure access
	if !utils.ValidateJWT(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var receipt models.Receipt
	// Parse JSON body into Receipt struct
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate receipt data before processing
	if err := validateReceipt(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate points based on receipt rules
	points := calculatePoints(&receipt)

	// Generate a unique ID for the processed receipt
	id := uuid.New().String()
	processedReceipt := &models.ProcessedReceipt{ID: id, Points: points}

	// Store the processed receipt in the in-memory store
	mu.Lock()
	receipts[id] = processedReceipt
	mu.Unlock()

	// Respond with the generated receipt ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// GetPoints handles the GET request to retrieve points for a specific receipt.
// It fetches the receipt by ID and returns the points awarded.
func GetPoints(w http.ResponseWriter, r *http.Request) {
	// Verify JWT token from Authorization header
	if !utils.ValidateJWT(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract the receipt ID from the request URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Safely retrieve receipt points with read-lock
	mu.RLock()
	receipt, exists := receipts[id]
	mu.RUnlock()

	// Handle case where receipt ID does not exist in the store
	if !exists {
		http.Error(w, "No receipt found for that ID", http.StatusNotFound)
		return
	}

	// Send points in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": receipt.Points})
}

// validateReceipt performs validation on the receipt data, ensuring required fields
// are present and correctly formatted.
func validateReceipt(r *models.Receipt) error {
	// Check for missing fields in receipt
	if r.Retailer == "" {
		return fmt.Errorf("retailer is required")
	}
	if r.PurchaseDate == "" {
		return fmt.Errorf("purchaseDate is required")
	}
	if r.PurchaseTime == "" {
		return fmt.Errorf("purchaseTime is required")
	}
	if len(r.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	if r.Total == "" {
		return fmt.Errorf("total is required")
	}

	// Regular expression validations for specific fields
	retailerRegex := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !retailerRegex.MatchString(r.Retailer) {
		return fmt.Errorf("invalid retailer name format")
	}

	// Validate date format (expected YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
		return fmt.Errorf("invalid purchase date format")
	}

	// Validate time format (expected HH:MM in 24-hour format)
	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return fmt.Errorf("invalid purchase time format")
	}

	// Validate total amount format (expected 0.00)
	totalRegex := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !totalRegex.MatchString(r.Total) {
		return fmt.Errorf("invalid total format")
	}

	// Validate each item in the receipt
	for _, item := range r.Items {
		if err := validateItem(&item); err != nil {
			return err
		}
	}

	return nil
}

// validateItem validates individual item data in the receipt, checking for
// required fields and proper formatting.
func validateItem(i *models.Item) error {
	// Check for missing fields in item
	if i.ShortDescription == "" {
		return fmt.Errorf("item short description is required")
	}
	if i.Price == "" {
		return fmt.Errorf("item price is required")
	}

	// Validate description format
	descRegex := regexp.MustCompile(`^[\w\s\-]+$`)
	if !descRegex.MatchString(i.ShortDescription) {
		return fmt.Errorf("invalid item short description format")
	}

	// Validate price format (expected 0.00)
	priceRegex := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !priceRegex.MatchString(i.Price) {
		return fmt.Errorf("invalid item price format")
	}

	return nil
}

// calculatePoints calculates the points for the receipt based on predefined rules.
func calculatePoints(r *models.Receipt) int {
	points := 0

	// Rule 1: One point per alphanumeric character in retailer name
	points += countAlphanumeric(r.Retailer)

	// Rule 2: 50 points if the total is a round dollar amount
	if strings.HasSuffix(r.Total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if isTotalMultipleOf25Cents(r.Total) {
		points += 25
	}

	// Rule 4: 5 points for every two items
	points += (len(r.Items) / 2) * 5

	// Rule 5: Extra points if item description length is multiple of 3
	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if purchase day is odd
	if isPurchaseDateOdd(r.PurchaseDate) {
		points += 6
	}

	// Rule 7: 10 points if purchase time is between 2:00pm and 4:00pm
	if isPurchaseTimeBetween2And4PM(r.PurchaseTime) {
		points += 10
	}

	return points
}

// Helper functions for calculating points

// countAlphanumeric counts alphanumeric characters in a string.
func countAlphanumeric(s string) int {
	count := 0
	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			count++
		}
	}
	return count
}

// isTotalMultipleOf25Cents checks if the total is a multiple of 0.25.
func isTotalMultipleOf25Cents(total string) bool {
	f, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return false
	}
	cents := int(f * 100)
	return cents%25 == 0
}

// isPurchaseDateOdd checks if the purchase date day is odd.
func isPurchaseDateOdd(date string) bool {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	return t.Day()%2 != 0
}

// isPurchaseTimeBetween2And4PM checks if purchase time is between 2:00pm and 4:00pm.
func isPurchaseTimeBetween2And4PM(timeStr string) bool {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	return t.Hour() >= 14 && t.Hour() < 16
}
