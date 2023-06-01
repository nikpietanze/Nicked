package models

import (
	"time"

	"github.com/kataras/iris/v12"

    "pricetracker/db"
)

type Tracker struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    Email string `bun:",notnull"`
    ItemId *int64 `bun:",notnull"`
    Active bool `bun:",notnull"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitTracker(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*Tracker)(nil)).
        IfNotExists().
        ForeignKey(`("item_id") REFERENCES "item" ("id") ON DELETE CASCADE`).
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}

