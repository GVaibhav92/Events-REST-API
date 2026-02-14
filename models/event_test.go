package models

import (
	"REST-API/config"
	"REST-API/db"
	"os"
	"testing"
	"time"
)

// setupTestDB creates a fresh in-memory database for each test
// so tests never interfere with each other or your real api.db
func setupTestDB(t *testing.T) {
	t.Helper()

	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("DB_PATH", ":memory:")

	config.Load()
	db.InitDB()
}

func TestSaveEvent(t *testing.T) {
	setupTestDB(t)

	event := Event{
		Name:        "Test Event",
		Description: "A test event description",
		Location:    "Test Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}

	err := event.Save()
	if err != nil {
		t.Errorf("expected no error saving event, got: %v", err)
	}

	if event.ID == 0 {
		t.Error("expected event ID to be set after save, got 0")
	}
}

func TestGetAllEvents(t *testing.T) {
	setupTestDB(t)

	// Save two events first
	for i := 0; i < 2; i++ {
		event := Event{
			Name:        "Test Event",
			Description: "A test event description",
			Location:    "Test Location",
			DateTime:    time.Now().Add(24 * time.Hour),
			UserID:      1,
		}
		event.Save()
	}

	events, total, err := GetAllEvents(1, 10)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventByID(t *testing.T) {
	setupTestDB(t)

	event := Event{
		Name:        "Findable Event",
		Description: "This event should be findable by ID",
		Location:    "Somewhere",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}
	event.Save()

	found, err := GetEventByID(event.ID)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find event, got nil")
	}
	if found.Name != event.Name {
		t.Errorf("expected name %s, got %s", event.Name, found.Name)
	}
}

func TestGetEventByID_NotFound(t *testing.T) {
	setupTestDB(t)

	found, err := GetEventByID(999)
	if err != nil {
		t.Errorf("expected no error for missing event, got: %v", err)
	}
	if found != nil {
		t.Error("expected nil for non-existent event, got a result")
	}
}

func TestUpdateEvent(t *testing.T) {
	setupTestDB(t)

	event := Event{
		Name:        "Original Name",
		Description: "Original description here",
		Location:    "Original Location",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}
	event.Save()

	event.Name = "Updated Name"
	event.Description = "Updated description here"
	err := event.Update()
	if err != nil {
		t.Errorf("expected no error on update, got: %v", err)
	}

	updated, _ := GetEventByID(event.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}
}

func TestDeleteEvent(t *testing.T) {
	setupTestDB(t)

	event := Event{
		Name:        "To Be Deleted",
		Description: "This event will be deleted",
		Location:    "Nowhere",
		DateTime:    time.Now().Add(24 * time.Hour),
		UserID:      1,
	}
	event.Save()

	err := event.Delete()
	if err != nil {
		t.Errorf("expected no error on delete, got: %v", err)
	}

	found, _ := GetEventByID(event.ID)
	if found != nil {
		t.Error("expected event to be deleted, but it still exists")
	}
}

func TestPagination(t *testing.T) {
	setupTestDB(t)

	// Create 5 events
	for i := 0; i < 5; i++ {
		event := Event{
			Name:        "Paginated Event",
			Description: "Testing pagination behavior",
			Location:    "Anywhere",
			DateTime:    time.Now().Add(24 * time.Hour),
			UserID:      1,
		}
		event.Save()
	}

	// Request page 1 with limit 3 — should get 3 results
	events, total, _ := GetAllEvents(1, 3)
	if len(events) != 3 {
		t.Errorf("expected 3 events on page 1, got %d", len(events))
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	// Request page 2 with limit 3 — should get remaining 2
	events, _, _ = GetAllEvents(2, 3)
	if len(events) != 2 {
		t.Errorf("expected 2 events on page 2, got %d", len(events))
	}
}
