package models

import (
	"time"

	"github.com/kataras/iris/v12"

    "pricetracker/db"
)

type Price struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    MerchantId *int64 `bun:",notnull"`
    Value *int64 `bun:",notnull"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitPrice(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*Price)(nil)).
        IfNotExists().
        ForeignKey(`("merchant_id") REFERENCES "merchant" ("id") ON DELETE CASCADE`).
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}
