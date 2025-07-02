package model_test

import (
	"os"
	"testing"
	"time"

	"github.com/socme-projects/backend/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	TMP_DB_PATH = "test.db"
	DB          *gorm.DB
)

// TestMain runs before all tests in the package to set up and tear down the database.
func TestMain(m *testing.M) {
	// Setup
	var err error
	DB, err = gorm.Open(sqlite.Open(TMP_DB_PATH), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database for tests: " + err.Error())
	}

	err = DB.AutoMigrate(&model.User{}, &model.Session{}, &model.Client{}, &model.Alert{})
	if err != nil {
		panic("Failed to auto migrate database for tests: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Teardown
	err = os.Remove(TMP_DB_PATH)
	if err != nil {
		println("Error deleting temporary database file:", err.Error())
	}

	os.Exit(code)
}

func User(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		initialCount := model.GetNumberOfUsers(DB)
		user, err := model.CreateUser(DB, "Test User", "github_id_123")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		if user == nil {
			t.Fatal("Created user is nil")
		}
		if user.ID == "" {
			t.Error("User ID should not be empty")
		}
		if user.Name != "Test User" {
			t.Errorf("Expected user name 'Test User', got '%s'", user.Name)
		}
		if user.GithubID != "github_id_123" {
			t.Errorf("Expected GithubID 'github_id_123', got '%s'", user.GithubID)
		}
		// First user should be admin
		if initialCount == 0 && user.Role != "admin" {
			t.Errorf("Expected first user to be 'admin', got '%s'", user.Role)
		}
		if initialCount > 0 && user.Role != "guest" {
			t.Errorf("Expected subsequent user to be 'guest', got '%s'", user.Role)
		}

		newCount := model.GetNumberOfUsers(DB)
		if newCount != initialCount+1 {
			t.Errorf("User count did not increment correctly. Expected %d, got %d", initialCount+1, newCount)
		}
	})

	t.Run("GetUserByID", func(t *testing.T) {
		user, err := model.CreateUser(DB, "Another User", "github_id_456")
		if err != nil {
			t.Fatalf("Failed to create user for lookup: %v", err)
		}

		foundUser, err := model.GetUserByID(DB, user.ID)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}
		if foundUser == nil {
			t.Fatal("GetUserByID returned nil for existing user")
		}
		if foundUser.ID != user.ID {
			t.Errorf("Expected user ID '%s', got '%s'", user.ID, foundUser.ID)
		}
		if foundUser.Name != user.Name {
			t.Errorf("Expected user name '%s', got '%s'", user.Name, foundUser.Name)
		}

		// Test for non-existent user
		nonExistentUser, err := model.GetUserByID(DB, "non-existent-id")
		if err == nil {
			t.Error("Expected an error for non-existent user, got none")
		}
		if nonExistentUser != nil {
			t.Error("GetUserByID returned a user for non-existent ID")
		}
	})

	t.Run("GetUserByGithubID", func(t *testing.T) {
		user, err := model.CreateUser(DB, "GitHub User", "unique_github_id_789")
		if err != nil {
			t.Fatalf("Failed to create user for GitHub ID lookup: %v", err)
		}

		foundUser, err := model.GetUserByGithubID(DB, user.GithubID)
		if err != nil {
			t.Fatalf("Failed to get user by GithubID: %v", err)
		}
		if foundUser == nil {
			t.Fatal("GetUserByGithubID returned nil for existing user")
		}
		if foundUser.GithubID != user.GithubID {
			t.Errorf("Expected GithubID '%s', got '%s'", user.GithubID, foundUser.GithubID)
		}

		// Test for non-existent GithubID
		nonExistentUser, err := model.GetUserByGithubID(DB, "non-existent-github-id")
		if err == nil {
			t.Error("Expected an error for non-existent GithubID, got none")
		}
		if nonExistentUser != nil {
			t.Error("GetUserByGithubID returned a user for non-existent GithubID")
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		user, err := model.CreateUser(DB, "User to Delete", "delete_id_111")
		if err != nil {
			t.Fatalf("Failed to create user for deletion: %v", err)
		}
		initialCount := model.GetNumberOfUsers(DB)

		err = model.DeleteUser(DB, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		deletedUser, err := model.GetUserByID(DB, user.ID)
		if err == nil {
			t.Error("Expected an error for a deleted user, got none")
		}
		if deletedUser != nil {
			t.Error("Deleted user was still found")
		}

		newCount := model.GetNumberOfUsers(DB)
		if newCount != initialCount-1 {
			t.Errorf("User count did not decrement correctly after deletion. Expected %d, got %d", initialCount-1, newCount)
		}

		// Test deleting non-existent user
		err = model.DeleteUser(DB, "non-existent-delete-id")
		if err == nil {
			t.Error("Expected an error when trying to delete a non-existent user, got none")
		}
	})
}

func Client(t *testing.T) {
	t.Run("CreateClient", func(t *testing.T) {
		client, err := model.CreateClient(DB,
			"Test Client", "logo.png", "192.168.1.1", "8080",
			"wazuhuser", "wazuhpass", "10.0.0.1", "9200", "indexeruser", "indexerpass")
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		if client == nil {
			t.Fatal("Created client is nil")
		}
		if client.ID == "" {
			t.Error("Client ID should not be empty")
		}
		if client.Name != "Test Client" {
			t.Errorf("Expected client name 'Test Client', got '%s'", client.Name)
		}
		if client.WazuhIP != "192.168.1.1" {
			t.Errorf("Expected WazuhIP '192.168.1.1', got '%s'", client.WazuhIP)
		}
	})

	t.Run("GetClientByID", func(t *testing.T) {
		client, err := model.CreateClient(DB,
			"Another Client", "logo2.png", "192.168.1.2", "8081",
			"wazuhuser2", "wazuhpass2", "10.0.0.2", "9201", "indexeruser2", "indexerpass2")
		if err != nil {
			t.Fatalf("Failed to create client for lookup: %v", err)
		}

		foundClient, err := model.GetClientByID(DB, client.ID)
		if err != nil {
			t.Errorf("Failed to get client by ID: %v", err)
		}
		if foundClient.ID != client.ID {
			t.Errorf("Expected client ID '%s', got '%s'", client.ID, foundClient.ID)
		}
		if foundClient.Name != client.Name {
			t.Errorf("Expected client name '%s', got '%s'", client.Name, foundClient.Name)
		}

		// Test for non-existent client
		_, err = model.GetClientByID(DB, "non-existent-client-id")
		if err == nil {
			t.Error("Expected an error for non-existent client, got none")
		}
	})

	t.Run("GetClientByName", func(t *testing.T) {
		client, err := model.CreateClient(DB,
			"Named Client", "logo3.png", "192.168.1.3", "8082",
			"wazuhuser3", "wazuhpass3", "10.0.0.3", "9202", "indexeruser3", "indexerpass3")
		if err != nil {
			t.Fatalf("Failed to create client for name lookup: %v", err)
		}

		foundClient, err := model.GetClientByName(DB, client.Name)
		if err != nil {
			t.Errorf("Failed to get client by name: %v", err)
		}
		if foundClient.Name != client.Name {
			t.Errorf("Expected client name '%s', got '%s'", client.Name, foundClient.Name)
		}

		// Test for non-existent name
		_, err = model.GetClientByName(DB, "non-existent-client-name")
		if err == nil {
			t.Error("Expected an error for non-existent client name, got none")
		}
	})

	t.Run("DeleteClient", func(t *testing.T) {
		client, err := model.CreateClient(DB,
			"Client to Delete", "logo4.png", "192.168.1.4", "8083",
			"wazuhuser4", "wazuhpass4", "10.0.0.4", "9203", "indexeruser4", "indexerpass4")
		if err != nil {
			t.Fatalf("Failed to create client for deletion: %v", err)
		}

		err = model.DeleteClient(DB, client.ID)
		if err != nil {
			t.Fatalf("Failed to delete client: %v", err)
		}

		_, err = model.GetClientByID(DB, client.ID)
		if err == nil {
			t.Error("Expected an error for a deleted client, got none")
		}

		// Test deleting non-existent client
		err = model.DeleteClient(DB, "non-existent-delete-client-id")
		if err == nil {
			t.Error("Expected an error when trying to delete a non-existent client, got none")
		}
	})
}

func Session(t *testing.T) {
	// First, create a user for the session
	user, err := model.CreateUser(DB, "Session User", "session_user_github_id")
	if err != nil {
		t.Fatalf("Failed to create user for session test: %v", err)
	}

	t.Run("CreateSession", func(t *testing.T) {
		session, err := model.CreateSession(DB, *user) // Pass the user value
		if err != nil {
			t.Fatal("Failed to create session:", err)
		}
		if session.Token == "" {
			t.Error("Session token should not be empty")
		}
		if session.User.ID != user.ID {
			t.Errorf("Session user ID mismatch. Expected %s, got %s", user.ID, session.User.ID)
		}
		if session.Exp.IsZero() {
			t.Error("Session expiration time should not be zero")
		}
		// Allow a small time difference for comparison
		if session.Exp.Before(time.Now().Add(model.SESSION_EXPIRATION-1*time.Second)) ||
			session.Exp.After(time.Now().Add(model.SESSION_EXPIRATION+1*time.Second)) {
			t.Errorf("Session expiration time is not approximately %v from now. Got %v", model.SESSION_EXPIRATION, session.Exp.Sub(time.Now()))
		}

		// Test creating a new session for the same user (should update token)
		updatedSession, err := model.CreateSession(DB, *user)
		if err != nil {
			t.Fatal("Failed to update session for existing user")
		}
		if updatedSession.Token == session.Token {
			t.Error("Session token should have been updated for existing user")
		}
		if updatedSession.User.ID != user.ID {
			t.Errorf("Updated session user ID mismatch. Expected %s, got %s", user.ID, updatedSession.User.ID)
		}
	})

	t.Run("GetSession", func(t *testing.T) {
		session, err := model.CreateSession(DB, *user)
		if err != nil {
			t.Fatalf("Failed to create session for lookup: %v", err)
		}

		foundSession, err := model.GetSession(DB, session.Token)
		if err != nil {
			t.Fatal("GetSession returned nil for existing session")
		}
		if foundSession.Token != session.Token {
			t.Errorf("Expected session token '%s', got '%s'", session.Token, foundSession.Token)
		}
		if foundSession.User.ID != user.ID {
			t.Errorf("Expected session user ID '%s', got '%s'", user.ID, foundSession.User.ID)
		}

		// Test for non-existent session
		_, err = model.GetSession(DB, "non-existent-token")
		if err == nil {
			t.Error("GetSession returned a session for non-existent token")
		}
	})

	t.Run("DeleteSession", func(t *testing.T) {
		session, err := model.CreateSession(DB, *user)
		if err != nil {
			t.Fatalf("Failed to create session for deletion: %v", err)
		}

		err = model.DeleteSession(DB, session.Token)
		if err != nil {
			t.Fatalf("Failed to delete session: %v", err)
		}

		_, err = model.GetSession(DB, session.Token)
		if err == nil {
			t.Error("Deleted session was still found")
		}

		// Test deleting non-existent session
		err = model.DeleteSession(DB, "non-existent-delete-token")
		if err == nil {
			t.Error("Expected an error when trying to delete a non-existent session, got none")
		}
	})

	t.Run("IsExpired", func(t *testing.T) {
		// Create a session that expires in the future
		futureSession := model.Session{
			Token: "future_token",
			Exp:   time.Now().Add(5 * time.Minute),
		}
		if futureSession.IsExpired() {
			t.Error("Session set to expire in the future is reported as expired")
		}

		// Create a session that is already expired
		expiredSession := model.Session{
			Token: "expired_token",
			Exp:   time.Now().Add(-5 * time.Minute),
		}
		if !expiredSession.IsExpired() {
			t.Error("Session set to expire in the past is not reported as expired")
		}
	})
}

func TestModels(t *testing.T) {
	t.Run("UserCRUD", User)
	t.Run("ClientCRUD", Client)
	t.Run("SessionCRUD", Session)
}
