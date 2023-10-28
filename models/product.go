package models

import (
	"context"
	"errors"
	"time"

	"nicked.io/db"
)

type Product struct {
	Id        int64     `bun:"id,pk,autoincrement"`
	ImageUrl  string    `bun:",notnull"`
	Name      string    `bun:",notnull"`
	OnSale    bool      `bun:",notnull"`
	Prices    []Price   `bun:"rel:has-many,join:id=product_id"`
	Sku       string    `bun:",notnull"`
	Store     string    `bun:",notnull"`
	Url       string    `bun:",notnull"`
	Users     []User    `bun:"m2m:user_to_products,join:Product=User"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func GetProduct(id int64, ctx context.Context) (*Product, error) {
	product := new(Product)
	err := db.Client.NewSelect().
		Model(product).
		Where("id = ?", id).
		Relation("Prices").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func GetProductBySku(sku string, store string, ctx context.Context) (*Product, error) {
	product := new(Product)
	err := db.Client.NewSelect().
		Model(product).
		Where("sku = ?", sku).
		Where("store = ?", store).
		Relation("Prices").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func GetAllProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	err := db.Client.NewSelect().
		Model(&products).
        Relation("Prices").
		Relation("Users").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func CreateProduct(product *Product, ctx context.Context) (*Product, error) {
	if product == nil {
		return nil, errors.New("invalid product data")
	}

	exists, err := db.Client.NewSelect().
		Model((*Product)(nil)).
		Where("sku = ?", product.Sku).
		Where("store = ?", product.Store).
		Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		_, err := db.Client.NewInsert().
			Model(product).
			Exec(ctx)
		if err != nil {
			return nil, err
		}

		newProduct, err := GetProductBySku(product.Sku, product.Store, ctx)
		if err != nil {
			return nil, err
		}

		return newProduct, nil
	} else {
		newProduct, err := GetProductBySku(product.Sku, product.Store, ctx)
		if err != nil {
			return nil, err
		}
		return newProduct, nil
	}
}

func UpdateProduct(product *Product, ctx context.Context) (*Product, error) {
	if product == nil {
		return nil, errors.New("invalid product data")
	}

	product.UpdatedAt = time.Now()
	_, err := db.Client.NewUpdate().
		Model(product).
		Column("active", "updated_at").
		Where("id = ?", product.Id).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func DeleteProduct(id int64, ctx context.Context) error {
	product := Product{
		Id: id,
	}
	_, err := db.Client.NewDelete().
		Model(&product).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
