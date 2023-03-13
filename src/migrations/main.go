package main

import (
	"github.com/alexkalak/pony_express-calculator/src/db"
	"github.com/alexkalak/pony_express-calculator/src/migrations/migration"
)

func main() {
	db.Init()
	migration.Migrate()
}
