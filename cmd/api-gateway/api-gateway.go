package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/repositories"
	authImpl "github.com/istt/api_gateway/internal/app/api-gateway/services/impl"
	"github.com/istt/api_gateway/internal/app/api-gateway/web/rest"
	"github.com/istt/api_gateway/pkg/fiber/authjwt"
	authApi "github.com/istt/api_gateway/pkg/fiber/authjwt/web/rest"
	"github.com/istt/api_gateway/pkg/fiber/middleware/filter"
	"github.com/markbates/pkger"
)

var configFile string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&configFile, "config", "configs/api-gateway.yaml", "API Gateway configuration file")
	flag.Parse()

	// 1 - set default settings for components.
	app.MongoDBConfig()

	// 2 - override defaults with configuration file and watch changes
	app.ConfigInit(configFile)
	app.ConfigWatch(configFile)

	// 3 - bring up components
	app.MongoDBInit()
	// + inject UserServiceMongoDB into application
	userRepo := repositories.NewUserRepositoryBuntDB(app.Config.MustString("buntdb.path"))
	userSvc := authImpl.NewUserServiceDummy()
	// 4 - setup the web server
	srv := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				return c.Status(code).JSON(e)
			}
			return c.Status(code).JSON(fiber.Map{"error": code, "message": err.Error()})
		},
	})
	configureFiber(srv)

	// + jwt secret support
	authjwt.USER_RESOURCE = authApi.NewDefaultUserResource(userSvc, userRepo)
	authjwt.ACCOUNT_RESOURCE = authApi.NewDefaultAccountResource(userSvc)
	authjwt.SetupAuthJWT(srv, app.Config.MustString("security.jwt-secret"), app.Config.Strings("security.skip-auth")...)
	authjwt.SetupRoutes(srv)
	setupRoutes(srv)
	setupProxy(srv)

	log.Fatal(srv.Listen(app.Config.String("http.listen")))
}

// configureFiber start the fiber with common settings
func configureFiber(srv *fiber.App) {
	staticAsset := filesystem.New(filesystem.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
		Root: pkger.Dir("/web"),
	})
	srv.Use(staticAsset)
	// + logging
	srv.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
	}))

	// srv.Use(recover.New())
}

// setupProxy setup the reverse proxy based on
func setupProxy(srv *fiber.App) {
	for k, v := range app.Config.StringsMap("http.proxy") {
		log.Printf("proxy request on /%s to %v", k, v)
		srv.Use(k, proxy.Balancer(proxy.Config{
			Servers: v,
			ModifyRequest: func(c *fiber.Ctx) error {
				c.Request().Header.Add("X-Real-IP", c.IP())
				c.Path(c.Path()[len(k)+1:])
				return nil
			},
			ModifyResponse: func(c *fiber.Ctx) error {
				c.Response().Header.Del(fiber.HeaderServer)
				return nil
			},
		}))
	}
}

// setupRoutes setup the route for application
func setupRoutes(app *fiber.App) {
	// + Jhipster endpoint for ROLE_USER
	app.Get("api/account", rest.GetAccount)                                 // getAccount
	app.Post("api/account", rest.SaveAccount)                               // saveAccount
	app.Post("api/account/change-password", rest.ChangePassword)            // ChangePassword
	app.Post("api/account/reset-password/finish", rest.FinishPasswordReset) // finishPasswordReset
	app.Post("api/account/reset-password/init", rest.RequestPasswordReset)  // requestPasswordReset

	// + account public end point
	app.Get("api/activate", rest.ActivateAccount)  // activateAccount
	app.Post("api/authenticate", rest.Login)       // isAuthenticated
	app.Post("api/register", rest.RegisterAccount) // registerAccount

	// + user Management routes
	app.Get("api/authorities", authjwt.HasAnyAuthority("ROLE_ADMIN"), rest.GetAuthorities)
	app.Get("api/admin/users", authjwt.HasAnyAuthority("ROLE_ADMIN"), filter.New(), rest.GetAllUser)
	app.Get("api/admin/users/:id", authjwt.HasAnyAuthority("ROLE_ADMIN"), rest.GetUser)
	app.Post("api/admin/users", authjwt.HasAnyAuthority("ROLE_ADMIN"), rest.CreateUser)
	app.Put("api/admin/users", authjwt.HasAnyAuthority("ROLE_ADMIN"), rest.UpdateUser)
	app.Delete("api/admin/users/:id", authjwt.HasAnyAuthority("ROLE_ADMIN"), rest.DeleteUser)
}
