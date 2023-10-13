package handlers

import (
	"github.com/kataras/iris/v12"

	"Nicked/models"
)

func CreateDataPoint(ctx iris.Context) {
	var dataPointJSON models.DataPoint
	if err := ctx.ReadJSON(&dataPointJSON); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("missing or invalid data point").DetailErr(err))
	}

    if err := models.CreateDataPoint(&dataPointJSON, ctx); err != nil {
		ctx.StopWithJSON(500, models.NewError(err))
	}

    ctx.ResponseWriter().WriteHeader(200)
}
