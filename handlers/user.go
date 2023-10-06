package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"pricetracker/models"
)

func GetUser(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("missing or invalid user id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.GetUser(&id, ctx))
}

func CreateUser(ctx iris.Context) {
    var user models.User
    err := ctx.ReadJSON(&user)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("missing or invalid user data").DetailErr(err))
        return
    }

    ctx.JSON(models.CreateUser(&user, ctx))
}

func UpdateUser(ctx iris.Context) {
    var user models.User
    err := ctx.ReadJSON(&user)
    if err != nil {
        ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
            Title("missing or invalid user data").DetailErr(err))
        return
    }

    ctx.JSON(models.UpdateUser(&user, ctx))
}

func DeleteUser(ctx iris.Context) {
    strId := ctx.Params().Get("id")
    if strId == "" {
        ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
            Title("missing or invalid user id"))
        return
    }

    id, err := strconv.ParseInt(strId, 10, 64)
    if err != nil {
        panic(err)
    }

    ctx.JSON(models.DeleteUser(&id, ctx))
}

