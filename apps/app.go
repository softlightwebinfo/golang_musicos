package apps

import (
	_ "../docs"
	"../graphiql"
	"../libs"
	"../middlewares"
	"../models"
	"../routers"
	"../settings"
	"fmt"
	"github.com/adejoux/grafanaclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag/example/celler/controller"
	"gopkg.in/robfig/cron.v3"
	"log"
	"time"
)

type App struct {
	g *gin.Engine
	c *controller.Controller
}

func (a *App) Initialize() {
	a.config()
	a.g = gin.New()
	a.Middleware()
	//a.g.MaxMultipartMemory = 8 << 20 // 8 MiB
	a.g.Use(gin.Logger())
	a.g.Use(gin.Recovery())
	a.g.Static("/assets", "./assets")
	a.g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//a.CronJob()
	v1 := a.g.Group("/api/v1/")
	routers.HomeRoute(v1)
}
func (a *App) Middleware() {
	a.g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3002", "http://localhost:8000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowFiles:       true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))
	//a.g.Use(middlewares.CORSMiddleware())
	a.g.Use(middlewares.IsAuthorized())
}
func (a *App) Run(addr string) {
	fmt.Println("Start app restful")
	log.Fatal(a.g.Run(addr))
	//log.Fatal(a.g.RunTLS(addr, "./ssl/server.pem", "./ssl/server.key"))
}
func (a *App) Grafana() {
	// create Grafana session
	config := settings.GetGrafanaConfig()
	session := grafanaclient.NewSession(config.User, config.Password, config.Url)
	// logon on the current session
	session.DoLogon()
	// let's list the existing data sources
	dataSources, listErr := session.GetDataSourceList()
	if listErr != nil {
		log.Fatal(listErr)
	}

	for _, dataSource := range dataSources {
		fmt.Printf("name: %s type: %s url: %s\n", dataSource.Name, dataSource.Type, dataSource.URL)
	}

}
func (a *App) config() {
	environment := settings.GetEnvironment()
	if environment == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func (a *App) Graphql() {
	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name:        "RootQuery",
		Description: "Root Query",
		Fields: graphql.Fields{
			"auth":  graphiql.GraphQueAuth(),
			"hello": graphiql.GraphHello(),
		},
	})
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"auth": graphiql.GraphMutAuth(),
		},
	})
	var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	model := libs.GraphqlModel{}
	model.Init(a.g, &Schema)
}

func (a *App) CronJob() {
	c := cron.New()
	_, _ = c.AddFunc("@every 6h30m", func() {
		settings.InstanceDb = libs.GetConnection()
		sitemap := libs.Sitemap{}
		sitemap.New()
		sitemap.Generate()
		sitemap.Save()
		defer settings.InstanceDb.Close()
	})
	_, _ = c.AddFunc("@every 1d", func() {
		settings.InstanceDb = libs.GetConnection()
		models.CronRenovated()
		defer settings.InstanceDb.Close()
	})
	c.Start()
}
