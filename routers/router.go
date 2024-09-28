package routers

import (
	"catApi/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/api/cat-images", &controllers.CatController{}, "get:GetCatImages")
	beego.Router("/api/breeds", &controllers.CatController{}, "get:GetBreeds")
	beego.Router("/api/cat-images/by-breed", &controllers.CatController{}, "get:GetCatImagesByBreed")

	// Favorites routes
	beego.Router("/api/favorites", &controllers.CatController{}, "post:AddFavorite")
	beego.Router("/api/favorites", &controllers.CatController{}, "get:GetFavorites")
	beego.Router("/api/favorites/:id", &controllers.CatController{}, "delete:DeleteFavorite")

	// New route for voting
	beego.Router("/api/votes", &controllers.CatController{}, "post:Vote")
	beego.Router("/api/votes", &controllers.CatController{}, "get:GetVotes")
	beego.Router("/*", &controllers.MainController{}, "get:Get")
}
