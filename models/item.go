package models

import (
	"errors"
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
)

type Item struct {
	Id        int64    `bun:"id,pk,autoincrement"`
	IsActive  bool      `bun:",notnull"`
	Name      string    `bun:",notnull"`
	Prices    []*Price  `bun:"rel:has-many,join:id=item_id"`
	Url       string    `bun:",notnull"`
	UserId    int64    `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitItem(ctx iris.Context) error {
	_, err := db.Client.NewCreateTable().
		Model((*Item)(nil)).
		IfNotExists().
		ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetItem(id *int64, ctx iris.Context) (*Item, error) {
    if (id == nil) {
        return nil, errors.New("missing or invalid user id")
    }

	item := new(Item)
	err := db.Client.NewSelect().
		Model(item).
		Where("id = ?", id).
        Relation("Prices").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
    return item, nil
}

func GetActiveItems(ctx iris.Context) ([]Item, error) {
	var items []Item
	err := db.Client.NewSelect().
		Model(items).
		Where("active LIKE ?", "true").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func GetItemsByUser(id *int64, email string, ctx iris.Context) ([]Item, error) {
	if id == nil && email == "" {
		return nil, errors.New("missing or invalid item id and email")
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

	var items []Item
	err := db.Client.NewSelect().
		Model(items).
		Where("user_id = ?", user.Id).
        Relation("Prices").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func GetLastItemByUser(id *int64, email string, ctx iris.Context) (*Item, error) {
	if id == nil && email == "" {
		return nil, errors.New("missing or invalid item id and email")
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

	var item Item
	err := db.Client.NewSelect().
		Model(item).
		Where("user_id = ?", user.Id).
        Relation("Prices").
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func CreateItem(item *Item, ctx iris.Context) (*Item, error) {
	if item == nil {
		return nil, errors.New("missing or invalid item data")
	}

	item.Name = strings.Title(item.Name)
	_, err := db.Client.NewInsert().
		Model(item).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

    newItem, err := GetLastItemByUser(&item.UserId, "", ctx)
    if (err != nil) {
        return nil, err
    }

    // queue the item to be crawled to add the current price

	return newItem, nil
}

func UpdateItem(item *Item, ctx iris.Context) (*Item, error) {
	if item == nil {
		return nil, errors.New("missing or invalid item data")
	}

	item.UpdatedAt = time.Now()
	_, err := db.Client.NewUpdate().
		Model(item).
		Column("name", "updated_at").
		OmitZero().
		WherePK().
		Returning("id").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func DeleteItem(id int64, ctx iris.Context) error {
	item := Item{
		Id: id,
	}
	_, err := db.Client.NewDelete().
		Model(&item).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
