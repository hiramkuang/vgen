package main

import (
	"testing"
)

func TestUserValidation(t *testing.T) {
	// Test Case 1: Valid User
	t.Run("ValidUser", func(t *testing.T) {
		validUser := &User{
			Name:   "Alice",
			Email:  "alice@example.com",
			Age:    30,
			City:   "Tokyo",  // 5个字符, 符合 len=5
			Status: "active", // 在 in=active,pending,disabled 列表中
		}
		if err := validUser.Validate(); err != nil {
			t.Errorf("Unexpected validation error for validUser: %v", err)
		}
	})

	// Test Case 2: Invalid User (Name too short)
	t.Run("InvalidUser_NameTooShort", func(t *testing.T) {
		invalidUser1 := &User{
			Name:   "A", // 违反 min=2
			Email:  "alice@example.com",
			Age:    30,
			City:   "Tokyo",
			Status: "active",
		}
		if err := invalidUser1.Validate(); err == nil {
			t.Error("Expected validation error for invalidUser1, but got none")
		}
	})

	// Test Case 3: Invalid User (Invalid Email)
	t.Run("InvalidUser_InvalidEmail", func(t *testing.T) {
		invalidUser2 := &User{
			Name:   "Alice",
			Email:  "invalid-email", // 违反 email
			Age:    30,
			City:   "Tokyo",
			Status: "active",
		}
		if err := invalidUser2.Validate(); err == nil {
			t.Error("Expected validation error for invalidUser2, but got none")
		}
	})

	// Test Case 4: Invalid User (Age out of range)
	t.Run("InvalidUser_AgeOutOfRange", func(t *testing.T) {
		invalidUser3 := &User{
			Name:   "Alice",
			Email:  "alice@example.com",
			Age:    200, // 违反 max=150
			City:   "Tokyo",
			Status: "active",
		}
		if err := invalidUser3.Validate(); err == nil {
			t.Error("Expected validation error for invalidUser3, but got none")
		}
	})

	// Test Case 5: Invalid User (City length wrong - 'len' rule)
	t.Run("InvalidUser_CityLengthWrong", func(t *testing.T) {
		userInvalidLen := &User{
			Name:   "Bob",
			Email:  "bob@example.com",
			Age:    25,
			City:   "NYC", // 长度为 3，违反 len=5
			Status: "pending",
		}
		if err := userInvalidLen.Validate(); err == nil {
			t.Error("Expected validation error for userInvalidLen (len rule), but got none")
		}
	})

	// Test Case 6: Invalid User (Status not in list - 'in' rule)
	t.Run("InvalidUser_StatusNotInList", func(t *testing.T) {
		userInvalidIn := &User{
			Name:   "Charlie",
			Email:  "charlie@example.com",
			Age:    35,
			City:   "Paris",    // 长度为 5，符合 len=5
			Status: "archived", // 不在 in=active,pending,disabled 列表中
		}
		if err := userInvalidIn.Validate(); err == nil {
			t.Error("Expected validation error for userInvalidIn (in rule), but got none")
		}
	})

	// Test Case 7: Valid User (All New Rules Pass)
	t.Run("ValidUser_AllNewRulesPass", func(t *testing.T) {
		userAllValid := &User{
			Name:   "Diana",
			Email:  "diana@example.com",
			Age:    28,
			City:   "Milan",   // 长度为 5，符合 len=5
			Status: "pending", // 在 in=active,pending,disabled 列表中
		}
		if err := userAllValid.Validate(); err != nil {
			t.Errorf("Unexpected validation error for userAllValid: %v", err)
		}
	})
}