package rendering

import "github.com/gin-contrib/multitemplate"

func CreateRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", "templates/layout.html", "templates/login.html")
	r.AddFromFiles("search", "templates/layout.html", "templates/search.html")
	r.AddFromFiles("register", "templates/layout.html", "templates/register.html")
	r.AddFromFiles("about", "templates/layout.html", "templates/about.html")
	r.AddFromFiles("weather", "templates/layout.html", "templates/weather.html")
	return r
}
