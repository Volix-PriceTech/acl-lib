# ACL Library for Go

A simple and flexible Access Control List (ACL) library for Go projects.

## Features

- Manage roles and permissions
- Assign permissions to roles
- Database-agnostic storage interface
- Thread-safe operations

## Installation

```bash
go get github.com/yourusername/acl-library
```

## Usage

```go
package main

import (
	"database/sql"

	"github.com/Volix-PriceTech/acl-lib"
	"github.com/Volix-PriceTech/acl-lib/storage"
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

	// Create a new role
	err = acl.CreateRole("admin")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new permission
	err = acl.CreatePermission("read_articles")
	if err != nil {
		log.Fatal(err)
	}

	// Assign permission to role
	err = acl.AssignPermission(roleID, permissionID)
	if err != nil {
		log.Fatal(err)
	}
}
```