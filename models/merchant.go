package models

import (
	"time"

	"github.com/kataras/iris/v12"

    "pricetracker/db"
)

type Merchant struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    ItemId *int64
    Name string `bun:",notnull"`
    Url string `bun:",notnull"`
    Prices []*Price `bun:"rel:has-many,join:id=merchant_id"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitMerchant(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*Merchant)(nil)).
        IfNotExists().
        ForeignKey(`("item_id") REFERENCES "item" ("id") ON DELETE CASCADE`).
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}

func GetMerchant(id *int64, ctx iris.Context) *Merchant {
    merchant := new(Merchant)
    err := db.Client.NewSelect().
        Model(&merchant).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    return merchant
}
