package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/istt/api_gateway/internal/app"
	"github.com/istt/api_gateway/internal/app/api-gateway/repositories"
	authImpl "github.com/istt/api_gateway/internal/app/api-gateway/services/impl"
	"github.com/istt/api_gateway/internal/app/s3proxy"
	"github.com/istt/api_gateway/pkg/fiber/authjwt"
	authApi "github.com/istt/api_gateway/pkg/fiber/authjwt/web/rest"
	"github.com/istt/api_gateway/pkg/fiber/fiberprometheus"
	"github.com/markbates/pkger"
)

var configFile string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&configFile, "config", "configs/api-gateway.yaml", "API Gateway configuration file")
	flag.Parse()

	// 1 - set default settings for components.

	// 2 - override defaults with configuration file and watch changes
	app.ConfigInit(configFile)
	app.ConfigWatch(configFile)

	// 3 - bring up components
	// + inject UserServiceMongoDB into application
	userRepo := repositories.NewUserRepositoryDummy()
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

	prometheus := fiberprometheus.New("api-gateway")
	prometheus.RegisterAt(srv, "/metrics")
	srv.Use(prometheus.Middleware)

	authjwt.USER_RESOURCE = authApi.NewDefaultUserResource(userSvc, userRepo)
	authjwt.ACCOUNT_RESOURCE = authApi.NewDefaultAccountResource(userSvc)
	authjwt.SetupAuthJWT(srv, app.Config.MustString("security.jwt-secret"), app.Config.Strings("security.skip-auth")...)
	authjwt.SetupRoutes(srv)

	s3proxy.SetupRoutes(srv)
	setupProxy(srv, prometheus.Middleware)
	if app.Config.String("https.listen") != "" {
		log.Fatal(srv.ListenTLS(app.Config.String("https.listen"), app.Config.MustString("https.cert"), app.Config.MustString("https.key")))
	} else {
		log.Fatal(srv.Listen(app.Config.String("http.listen")))
	}
}

// configureFiber start the fiber with common settings
func configureFiber(srv *fiber.App) {
	staticAsset := filesystem.New(filesystem.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api") || strings.HasPrefix(c.Path(), "/services")
		},
		Root: pkger.Dir("/web"),
	})
	srv.Use(staticAsset)
	// + logging
	srv.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05-0700",
	}))

	srv.Use(recover.New())
	srv.Get("/management/info", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusOK, "ok")
	})
	srv.Get("/management/health", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusOK, "ok")
	})
}

// setupProxy setup the reverse proxy based on
func setupProxy(srv *fiber.App, handlers interface{}) {
	for k, v := range app.Config.StringsMap("services") {
		serviceNamespace := fmt.Sprintf("/services/%s/", strings.TrimRight(k, "/"))
		log.Printf("proxy request on %s to %v", serviceNamespace, v)
		srv.Use(serviceNamespace, handlers, proxy.Balancer(proxy.Config{
			Servers: v,
			ModifyRequest: func(c *fiber.Ctx) error {
				c.Request().Header.Add("X-Real-IP", c.IP())
				c.Path(c.Path()[len(serviceNamespace):])
				// FIXME: not sure why but for uploading files, this must be set
				if strings.Contains(string(c.Request().Header.ContentType()), fiber.MIMEMultipartForm) {
					b := c.Body()
					c.Request().SetBodyRaw(b)
				}
				return nil
			},
			ModifyResponse: func(c *fiber.Ctx) error {
				c.Response().Header.Del(fiber.HeaderServer)
				return nil
			},
		}))
	}
}
