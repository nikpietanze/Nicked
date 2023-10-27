package models

import (
	"context"
	"errors"
	"time"

	"Nicked/db"
)

type ProductSetting struct {
	Id        int64     `bun:"id,pk,autoincrement"`
	Active    bool      `bun:",notnull"`
	ProductId int64     `bun:",notnull"`
	UserId    int64     `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func GetProductSetting(id int64, ctx context.Context) (*ProductSetting, error) {
	productSetting := new(ProductSetting)
	err := db.Client.NewSelect().
		Model(productSetting).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return productSetting, nil
}

func CreateProductSetting(productSetting *ProductSetting, ctx context.Context) (*ProductSetting, error) {
	if productSetting == nil {
		return nil, errors.New("invalid product setting data")
	}

	exists, err := db.Client.NewSelect().
		Model((*ProductSetting)(nil)).
		Where("product_id = ?", productSetting.ProductId).
		Where("user_id = ?", productSetting.UserId).
		Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err := db.Client.NewInsert().
			Model(productSetting).
			Exec(ctx)
		if err != nil {
			return nil, err
		}

		return productSetting, nil
	}
    return productSetting, nil;
}

func UpdateProductSetting(productSetting *ProductSetting, ctx context.Context) (*ProductSetting, error) {
	if productSetting == nil {
		return nil, errors.New("invalid product setting data")
	}

	productSetting.UpdatedAt = time.Now()
	_, err := db.Client.NewUpdate().
		Model(productSetting).
		Column("active", "updated_at").
        Where("id = ?", productSetting.Id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return productSetting, nil
}

func DeleteProductSetting(id int64, ctx context.Context) error {
	productSetting := ProductSetting{
		Id: id,
	}
	_, err := db.Client.NewDelete().
		Model(&productSetting).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
