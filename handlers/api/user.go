package handlers

import (
	"strconv"

	"github.com/kataras/iris/v12"

	"Nicked/models"
)

func GetUser(ctx iris.Context) {
	strId := ctx.Params().Get("id")
	if strId == "" {
		ctx.StopWithProblem(424, iris.NewProblem().
			Title("missing or invalid user id"))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	user, getErr := models.GetUser(&id, ctx)
	if getErr != nil {
		ctx.StopWithJSON(500, models.NewError(getErr))
	}

	ctx.JSON(user)
}

func CreateUser(ctx iris.Context) {
	var userJSON models.User
	if err := ctx.ReadJSON(&userJSON); err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	user, err := models.CreateUser(&userJSON, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    ctx.JSON(user)
}

func UpdateUser(ctx iris.Context) {
	var userJSON models.User
	if err := ctx.ReadJSON(&userJSON); err != nil {
		ctx.StopWithJSON(400, models.NewError(err))
	}

	user, err := models.UpdateUser(&userJSON, ctx)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    ctx.JSON(user)
}

func DeleteUser(ctx iris.Context) {
	strId := ctx.Params().Get("id")
	if strId == "" {
		ctx.StopWithProblem(iris.StatusFailedDependency, iris.NewProblem().
			Title("missing or invalid user id"))
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

	if err := models.DeleteUser(id, ctx); err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}
}
