package web

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	_ "github.com/lib/pq"
	"github.com/nycdavid/ziptie"
	"github.com/satori/go.uuid"
	"github.com/velvetreactor/postapocalypse/testhelper"
)

func TestMain(m *testing.M) {
	dbo, err := sql.Open("postgres", os.Getenv("PGCONN"))
	if err != nil {
		log.Print(err)
		panic(err)
	}
	testhelper.CreateTestTables(dbo)
	testhelper.SeedDb(dbo, "../testhelper/seeds.csv")
	code := m.Run()
	teardownDb(dbo)
	os.Exit(code)
}

func teardownDb(dbo *sql.DB) {
	_, err := dbo.Exec("TRUNCATE items")
	if err != nil {
		log.Print(err)
		panic(err)
	}
}

func getDbo(t *testing.T) *sql.DB {
	pgConn, err := openPgConnection()
	if err != nil {
		t.Error(err)
	}
	return pgConn
}

func openPgConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("PGCONN"))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupSessionStore(e *echo.Echo) *sessions.CookieStore {
	cookieStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	e.Use(session.Middleware(cookieStore))
	return cookieStore
}

func authenticateContext(ctx echo.Context, store *sessions.CookieStore) error {
	ctx.Set("_session_store", store)
	uuid := uuid.NewV4()
	dbo, err := openPgConnection()
	if err != nil {
		return err
	}
	sesn, _ := session.Get("session", ctx)
	sesn.Values["uuid"] = uuid.String()
	sesn.Save(ctx.Request(), ctx.Response())
	DBObjects[uuid] = dbo
	return nil
}

func echoInit(ctrl ziptie.Ctrl) (*echo.Echo, *sessions.CookieStore) {
	e := echo.New()
	ziptie.Fasten(ctrl, e)
	return e, setupSessionStore(e)
}
