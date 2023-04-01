package main

import (
	"github.com/alexkalak/migration/src/db"
	"github.com/alexkalak/migration/src/migrations/migration"
)

func main() {
	db.Init()
	migration.Migrate()
}
