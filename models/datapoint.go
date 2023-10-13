package models

import (
	"Nicked/db"
	"errors"
	"log"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/uptrace/bun"
)

type DataPoint struct {
	bun.BaseModel `bun:"table:analytics"`
	Id            int64     `bun:"id,pk,autoincrement"`
	Event         string    `bun:",notnull"`
	Location      string    `bun:",notnull"`
	Page          string    `bun:",notnull"`
	Details       string    `bun:",notnull"`
	Data1         string    `bun:",notnull"`
	Data2         string    `bun:",notnull"`
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func InitAnalytics(ctx iris.Context) error {
	_, err := db.Client.NewCreateTable().
		Model((*DataPoint)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreateDataPoint(datapoint *DataPoint, ctx iris.Context) error {
	if datapoint == nil {
		return errors.New("missing or invalid data point")
	}

	if err := InitAnalytics(ctx); err != nil {
		log.Print(err)
	}

	_, err := db.Client.NewInsert().
		Model(datapoint).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
