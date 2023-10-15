package models

import (
	"context"
	"time"

	"Nicked/db"
)

type Price struct {
	Id        int64     `bun:"id,pk,autoincrement"`
	Amount    float64   `bun:",notnull"`
	Currency  string    `bun:",notnull"`
	ProductId int64     `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitPrice(ctx context.Context) error {
	_, err := db.Client.NewCreateTable().
		Model((*Price)(nil)).
		IfNotExists().
		ForeignKey(`("product_id") REFERENCES "products" ("id") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetPrice(id int64, ctx context.Context) (*Price, error) {
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

func GetPricesByProduct(productId int64, ctx context.Context) ([]Price, error) {
	var prices []Price
	err := db.Client.NewSelect().
		Model(&prices).
		Where("product_id = ?", productId).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func GetLatestPriceByProduct(productId int64, ctx context.Context) (*Price, error) {
	price := new(Price)
	err := db.Client.NewSelect().
		Model(price).
		Where("product_id = ?", productId).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return price, nil
}

func CreatePrice(price Price, ctx context.Context) (*Price, error) {
	_, err := db.Client.NewInsert().
		Model(&price).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	newPrice, err := GetLatestPriceByProduct(price.ProductId, ctx)
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
