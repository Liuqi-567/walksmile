package main

//go:generate go-bindata -pkg views  -o asset/views/views_gen.go views/...
//go:generate go-bindata-assetfs -pkg static -o asset/static/static_gen.go  static/...
import (
	"flag"
	"html/template"
	"net/http"
	"strings"
	staticinternal "walksmile/asset/static"
	"walksmile/asset/views"
	"walksmile/conf"
	"walksmile/controllers"
	"walksmile/helpers/baiduyun"
	"walksmile/models"

	"github.com/DeanThompson/ginpprof"
	"github.com/claudiu/gocron"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	cfgFile := flag.String("conf", "conf/config.toml", "config file")
	flag.Parse()
	conf.CfgFile = *cfgFile
	conf.InitConfig()
	models.InitDB()
	controllers.InitDATA()
	gocron.Every(1).Day().At("01:00").Do(controllers.FetchHotword)
	gocron.Every(1).Hours().Do(controllers.Last24Update)
	gocron.Every(3).Minutes().Do(baiduyun.Fetcher)
	gocron.Start()
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"isActive":       controllers.IsActive,
		"stringInSlice":  controllers.StringInSlice,
		"dateTime":       controllers.DateTime,
		"pagetofrom":     controllers.PageToFrom,
		"GoodPageToFrom": controllers.GoodPageToFrom,
		"ToTemplate":     controllers.ToTemplate,
	})
	if gin.Mode() == gin.DebugMode {
		views.InitTemplate(true, r.FuncMap)
	} else {
		views.InitTemplate(false, r.FuncMap)
	}
	r.SetHTMLTemplate(views.Temp)

	r.Use(controllers.Domain(), controllers.IsMobile())
	r.Use(static.Serve("/static", BinaryFileSystem("static")))
	r.Static("/static", "./static")
	r.GET("/", controllers.Home)
	r.GET("/search", controllers.FileList)
	//r.GET("/last24", controllers.Last24hour)
	r.GET("/sharedfile/:id", controllers.FileDetail)
	r.GET("/wallpaper/:id", controllers.Wallpaper)
	r.POST("/comment", controllers.Comment)
	ginpprof.Wrapper(r)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080

}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{staticinternal.Asset, staticinternal.AssetDir, staticinternal.AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}
