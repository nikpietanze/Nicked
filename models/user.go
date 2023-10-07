package models

import (
	"errors"
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
)

type User struct {
    Id        int64    `bun:"id,pk,autoincrement"`
	Email     string    `bun:",notnull"`
	Items     []*Item   `bun:"rel:has-many,join:id=user_id"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitUser(ctx iris.Context) {
	_, err := db.Client.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		ctx.StopWithJSON(500, NewError(err))
	}
}

func GetUser(id *int64, ctx iris.Context) (*User, error) {
    if (id == nil) {
        return nil, errors.New("missing or invalid user id")
    }

	user := new(User)
	err := db.Client.NewSelect().
		Model(user).
		Where("id = ?", id).
        Relation("Items").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(email string, ctx iris.Context) (*User, error) {
    if (email == "") {
        return nil, errors.New("missing or invalid email")
    }

	user := new(User)
	err := db.Client.NewSelect().
		Model(user).
		Where("email = ?", email).
        Relation("Items").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(user *User, ctx iris.Context) (*User, error) {
    if (user == nil) {
        return nil, errors.New("missing or invalid user data")
    }

	user.Email = strings.ToLower(user.Email)
	_, err := db.Client.NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

    newUser, err := GetUserByEmail(user.Email, ctx)
    if (err != nil) {
        return nil, err
    }

	return newUser, nil
}

func UpdateUser(user *User, ctx iris.Context) (*User, error) {
    if (user == nil) {
        return nil, errors.New("missing or invalid user data")
    }

	user.UpdatedAt = time.Now()
	_, err := db.Client.NewUpdate().
		Model(user).
		Column("name", "updated_at").
		OmitZero().
		WherePK().
		Returning("id").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(id int64, ctx iris.Context) error {
	user := User{
		Id: id,
	}
	_, err := db.Client.NewDelete().
		Model(&user).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
