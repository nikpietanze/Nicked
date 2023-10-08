package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

    "Nicked/config"
)

var Client *bun.DB

func Init() {
    sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DSN)))
    Client = bun.NewDB(sqldb, pgdialect.New())
}
