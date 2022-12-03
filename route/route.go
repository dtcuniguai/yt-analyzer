package route

import (
	"ytanalyzer/app/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoute(app *fiber.App) {

	//access log
	// app.Use(middleware.AccessLog)

	//handler init
	var oauth handler.Oauth

	//oauth api
	app.Get("oauth/access", oauth.Auth)
	app.Get("oauth/redirect", oauth.Redirect)

}
