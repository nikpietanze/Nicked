package models

import (
	"context"
	"errors"
	"log"

	"nicked.io/db"
)

type UserToProduct struct {
	UserId    int64    `bun:",pk"`
	User      *User    `bun:"rel:belongs-to,join:user_id=id"`
	ProductId int64    `bun:",pk"`
	Product   *Product `bun:"rel:belongs-to,join:product_id=id"`
}

func CreateUserToProduct(userToProduct *UserToProduct, ctx context.Context) (*UserToProduct, error) {
	if userToProduct == nil {
		return nil, errors.New("invalid user_to_product data")
	}

    log.Println(userToProduct);

	exists, err := db.Client.NewSelect().
		Model((*UserToProduct)(nil)).
		Where("product_id = ?", userToProduct.ProductId).
		Where("user_id = ?", userToProduct.UserId).
		Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
        log.Println("doesn't exist, inserting...");
		_, err := db.Client.NewInsert().
			Model(userToProduct).
			Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
    return userToProduct, nil;
}

