package rendering

import "github.com/gin-contrib/multitemplate"

func CreateRenderer() multitemplate.Renderer {
	LAYOUT_PATH := "templates/layout.html"

	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", LAYOUT_PATH, "templates/login.html")
	r.AddFromFiles("search", LAYOUT_PATH, "templates/search.html")
	r.AddFromFiles("register", LAYOUT_PATH, "templates/register.html")
	r.AddFromFiles("weather", LAYOUT_PATH, "templates/weather.html")
	return r
}
