package web

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	_ "github.com/lib/pq"
)

type TablesResp struct {
	Tables []string `json:"tables"`
}

type TablesCtrl struct {
	Namespace string
	Index     interface{} `path:"" method:"GET"`
}

func (ctrl *TablesCtrl) IndexFunc(ctx echo.Context) error {
	var tablesResp TablesResp
	sesn, _ := session.Get("session", ctx)
	dbo := sesn.Values["dbo"].(*sql.DB)
	rows, err := dbo.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	appendRows(rows, &tablesResp)

	return ctx.JSON(http.StatusOK, tablesResp)
}

func appendRows(rows *sql.Rows, tablesResp *TablesResp) {
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tablesResp.Tables = append(tablesResp.Tables, name)
	}
}