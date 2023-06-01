package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

func GetItem(ctx iris.Context) {
    models.InitItem(ctx)

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
    models.InitItem(ctx)

    var item models.Item
    err := ctx.ReadJSON(&item)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid item data").DetailErr(err))
        return
    }

    ctx.JSON(models.CreateItem(&item, ctx))
}

func UpdateItem(ctx iris.Context) {
    models.InitItem(ctx)

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
    models.InitItem(ctx)

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

