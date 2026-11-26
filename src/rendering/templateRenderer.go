package rendering

import "github.com/gin-contrib/multitemplate"

func CreateRenderer() multitemplate.Renderer {
	const layoutPath = "templates/layout.html"

	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", layoutPath, "templates/login.html")
	r.AddFromFiles("search", layoutPath, "templates/search.html")
	r.AddFromFiles("register", layoutPath, "templates/register.html")
	r.AddFromFiles("weather", layoutPath, "templates/weather.html")
	return r
}
