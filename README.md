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
	"database/sql"
	"log"

	"github.com/yourusername/acl-library"
	"github.com/yourusername/acl-library/storage"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage := storage.NewMySQLStorage(db)
	acl := acl.NewACL(storage)

	// Migrate tables
	if err := acl.Migrate(); err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	// Create a new role
	if err := acl.CreateRole("admin"); err != nil {
		log.Fatal(err)
	}

	// Create a new permission
	if err := acl.CreatePermission("read_articles"); err != nil {
		log.Fatal(err)
	}

	// Assign permission to role
	if err := acl.AssignPermission(roleID, permissionID); err != nil {
		log.Fatal(err)
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
    "testing"

    "github.com/yourusername/acl-library/storage"
    _ "github.com/go-sql-driver/mysql"
)

func TestMySQLStorage_MigrateTables(t *testing.T) {
    db, err := sql.Open("mysql", "user:password@/testdb")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    storage := storage.NewMySQLStorage(db)

    // Clean up tables before testing
    db.Exec("DROP TABLE IF EXISTS roles_permissions")
    db.Exec("DROP TABLE IF EXISTS roles")
    db.Exec("DROP TABLE IF EXISTS permissions")

    // Test migration
    if err := storage.MigrateTables(); err != nil {
        t.Fatalf("Failed to migrate tables: %v", err)
    }

    // Check if tables exist
    tables := []string{"roles", "permissions", "roles_permissions"}
    for _, table := range tables {
        var exists string
        query := fmt.Sprintf("SHOW TABLES LIKE '%s'", table)
        if err := db.QueryRow(query).Scan(&exists); err != nil {
            t.Errorf("Table %s does not exist after migration", table)
        }
    }
}
```