package models

import (
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
)

type Price struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    MerchantId *int64
    Value *int64 `bun:",notnull"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitPrice(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*Price)(nil)).
        IfNotExists().
        ForeignKey(`("merchant_id") REFERENCES "merchants" ("id") ON DELETE CASCADE`).
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}

func GetPrice(id *int64, ctx iris.Context) *Price {
    price := new(Price)
    err := db.Client.NewSelect().
        Model(&price).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    return price
}

func CreatePrice(price *Price, ctx iris.Context) int64 {
    res, err := db.Client.NewInsert().
        Model(price).
        Returning("id").
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    id, err := res.RowsAffected()
    if err != nil {
        panic(err)
    }

    return id
}

func UpdatePrice(price *Price, ctx iris.Context) int64 {
    price.UpdatedAt = time.Now()

    res, err := db.Client.NewUpdate().
        Model(price).
        Column("name", "updated_at").
        OmitZero().
        WherePK().
        Returning("id").
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    id, err := res.RowsAffected()
    if err != nil {
        panic(err)
    }

    return id
}

func DeletePrice(id *int64, ctx iris.Context) *int64 {
    price := Price {
        Id: id,
    }

    _, err := db.Client.NewDelete().
        Model(&price).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    return id
}
