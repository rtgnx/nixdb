package main

import (
	"log"
	"os"

	"github.com/Reverse-Labs/nixdb"
	cli "github.com/jawher/mow.cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e   = echo.New()
	app = cli.App("nixdbd", "Nix Database Daemon")
)

func main() {
	app.Command("serve", "start http server", cmdServe)
	app.Command("genkey", "generate secret key", cmdGenKey)

	app.Run(os.Args)
}

func registerEndpoints(e *echo.Echo, db nixdb.Database, secret JWTSecretKey) {
	api := e.Group("/v1/api")

	config := middleware.JWTConfig{Claims: &JWTClaim{}, SigningKey: []byte(secret)}
	api.Use(middleware.JWTWithConfig(config))

	H := HTTP{db: db}

	e.GET("login", H.GETLogin)

	api.GET("/passwd", H.GETUsers)
	api.GET("/group", H.GETGroups)
	api.GET("/hosts", H.GETHosts)
}

func cmdServe(cmd *cli.Cmd) {
	var (
		addr    = cmd.StringOpt("addr", ":8080", "bind address")
		baseDir = cmd.StringOpt("baseDir", "/etc", "directory containing passwd group and hosts files")
		minUID  = cmd.IntOpt("minUID", 1000, "minimum uid")
		minGID  = cmd.IntOpt("minGID", 1000, "minimum gid")
		groups  = cmd.StringsOpt("groups", []string{"testuser"}, "comma separated list of authorized groups")

		secretKey = cmd.StringOpt("secretKey", "/etc/nixdb.key", "path to secret key")

		autoTLS   = cmd.BoolOpt("autoTLS", false, "enable auto tls (Let' Encrypt)")
		enableTLS = cmd.BoolOpt("enableTLS", false, "enable tls")
		tlsCert   = cmd.StringOpt("tlsCert", "", "path to tls certificate")
		tlsKey    = cmd.StringOpt("tlsKey", "", "path to tls key")
	)

	cmd.Action = func() {

		key, err := ReadOrGenerateJWTSecret(*secretKey)

		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Printf("[*] Serving Directory: %s", *baseDir)
		log.Printf("[*] Min UID: %d, Min GID: %d", *minUID, *minGID)

		e.Use(middleware.Recover())
		e.Use(middleware.Logger())
		db := nixdb.NewDB(*baseDir, uint(*minUID), uint(*minGID))

		if err := db.Update(); err != nil {
			log.Fatalln(err.Error())
		}

		e.Use(AuthMiddleware(key, db, *groups))

		registerEndpoints(e, db, key)

		switch true {
		case *autoTLS:
			e.Logger.Info(e.StartAutoTLS(*addr))
		case *enableTLS:
			e.Logger.Info(e.StartTLS(*addr, *tlsCert, *tlsKey))
		default:
			e.Logger.Info(e.Start(*addr))
		}

	}
}

func cmdGenKey(cmd *cli.Cmd) {
	var (
		secretKey  = cmd.StringOpt("secretKey", "/etc/nixdb.key", "path to secret key")
		secretSize = cmd.IntOpt("secretSize", 4096, "bit size of secret key")
	)

	cmd.Action = func() {
		WriteNewKey(*secretKey, uint(*secretSize))
	}
}
