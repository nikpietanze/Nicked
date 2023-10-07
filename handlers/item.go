package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

func GetItem(ctx iris.Context) {
	strId := ctx.Params().Get("id")
	if strId == "" {
		ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
			Title("missing or invalid item id"))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	item, getErr := models.GetItem(&id, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(getErr))
	}

	ctx.JSON(item)
}

func CreateItem(ctx iris.Context) {
	var itemJSON models.Item
	err := ctx.ReadJSON(&itemJSON)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("missing or invalid item data").DetailErr(err))
	}

	item, err := models.CreateItem(&itemJSON, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	ctx.JSON(item)
}

func UpdateItem(ctx iris.Context) {
	var itemJSON models.Item
	err := ctx.ReadJSON(&itemJSON)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("missing or invalid item data").DetailErr(err))
	}

	item, err := models.UpdateItem(&itemJSON, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    ctx.JSON(item)
}

func DeleteItem(ctx iris.Context) {
	strId := ctx.Params().Get("id")
	if strId == "" {
		ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
			Title("missing or invalid item id"))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	deleteErr := models.DeleteItem(id, ctx)
	if deleteErr != nil {
		ctx.StopWithJSON(500, models.NewError(deleteErr))
	}
}
