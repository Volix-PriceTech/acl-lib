# ACL Library for Go

A simple and flexible Access Control List (ACL) library for Go projects.

## Features

- Manage roles and permissions
- Assign permissions to roles
- Database-agnostic storage interface
- Thread-safe operations

## Installation

```bash
go get github.com/Volix-PriceTech/acl-lib
```

## Usage

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Volix-PriceTech/acl-library"
	"github.com/Volix-PriceTech/acl-library/storage"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create a context for the operations
	ctx := context.Background()

	// Retrieve database credentials from environment variables for security
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Verify the database connection is alive
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// Initialize MySQL storage
	storage := storage.NewMySQLStorage(db)

	// Initialize the ACL system
	aclSystem := acl_lib.NewACL(storage)

	// Migrate tables (create necessary tables if they don't exist)
	if err := aclSystem.Migrate(ctx); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}
	log.Println("Database tables migrated successfully.")

	// Create a new role
	adminRole, err := aclSystem.CreateRole(ctx, "admin")
	if err != nil {
		log.Fatalf("Failed to create role 'admin': %v", err)
	}
	log.Printf("Created role: ID=%d, Name=%s\n", adminRole.ID, adminRole.Name)

	// Create a new permission
	readPermission, err := aclSystem.CreatePermission(ctx, "read_articles")
	if err != nil {
		log.Fatalf("Failed to create permission 'read_articles': %v", err)
	}
	log.Printf("Created permission: ID=%d, Name=%s\n", readPermission.ID, readPermission.Name)

	// Assign permission to role
	if err := aclSystem.AssignPermission(ctx, adminRole.ID, readPermission.ID); err != nil {
		log.Fatalf("Failed to assign permission to role: %v", err)
	}
	log.Printf("Assigned permission ID=%d to role ID=%d\n", readPermission.ID, adminRole.ID)

	// Assign role to user
	userID := "user123"
	if err := aclSystem.AssignRoleToUser(ctx, adminRole.ID, userID); err != nil {
		log.Fatalf("Failed to assign role to user: %v", err)
	}
	log.Printf("Assigned role ID=%d to user ID=%s\n", adminRole.ID, userID)

	// Retrieve role by user
	userRole, err := aclSystem.GetRoleByUser(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to retrieve role for user ID=%s: %v", userID, err)
	}
	log.Printf("User ID=%s has role: ID=%d, Name=%s\n", userID, userRole.ID, userRole.Name)

	// Remove role from user
	if err := aclSystem.RemoveRoleFromUser(ctx, adminRole.ID, userID); err != nil {
		log.Fatalf("Failed to remove role from user: %v", err)
	}
	log.Printf("Removed role ID=%d from user ID=%s\n", adminRole.ID, userID)

	// Verify role removal by attempting to retrieve the role again
	userRole, err = aclSystem.GetRoleByUser(ctx, userID)
	if err != nil {
		log.Printf("As expected, user ID=%s has no role assigned: %v\n", userID, err)
	} else {
		log.Printf("Unexpectedly, user ID=%s still has role: ID=%d, Name=%s\n", userID, userRole.ID, userRole.Name)
	}
}

```

## **Testing the MigrateTables Method**

It's important to add tests for the new `MigrateTables()` method to ensure it works as expected.

**`tests/storage_test.go`**

```go
package tests

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/Volix-PriceTech/acl-library/storage"
	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLStorage_MigrateTables(t *testing.T) {
	// Define your test database credentials
	// It's recommended to use environment variables or a configuration file for sensitive information
	dsn := "user:password@tcp(127.0.0.1:3306)/testdb?parseTime=true"

	// Open a connection to the test database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database connection: %v", err)
		}
	}()

	// Verify the database connection is alive
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize MySQLStorage
	storageMySQL := storage.NewMySQLStorage(db)

	// Clean up tables before testing to ensure a fresh state
	tablesToDrop := []string{"user_roles", "roles_permissions", "roles", "permissions"}
	for _, table := range tablesToDrop {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			t.Fatalf("Failed to drop table %s: %v", table, err)
		}
	}

	// Run migration to create tables
	if err := storageMySQL.MigrateTables(); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	// Define the expected tables after migration
	expectedTables := []string{"users", "roles", "permissions", "user_roles", "roles_permissions"}

	// Verify that each expected table exists
	for _, table := range expectedTables {
		var exists string
		query := fmt.Sprintf("SHOW TABLES LIKE '%s'", table)
		err := db.QueryRow(query).Scan(&exists)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				t.Errorf("Expected table '%s' to exist, but it does not", table)
			} else {
				t.Errorf("Error checking existence of table '%s': %v", table, err)
			}
			continue
		}
		if exists != table {
			t.Errorf("Expected table name '%s', got '%s'", table, exists)
		}
	}
}
```