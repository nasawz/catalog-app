package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// //go:embed pb_public
// var content embed.FS

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		subFs := echo.MustSubFS(e.Router.Filesystem, "pb_public")
		e.Router.GET("/*", apis.StaticDirectoryHandler(subFs, false))

		// redirect to homepage on missing file
		originalErrorHandler := e.Router.HTTPErrorHandler
		e.Router.HTTPErrorHandler = func(c echo.Context, err error) {
			if c.Path() == "/*" && err == echo.ErrNotFound {
				err = c.Redirect(http.StatusTemporaryRedirect, "/")
			}
			originalErrorHandler(c, err)
		}

		return nil
	})

	// app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
	// 	// serves static files from the provided public dir (if exists)
	// 	subFs := echo.MustSubFS(content, "pb_public")
	// 	e.Router.GET("/*", apis.StaticDirectoryHandler(subFs, false))

	// 	// redirect to homepage on missing file
	// 	originalErrorHandler := e.Router.HTTPErrorHandler
	// 	e.Router.HTTPErrorHandler = func(c echo.Context, err error) {
	// 		if c.Path() == "/*" && err == echo.ErrNotFound {
	// 			err = c.Redirect(http.StatusTemporaryRedirect, "/")
	// 		}
	// 		originalErrorHandler(c, err)
	// 	}

	// 	return nil
	// })

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /api/hello" route to the app router (echo)
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/hello",
			Handler: func(c echo.Context) error {
				type User struct {
					Name  string `json:"name" xml:"name"`
					Email string `json:"email" xml:"email"`
				}
				u := &User{
					Name:  "Jon",
					Email: "jon@labstack.com",
				}
				return c.JSON(http.StatusOK, u)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireGuestOnly(),
			},
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
