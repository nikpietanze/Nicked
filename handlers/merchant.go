package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

func GetMerchant(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid merchant id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.GetMerchant(&id, ctx))
}

func CreateMerchant(ctx iris.Context) {
    var merchant models.Merchant
    err := ctx.ReadJSON(&merchant)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid merchant data").DetailErr(err))
        return
    }

    ctx.JSON(models.CreateMerchant(&merchant, ctx))
}

func UpdateMerchant(ctx iris.Context) {
    var merchant models.Merchant
    err := ctx.ReadJSON(&merchant)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid merchant data").DetailErr(err))
        return
    }

    ctx.JSON(models.UpdateMerchant(&merchant, ctx))
}

func DeleteMerchant(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid merchant id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.DeleteMerchant(&id, ctx))
}
