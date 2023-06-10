package handlers

import (
	"context"
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

func GetItem(ctx context.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid item id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.GetItem(&id, ctx))
}

func CreateItem(ctx iris.Context) {
    // TODO: verify if this works
    var item models.Item
    err := ctx.ReadJSON(&item)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid item data").DetailErr(err))
        return
    }

    for _, v := range item.Merchants {
        models.CreateMerchant(v, ctx)
	}


    ctx.JSON(models.CreateItem(&item, ctx))
}

func UpdateItem(ctx iris.Context) {
    var item models.Item
    err := ctx.ReadJSON(&item)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid item data").DetailErr(err))
        return
    }

    ctx.JSON(models.UpdateItem(&item, ctx))
}

func DeleteItem(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid item id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.DeleteItem(&id, ctx))
}

