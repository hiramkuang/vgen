// examples/main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	// --- 保留的旧测试用例 ---
	fmt.Println("--- Test Case 1: Valid User (Old Tests) ---")
	validUser := &User{
		Name:   "Alice",
		Email:  "alice@example.com",
		Age:    30,
		City:   "Tokyo",  // 5个字符, 符合 len=5
		Status: "active", // 在 in=active,pending,disabled 列表中
	}
	if err := validUser.Validate(); err != nil {
		log.Printf("Unexpected validation error for validUser: %v", err)
	} else {
		fmt.Println("validUser is valid!")
	}

	fmt.Println("\n--- Test Case 2: Invalid User (Name too short) (Old Tests) ---")
	invalidUser1 := &User{
		Name:   "A", // 违反 min=2
		Email:  "alice@example.com",
		Age:    30,
		City:   "Tokyo",
		Status: "active",
	}
	if err := invalidUser1.Validate(); err != nil {
		fmt.Printf("Validation failed for invalidUser1 as expected: %v\n", err)
	} else {
		log.Println("Expected validation error for invalidUser1, but got none")
	}

	fmt.Println("\n--- Test Case 3: Invalid User (Invalid Email) (Old Tests) ---")
	invalidUser2 := &User{
		Name:   "Alice",
		Email:  "invalid-email", // 违反 email
		Age:    30,
		City:   "Tokyo",
		Status: "active",
	}
	if err := invalidUser2.Validate(); err != nil {
		fmt.Printf("Validation failed for invalidUser2 as expected: %v\n", err)
	} else {
		log.Println("Expected validation error for invalidUser2, but got none")
	}

	fmt.Println("\n--- Test Case 4: Invalid User (Age out of range) (Old Tests) ---")
	invalidUser3 := &User{
		Name:   "Alice",
		Email:  "alice@example.com",
		Age:    200, // 违反 max=150
		City:   "Tokyo",
		Status: "active",
	}
	if err := invalidUser3.Validate(); err != nil {
		fmt.Printf("Validation failed for invalidUser3 as expected: %v\n", err)
	} else {
		log.Println("Expected validation error for invalidUser3, but got none")
	}
	// --- 旧测试用例结束 ---

	// --- 新增测试用例：测试 'len' 规则 ---
	fmt.Println("\n--- Test Case 5: Invalid User (City length wrong - 'len' rule) ---")
	userInvalidLen := &User{
		Name:   "Bob",
		Email:  "bob@example.com",
		Age:    25,
		City:   "NYC", // 长度为 3，违反 len=5
		Status: "pending",
	}
	if err := userInvalidLen.Validate(); err != nil {
		fmt.Printf("Validation failed for userInvalidLen as expected: %v\n", err)
	} else {
		log.Println("Expected validation error for userInvalidLen (len rule), but got none")
	}

	// --- 新增测试用例：测试 'in' 规则 ---
	fmt.Println("\n--- Test Case 6: Invalid User (Status not in list - 'in' rule) ---")
	userInvalidIn := &User{
		Name:   "Charlie",
		Email:  "charlie@example.com",
		Age:    35,
		City:   "Paris",    // 长度为 5，符合 len=5
		Status: "archived", // 不在 in=active,pending,disabled 列表中
	}
	if err := userInvalidIn.Validate(); err != nil {
		fmt.Printf("Validation failed for userInvalidIn as expected: %v\n", err)
	} else {
		log.Println("Expected validation error for userInvalidIn (in rule), but got none")
	}

	// --- 新增测试用例：验证所有规则都通过 ---
	fmt.Println("\n--- Test Case 7: Valid User (All New Rules Pass) ---")
	userAllValid := &User{
		Name:   "Diana",
		Email:  "diana@example.com",
		Age:    28,
		City:   "Milan",   // 长度为 5，符合 len=5
		Status: "pending", // 在 in=active,pending,disabled 列表中
	}
	if err := userAllValid.Validate(); err != nil {
		log.Printf("Unexpected validation error for userAllValid: %v", err)
	} else {
		fmt.Println("userAllValid is valid!")
	}
}
