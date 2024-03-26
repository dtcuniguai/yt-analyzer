package route

import (
	"ytanalyzer/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoute(app *fiber.App) {

	//access log
	// app.Use(middleware.AccessLog)

	//handler init
	var oauth handlers.Oauth

	//oauth api
	app.Get("oauth/access", oauth.Auth)
	app.Get("oauth/redirect", oauth.Redirect)

}
