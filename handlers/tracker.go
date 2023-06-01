package handlers

import (
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"

	"pricetracker/db"
	"pricetracker/models"
)

func GetTracker(ctx iris.Context) {
    models.InitTracker(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid account id"))
        return
    }

    id, err := strconv.Atoi(strId)
    if err != nil {
        panic(err)
    }

    tracker := new(models.Tracker)
    err = db.Client.NewSelect().
        Model(tracker).
        Where("id = ?", id).
        Scan(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(tracker)
}

func CreateTracker(ctx iris.Context) {
    models.InitTracker(ctx)

    var tracker models.Tracker
    err := ctx.ReadJSON(&tracker)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("Missing or invalid tracker data").DetailErr(err))
        return
    }

    tracker.Email = strings.ToLower(tracker.Email)
    tracker.Active = true

    res, err := db.Client.NewInsert().
        Model(&tracker).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(res)
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

    res, err := db.Client.NewUpdate().
        Model(&tracker).
        Column("name", "url", "email", "active").
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

func DeleteTracker(ctx iris.Context) {
    models.InitTracker(ctx)

    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("Missing or invalid tracker id"))
        return
    }

    id, err := strconv.Atoi(strId)
    if err != nil {
        panic(err)
    }

    tracker := models.Tracker {
        Id: &id,
    }

    _, err = db.Client.NewDelete().
        Model(&tracker).
        Where("id = ?", id).
        Exec(ctx)
    if err != nil {
        panic(err)
    }

    ctx.JSON(true)
}
