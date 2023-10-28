package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"nicked.io/db"
)

type User struct {
	Id              int64            `bun:"id,pk,autoincrement"`
	Email           string           `bun:",notnull"`
	EmailAlerts     bool             `bun:",notnull"`
	Products        []Product        `bun:"m2m:user_to_products,join:User=Product"`
	ProductSettings []ProductSetting `bun:"rel:has-many,join:id=user_id"`
	CreatedAt       time.Time        `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt       time.Time        `bun:",nullzero,notnull,default:current_timestamp"`
}

func GetUser(id *int64, ctx context.Context) (*User, error) {
	if id == nil {
		return nil, errors.New("invalid user id")
	}

	user := new(User)
	err := db.Client.NewSelect().
		Model(user).
		Where("id = ?", id).
		Relation("Products").
		Relation("ProductSettings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(email string, ctx context.Context) (*User, error) {
	if email == "" {
		return nil, errors.New("invalid email")
	}

	user := new(User)
	err := db.Client.NewSelect().
		Model(user).
		Where("email = ?", email).
		Relation("Products").
		Relation("ProductSettings").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUser(user *User, ctx context.Context) (*User, error) {
	if user == nil {
		return nil, errors.New("invalid user data")
	}

	user.Email = strings.ToLower(user.Email)
	user.EmailAlerts = true

	exists, err := db.Client.NewSelect().
		Model((*User)(nil)).
		Where("email = ?", user.Email).
		Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err := db.Client.NewInsert().
			Model(user).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	newUser, err := GetUserByEmail(user.Email, ctx)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func UpdateUser(user *User, ctx context.Context) (*User, error) {
	if user == nil {
		return nil, errors.New("invalid user data")
	}

	user.UpdatedAt = time.Now()
	_, err := db.Client.NewUpdate().
		Model(user).
		Column("email", "email_alerts", "updated_at").
		Where("id = ?", user.Id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(id int64, ctx context.Context) error {
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
