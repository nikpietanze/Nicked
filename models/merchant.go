package models

import (
	"strings"
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
        ForeignKey(`("item_id") REFERENCES "items" ("id") ON DELETE CASCADE`).
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

func CreateMerchant(merchant *Merchant, ctx iris.Context) int64 {
    merchant.Name = strings.Title(merchant.Name)
    merchant.Url = strings.ToLower(merchant.Url)

    res, err := db.Client.NewInsert().
        Model(merchant).
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

func UpdateMerchant(merchant *Merchant, ctx iris.Context) int64 {
    merchant.UpdatedAt = time.Now()

    res, err := db.Client.NewUpdate().
        Model(&merchant).
        Column("name", "item_id", "updated_at").
        OmitZero().
        WherePK().
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

func DeleteMerchant(id *int64, ctx iris.Context) *int64 {
    merchant := Merchant {
        Id: id,
    }

    _, err := db.Client.NewDelete().
        Model(&merchant).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    return id
}

func DetermineMerchant(url string) string {
    if strings.Contains(url, "amazon") {
        return "Amazon"
    }
    if strings.Contains(url, "wayfair") {
        return "Wayfair"
    }
    return ""
}
