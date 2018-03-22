package web

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/satori/go.uuid"
)

type QueriesCtrl struct {
	Namespace string
	Create    interface{} `path:"" method:"POST"`
}

type Query struct {
	String string `json:"query"`
}

func (ctrl *QueriesCtrl) CreateFunc(ctx echo.Context) error {
	var query Query
	sesn, _ := session.Get("session", ctx)
	uuidStr, ok := sesn.Values["uuid"].(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, false)
	}
	uuid, err := uuid.FromString(uuidStr)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, false)
	}
	dbo := DBObjects[uuid]
	json.NewDecoder(ctx.Request().Body).Decode(&query)
	rows, err := dbo.Query(query.String)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	trs := &TableRows{}

	for rows.Next() {
		mapPGRowToTableRow(rows, trs) // at this point in time, represents the cursor at a specific row
	}
	return ctx.JSON(http.StatusOK, trs)
}