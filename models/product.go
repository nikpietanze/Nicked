package models

import (
	"context"
	"errors"
	"time"

	"Nicked/db"
)

type Product struct {
	Id        int64     `bun:"id,pk,autoincrement"`
	Active    bool      `bun:",notnull"`
	ImageUrl  string    `bun:",notnull"`
	Name      string    `bun:",notnull"`
	Prices    []*Price  `bun:"rel:has-many,join:id=product_id"`
	Sku       string    `bun:",notnull"`
	Store     string    `bun:",notnull"`
	Url       string    `bun:",notnull"`
	UserId    int64     `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitProduct(ctx context.Context) error {
	_, err := db.Client.NewCreateTable().
		Model((*Product)(nil)).
		IfNotExists().
		ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
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

func GetProductBySku(sku string, store string, userId int64, ctx context.Context) (*Product, error) {
	product := new(Product)
	err := db.Client.NewSelect().
		Model(product).
		Where("sku = ?", sku).
		Where("store = ?", store).
		Where("user_id = ?", userId).
		Relation("Prices").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func GetActiveProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	err := db.Client.NewSelect().
		Model(&products).
		Where("active = ?", true).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductsByUser(id *int64, email string, ctx context.Context) ([]Product, error) {
	if id == nil && email == "" {
		return nil, errors.New("invalid user id and/or email")
	}

	user := new(User)
	if id != nil {
		usr, err := GetUser(id, ctx)
		if err != nil {
			panic(err)
		}
		user = usr
	} else if email != "" {
		usr, err := GetUserByEmail(email, ctx)
		if err != nil {
			panic(err)
		}
		user = usr
	}

	var products []Product
	err := db.Client.NewSelect().
		Model(&products).
		Where("user_id = ?", user.Id).
		Relation("Prices").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetLastProductByUser(id *int64, email string, ctx context.Context) (*Product, error) {
	if id == nil && email == "" {
		return nil, errors.New("invalid user id and/or email")
	}

	user := new(User)
	if id != nil {
		usr, err := GetUser(id, ctx)
		if err != nil {
			panic(err)
		}
		user = usr
	} else if email != "" {
		usr, err := GetUserByEmail(email, ctx)
		if err != nil {
			panic(err)
		}
		user = usr
	}

	var product Product
	err := db.Client.NewSelect().
		Model(&product).
		Where("user_id = ?", user.Id).
		Relation("Prices").
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func CreateProduct(product *Product, ctx context.Context) (*Product, error) {
	if product == nil {
		return nil, errors.New("invalid product data")
	}

	exists, err := db.Client.NewSelect().
		Model((*Product)(nil)).
		Where("sku = ?", product.Sku).
		Where("store = ?", product.Store).
		Where("user_id = ?", product.UserId).
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

		newProduct, err := GetLastProductByUser(&product.UserId, "", ctx)
		if err != nil {
			return nil, err
		}
		return newProduct, nil
	} else {
		newProduct, err := GetProductBySku(product.Sku, product.Store, product.UserId, ctx)
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
		Where("sku = ?", product.Sku).
		Where("store = ?", product.Store).
		Where("user_id = ?", product.UserId).
		Column("active", "updated_at").
		OmitZero().
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
