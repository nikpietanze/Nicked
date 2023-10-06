package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

func GetPrice(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("missing or invalid price id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.GetPrice(&id, ctx))
}

func CreatePrice(ctx iris.Context) {
    var price models.Price
    err := ctx.ReadJSON(&price)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("missing or invalid price data").DetailErr(err))
        return
    }

    ctx.JSON(models.CreatePrice(&price, ctx))
}

func UpdatePrice(ctx iris.Context) {
    var price models.Price
    err := ctx.ReadJSON(&price)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("missing or invalid price data").DetailErr(err))
        return
    }

    ctx.JSON(models.UpdatePrice(&price, ctx))
}

func DeletePrice(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid price id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.DeletePrice(&id, ctx))
}

