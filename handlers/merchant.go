package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
	"pricetracker/models"
)

func GetMerchant(ctx iris.Context) {
    models.InitMerchant(ctx)

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
    models.InitMerchant(ctx)

    var merchant models.Merchant
    err := ctx.ReadJSON(&merchant)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid merchant data").DetailErr(err))
        return
    }

    merchant.Name = strings.Title(merchant.Name)
    merchant.Url = strings.ToLower(merchant.Url)

    res, err := db.Client.NewInsert().
        Model(&merchant).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(res)
}

func UpdateMerchant(ctx iris.Context) {
    models.InitMerchant(ctx)

    var merchant models.Merchant
    err := ctx.ReadJSON(&merchant)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid merchant data").DetailErr(err))
        return
    }

    merchant.UpdatedAt = time.Now()

    res, err := db.Client.NewUpdate().
        Model(&merchant).
        Column("name", "item_id", "updated_at").
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

func DeleteMerchant(ctx iris.Context) {
    models.InitMerchant(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid merchant id"))
        return
    }

    id, err := strconv.Atoi(strId)
    if err != nil {
        panic(err)
    }

    merchant := models.Merchant {
        Id: &id,
    }

    _, err = db.Client.NewDelete().
        Model(&merchant).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(true)
}
