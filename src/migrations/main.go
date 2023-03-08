package main

import (
	"github.com/alexkalak/pony_express/src/db"
	"github.com/alexkalak/pony_express/src/migrations/migration"
)

func main() {
	db.Init()
	migration.Migrate()
}
