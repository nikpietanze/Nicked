package models

import (
	"context"
	"errors"
	"log"
	"time"

	"Nicked/db"
)

type Price struct {
	Id        int64     `bun:"id,pk,autoincrement"`
	Amount    float64   `bun:",notnull"`
	Currency  string    `bun:",notnull"`
	Store     string    `bun:",notnull"`
	ItemId    int64     `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitPrice(ctx context.Context) error {
	_, err := db.Client.NewCreateTable().
		Model((*Price)(nil)).
		IfNotExists().
		ForeignKey(`("item_id") REFERENCES "items" ("id") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetPrice(id *int64, ctx context.Context) (*Price, error) {
	if id == nil {
		return nil, errors.New("missing or invalid price id")
	}

	price := new(Price)
	err := db.Client.NewSelect().
		Model(price).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return price, nil
}

func GetPricesByItem(itemId *int64, ctx context.Context) ([]Price, error) {
	if itemId == nil {
		return nil, errors.New("missing or invalid item id")
	}

	var prices []Price
	err := db.Client.NewSelect().
		Model(&prices).
		Where("item_id = ?", itemId).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func GetLatestPriceByItem(itemId *int64, ctx context.Context) (*Price, error) {
	if itemId == nil {
		return nil, errors.New("missing or invalid item id")
	}

	price := new(Price)
	err := db.Client.NewSelect().
		Model(price).
		Where("item_id = ?", itemId).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return price, nil
}

func CreatePrice(price *Price, ctx context.Context) (*Price, error) {
	if price == nil {
		return nil, errors.New("missing or invalid price data")
	}

	if err := InitPrice(ctx); err != nil {
		log.Print(err)
	}

	_, err := db.Client.NewInsert().
		Model(price).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	newPrice, err := GetLatestPriceByItem(&price.ItemId, ctx)
	if err != nil {
		return nil, err
	}
	return newPrice, nil
}

func DeletePrice(id int64, ctx context.Context) error {
	price := Price{
		Id: id,
	}
	_, err := db.Client.NewDelete().
		Model(&price).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
