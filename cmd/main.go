package main

import (
	"flag"
	"log"
	_ "project/docs"
	"project/infra"
	"project/routes"
)

// @title
// @version 1.0
// @description
// @termsOfService http://example.com/terms/
// @contact.name Team
// @contact.url https://academy.lumoshive.com/contact-us
// @contact.email lumoshive.academy@gmail.com
// @license.name Lumoshive Academy
// @license.url https://academy.lumoshive.com
// @host localhost:8080
// @schemes http
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	ctx, err := infra.NewServiceContext()
	if err != nil {
		log.Fatal("can't init service context %w", err)
	}

	if shouldNotLaunchServer() {
		return
	}

	routes.NewRoutes(*ctx)
}

func shouldNotLaunchServer() bool {
	shouldNotLaunch := false

	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "m" {
			shouldNotLaunch = true
		}

		if f.Name == "s" {
			shouldNotLaunch = true
		}
	})

	return shouldNotLaunch
}
