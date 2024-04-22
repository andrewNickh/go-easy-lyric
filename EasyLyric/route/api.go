package route

import (
	"easy-lyric/EasyLyric/controller"
	"easy-lyric/EasyLyric/model/response"
	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app *iris.Application) {
	//opts := basicauth.Options{
	//	Allow: basicauth.AllowUsers(map[string]string{
	//		config.Instance.Username: config.Instance.Password,
	//	}),
	//	Realm:        "Authorization Required",
	//	ErrorHandler: basicauth.DefaultErrorHandler,
	//}
	//auth := basicauth.New(opts)
	//app.Use(auth)

	mainGroup := app.Party("/api")
	mainGroup.Get("/ping", func(ctx iris.Context) {
		response.OkWithMessageV2("ok", "ok", ctx)
	})

	// api/scrape
	scrapeGroup := mainGroup.Party("/scrape")
	{
		scrapeGroup.Post("/lyric", controller.ScrapController.GetLyrics)
	}

	// api/song
	songGroup := mainGroup.Party("/song")
	{
		songGroup.Get("/search", controller.SongController.SearchSong)
		songGroup.Post("/create", controller.SongController.CreateSong)
		songGroup.Post("/update", controller.SongController.UpdateSong)
		songGroup.Post("/delete", controller.SongController.DeleteSong)
	}

	// api/resource
	resourceGroup := mainGroup.Party("/resource")
	{
		resourceGroup.Get("/list", controller.ResourceController.GetResourceList)
		resourceGroup.Get("/detail", controller.ResourceController.GetResource)
		resourceGroup.Post("/create", controller.ResourceController.CreateResource)
		resourceGroup.Post("/update", controller.ResourceController.UpdateResource)
		resourceGroup.Post("/delete", controller.ResourceController.DeleteResource)
	}
}
