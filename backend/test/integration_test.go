package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	authServiceURL    = "http://auth-service:8080"
	bookingServiceURL = "http://booking-service:8080"

	testEmail     = "test@example.com"
	testPassword  = "password1234"
	testFirstName = "Test"
	testLastName  = "User"
	testRole      = "student"

	testCoworkingID = "550e8400-e29b-41d4-a716-446655440000"
	testPlaceID     = "550e8400-e29b-41d4-a716-446655441001"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Response structures
type RegisterResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type CoworkingsResponse struct {
	Coworkings []Coworking `json:"coworkings"`
}

type Coworking struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AvailablePlacesResponse struct {
	Places []Place `json:"places"`
}

type Place struct {
	ID            uuid.UUID `json:"id"`
	CoworkingID   string    `json:"coworkingId"`
	CoworkingName string    `json:"coworkingName"`
	Label         string    `json:"label"`
	PlaceType     string    `json:"placeType"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BookingsResponse struct {
	Bookings   []Booking      `json:"bookings"`
	Pagination PaginationMeta `json:"pagination"`
}

type Booking struct {
	ID           uuid.UUID  `json:"id"`
	UserID       string     `json:"userId"`
	UserName     string     `json:"userName"`
	Place        Place      `json:"place"`
	StartTime    time.Time  `json:"startTime"`
	EndTime      time.Time  `json:"endTime"`
	Status       string     `json:"status"`
	CancelReason *string    `json:"cancelReason,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	CancelledAt  *time.Time `json:"cancelledAt,omitempty"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// Helper functions
func doRequest(method, url string, body interface{}, accessToken string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	return httpClient.Do(req)
}

func parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// Test scenario
func TestIntegrationBookingFlow(t *testing.T) {
	// 1. Register user
	t.Log("Step 1: Registering user...")
	registerReq := map[string]string{
		"email":      testEmail,
		"password":   testPassword,
		"first_name": testFirstName,
		"last_name":  testLastName,
		"roleCode":   testRole,
	}

	resp, err := doRequest("POST", authServiceURL+"/auth/register", registerReq, "")
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Register failed with status %d: %s", resp.StatusCode, string(body))
	}

	var registerResp RegisterResponse
	if err := parseResponse(resp, &registerResp); err != nil {
		t.Fatalf("Failed to parse register response: %v", err)
	}
	accessToken := registerResp.AccessToken
	t.Logf("✓ User registered successfully. Token: %s...", accessToken[:20])

	// 2. Get list of coworkings
	t.Log("Step 2: Getting coworkings list...")
	resp, err = doRequest("GET", bookingServiceURL+"/coworkings", nil, accessToken)
	if err != nil {
		t.Fatalf("Failed to get coworkings: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Get coworkings failed with status %d: %s", resp.StatusCode, string(body))
	}

	var coworkingsResp CoworkingsResponse
	if err := parseResponse(resp, &coworkingsResp); err != nil {
		t.Fatalf("Failed to parse coworkings response: %v", err)
	}

	if len(coworkingsResp.Coworkings) == 0 {
		t.Fatalf("No coworkings found")
	}
	t.Logf("✓ Got %d coworking(s). First: %s", len(coworkingsResp.Coworkings), coworkingsResp.Coworkings[0].Name)

	// 3. Get available places for specific time
	t.Log("Step 3: Getting available places...")
	now := time.Now().UTC()
	startTime := now.Add(2 * time.Hour).Truncate(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	url := fmt.Sprintf("%s/coworkings/%s/available-places?startTime=%s&endTime=%s",
		bookingServiceURL,
		testCoworkingID,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	resp, err = doRequest("GET", url, nil, accessToken)
	if err != nil {
		t.Fatalf("Failed to get available places: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Get available places failed with status %d: %s", resp.StatusCode, string(body))
	}

	var placesResp AvailablePlacesResponse
	if err := parseResponse(resp, &placesResp); err != nil {
		t.Fatalf("Failed to parse available places response: %v", err)
	}

	if len(placesResp.Places) == 0 {
		t.Fatalf("No available places found")
	}
	t.Logf("✓ Got %d available place(s). First: %s", len(placesResp.Places), placesResp.Places[0].Label)

	// 4. Create booking
	t.Log("Step 4: Creating booking...")
	bookingReq := map[string]interface{}{
		"placeId":   testPlaceID,
		"startTime": startTime,
		"endTime":   endTime,
	}

	resp, err = doRequest("POST", bookingServiceURL+"/bookings", bookingReq, accessToken)
	if err != nil {
		t.Fatalf("Failed to create booking: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create booking failed with status %d: %s", resp.StatusCode, string(body))
	}
	t.Log("✓ Booking created successfully")

	// 5. Get active bookings
	t.Log("Step 5: Getting user's active bookings...")
	respURL := fmt.Sprintf("%s/bookings/active?page=1&pageSize=10", bookingServiceURL)
	resp, err = doRequest("GET", respURL, nil, accessToken)
	if err != nil {
		t.Fatalf("Failed to get bookings: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Get bookings failed with status %d: %s", resp.StatusCode, string(body))
	}

	var bookingsResp BookingsResponse
	if err := parseResponse(resp, &bookingsResp); err != nil {
		t.Fatalf("Failed to parse bookings response: %v", err)
	}

	if len(bookingsResp.Bookings) == 0 {
		t.Fatalf("No active bookings found")
	}

	activeBooking := bookingsResp.Bookings[0]
	bookingID := activeBooking.ID
	t.Logf("✓ Got %d active booking(s). ID: %s, Place: %s, Status: %s",
		len(bookingsResp.Bookings), bookingID, activeBooking.Place.Label, activeBooking.Status)

	if activeBooking.Status != "active" {
		t.Fatalf("Expected booking status 'active', got '%s'", activeBooking.Status)
	}

	// 6. Cancel booking
	t.Log("Step 6: Cancelling booking...")
	cancelReq := map[string]string{
		"reason": "Test cancellation",
	}

	cancelURL := fmt.Sprintf("%s/bookings/%s", bookingServiceURL, bookingID)
	resp, err = doRequest("DELETE", cancelURL, cancelReq, accessToken)
	if err != nil {
		t.Fatalf("Failed to cancel booking: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Cancel booking failed with status %d: %s", resp.StatusCode, string(body))
	}
	t.Log("✓ Booking cancelled successfully")

	// 7. Verify booking is cancelled
	t.Log("Step 7: Verifying booking is cancelled...")
	respURL = fmt.Sprintf("%s/bookings/history?page=1&pageSize=10", bookingServiceURL)
	resp, err = doRequest("GET", respURL, nil, accessToken)
	if err != nil {
		t.Fatalf("Failed to get all bookings: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Get all bookings failed with status %d: %s", resp.StatusCode, string(body))
	}

	var allBookingsResp BookingsResponse
	if err := parseResponse(resp, &allBookingsResp); err != nil {
		t.Fatalf("Failed to parse all bookings response: %v", err)
	}

	// Find our cancelled booking
	var cancelledBooking *Booking
	for _, b := range allBookingsResp.Bookings {
		if b.ID == activeBooking.ID {
			cancelledBooking = &b
			break
		}
	}

	if cancelledBooking == nil {
		t.Fatalf("Booking with ID %s not found in all bookings list", activeBooking.ID)
	}

	if cancelledBooking.Status != "cancelled" {
		t.Fatalf("Expected booking status 'cancelled', got '%s'", cancelledBooking.Status)
	}

	if cancelledBooking.CancelledAt == nil {
		t.Fatalf("Expected CancelledAt to be set, but it's nil")
	}

	t.Logf("✓ Booking successfully cancelled. Status: %s, Cancelled at: %s",
		cancelledBooking.Status, cancelledBooking.CancelledAt.Format(time.RFC3339))

	t.Log("\n✅ All tests passed!")
}
