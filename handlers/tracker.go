package handlers

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

type CreateTrackerBody struct {
    Name string `json:"name"`
    Url string `json:"url"`
    Email string `json:"email"`
}

func GetTracker(ctx iris.Context) {
    models.InitTracker(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid account id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.GetTracker(&id, ctx))
}

func CreateTracker(ctx iris.Context) {
    var reqBody CreateTrackerBody
    err := ctx.ReadJSON(&reqBody)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid tracker data").DetailErr(err))
        return
    }

    item := models.Item{
        Name: reqBody.Name,
    }
    itemId := models.CreateItem(&item, ctx)

    fmt.Printf("itemId: %v\n", itemId)

    merchant := models.Merchant{
        Name: models.DetermineMerchant(reqBody.Url),
        Url: reqBody.Url,
        ItemId: &itemId,
    }
    models.CreateMerchant(&merchant, ctx)

    tracker := models.Tracker{
        Email: reqBody.Email,
        ItemId: &itemId,
    }
    ctx.JSON(models.CreateTracker(&tracker, ctx))
}

func UpdateTracker(ctx iris.Context) {
    models.InitTracker(ctx)

    var tracker models.Tracker
    err := ctx.ReadJSON(&tracker)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid tracker data").DetailErr(err))
        return
    }

    ctx.JSON(models.UpdateTracker(&tracker, ctx))
}

func DeleteTracker(ctx iris.Context) {
    models.InitTracker(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid tracker id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.DeleteTracker(&id, ctx))
}
