package models

import (
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
)

type Tracker struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    Email string `bun:",notnull"`
    ItemId *int64
    Active bool `bun:",notnull"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitTracker(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*Tracker)(nil)).
        IfNotExists().
        ForeignKey(`("item_id") REFERENCES "items" ("id") ON DELETE CASCADE`).
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}

func GetTracker(id *int64, ctx iris.Context) *Tracker {
    tracker := new(Tracker)
    err := db.Client.NewSelect().
        Model(tracker).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    return tracker
}

func GetActiveTrackers(ctx iris.Context) []Tracker {
    var trackers []Tracker
    err := db.Client.NewSelect().
        Model(trackers).
        Where("active LIKE ?", "true").
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    return trackers
}

func CreateTracker(tracker *Tracker, ctx iris.Context) int64 {
    tracker.Email = strings.ToLower(tracker.Email)
    tracker.Active = true

    res, err := db.Client.NewInsert().
        Model(tracker).
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

func UpdateTracker(tracker *Tracker, ctx iris.Context) int64 {
    res, err := db.Client.NewUpdate().
        Model(tracker).
        Column("name", "url", "email", "active").
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

func DeleteTracker(id *int64, ctx iris.Context) *int64 {
    tracker := Tracker {
        Id: id,
    }

    _, err := db.Client.NewDelete().
        Model(&tracker).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    return id
}
