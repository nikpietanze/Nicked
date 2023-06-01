package handlers

import (
	"strconv"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
	"pricetracker/models"
)

func GetPrice(ctx iris.Context) {
    models.InitPrice(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid price id"))
        return
    }

    id, err := strconv.Atoi(strId)
    if err != nil {
        panic(err)
    }

    price := new(models.Price)
    err = db.Client.NewSelect().
        Model(&price).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(price)
}

func CreatePrice(ctx iris.Context) {
    models.InitPrice(ctx)

    var price models.Price
    err := ctx.ReadJSON(&price)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid price data").DetailErr(err))
        return
    }

    res, err := db.Client.NewInsert().
        Model(&price).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(res)
}

func UpdatePrice(ctx iris.Context) {
    models.InitPrice(ctx)

    var price models.Price
    err := ctx.ReadJSON(&price)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid price data").DetailErr(err))
        return
    }

    price.UpdatedAt = time.Now()

    res, err := db.Client.NewUpdate().
        Model(&price).
        Column("name", "updated_at").
        OmitZero().
        WherePK().
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        panic(err)
    }

    ctx.JSON(id)
}

func DeletePrice(ctx iris.Context) {
    models.InitPrice(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid price id"))
        return
    }

    id, err := strconv.Atoi(strId)
    if err != nil {
        panic(err)
    }

    price := models.Price {
        Id: &id,
    }

    _, err = db.Client.NewDelete().
        Model(&price).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(true)
}
