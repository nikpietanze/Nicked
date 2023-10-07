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
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	price, err := models.GetPrice(&id, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	ctx.JSON(price)
}

func CreatePrice(ctx iris.Context) {
	var priceJSON models.Price
	err := ctx.ReadJSON(&priceJSON)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("missing or invalid price data").DetailErr(err))
	}

	price, err := models.CreatePrice(&priceJSON, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    ctx.JSON(price)
}

func DeletePrice(ctx iris.Context) {
	strId := ctx.Params().Get("id")
	if strId == "" {
		ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
			Title("Missing or invalid price id"))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	deleteErr := models.DeletePrice(id, ctx)
	if deleteErr != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}
}
