package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"Nicked/models"
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
	type ItemJSON struct {
		Currency string
		Email    string
		Name     string
		Price    string
		Sku      string
		Store    string
		Url      string
	}

	var itemJSON ItemJSON
	if err := ctx.ReadJSON(&itemJSON); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("missing or invalid item data").DetailErr(err))
	}

    user, err := models.GetUserByEmail(itemJSON.Email, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    newItem := models.Item{
        Name: itemJSON.Name,
        IsActive: true,
        Sku: itemJSON.Sku,
        Url: itemJSON.Url,
        UserId: user.Id,
    }

	item, err := models.CreateItem(&newItem, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    priceFlt, err := strconv.ParseFloat(itemJSON.Price, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    newPrice := models.Price{
        Amount: priceFlt,
        Currency: itemJSON.Currency,
        Store: itemJSON.Store,
        ItemId: item.Id,
    }

    _, err = models.CreatePrice(&newPrice, ctx)
    if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	ctx.JSON(item)
}

func UpdateItem(ctx iris.Context) {
	var itemJSON models.Item
	if err := ctx.ReadJSON(&itemJSON); err != nil {
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

	if err := models.DeleteItem(id, ctx); err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}
}
