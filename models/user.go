package models

import (
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
)

type User struct {
    Id *int64 `bun:"id,pk,autoincrement"`
    Email string `bun:",notnull"`
    Items []*Item `bun:"rel:has-many,join:id=user_id"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitUser(ctx iris.Context) {
    _, err := db.Client.NewCreateTable().
        Model((*User)(nil)).
        IfNotExists().
        Exec(ctx)
    if err != nil {
        panic(err)
    }
}

func GetUser(id *int64, ctx iris.Context) *User {
    user := new(User)
    err := db.Client.NewSelect().
        Model(&user).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }
    return user
}

func CreateUser(user *User, ctx iris.Context) int64 {
    user.Email = strings.ToLower(user.Email)

    res, err := db.Client.NewInsert().
        Model(user).
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

func UpdateUser(user *User, ctx iris.Context) int64 {
    user.UpdatedAt = time.Now()

    res, err := db.Client.NewUpdate().
        Model(user).
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

func DeleteUser(id *int64, ctx iris.Context) *int64 {
    user := User{
        Id: id,
    }

    _, err := db.Client.NewDelete().
        Model(&user).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    return id
}

