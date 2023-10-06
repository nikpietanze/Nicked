package models

import (
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
)

type Item struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    IsActive bool `bun:",notnull"`
    Name string `bun:",notnull"`
    Prices []*Price `bun:"rel:has-many,join:id=item_id"`
    Url string `bun:",notnull"`
    UserId *int64 `bun:",notnull"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitItem(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*Item)(nil)).
        IfNotExists().
        ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}

func GetItem(id *int64, ctx iris.Context) *Item {
    item := new(Item)
    err := db.Client.NewSelect().
        Model(&item).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }
    return item
}

func GetActiveItems(ctx iris.Context) []Item {
    var items []Item
    err := db.Client.NewSelect().
        Model(items).
        Where("active LIKE ?", "true").
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    return items
}

func CreateItem(item *Item, ctx iris.Context) int64 {
    item.Name = strings.Title(item.Name)

    res, err := db.Client.NewInsert().
        Model(item).
        Returning("id").
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        panic(err)
    }

    return id
}

func UpdateItem(item *Item, ctx iris.Context) int64 {
    item.UpdatedAt = time.Now()

    res, err := db.Client.NewUpdate().
        Model(item).
        Column("name", "updated_at").
        OmitZero().
        WherePK().
        Returning("id").
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        panic(err)
    }

    return id
}

func DeleteItem(id *int64, ctx iris.Context) *int64 {
    item := Item{
        Id: id,
    }

    _, err := db.Client.NewDelete().
        Model(&item).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    return id
}

