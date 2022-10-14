package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pocketbase/pocketbase/cmd"

	"github.com/pocketbase/pocketbase"
	"github.com/stewie1520/ecommerce-backend/internal/cache"
	"github.com/stewie1520/ecommerce-backend/internal/config"
	service_account "github.com/stewie1520/ecommerce-backend/internal/core/services/account"
	service_hash "github.com/stewie1520/ecommerce-backend/internal/core/services/hash"
	"github.com/stewie1520/ecommerce-backend/internal/infrastructure/web/routes"
	"github.com/stewie1520/ecommerce-backend/internal/tools/migration"
	"github.com/stewie1520/ecommerce-backend/internal/tools/path"
)

func init() {
	fmt.Printf("✈️ App starting in %s mode\n", config.AppConfig.Environment)

	err := cache.InitRedisClient(config.AppConfig.Redis.Host, config.AppConfig.Redis.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: path.PathCWD("data"),
	})

	hashService := service_hash.NewBcryptHashService()
	accountService := service_account.NewAccountService(cache.RedisClient, hashService)

	routes.RegisterRoutes(app, accountService)

	if err := app.Bootstrap(); err != nil {
		log.Fatalln(err)
	}

	migration.RunMigrations(app)
	config.OverwritePocketBaseConfig(app)

	gratefulStartApp(app)
}

func gratefulStartApp(app *pocketbase.PocketBase) error {
	// clear all args, if we want to customize pocketbase with any args, it should be written in code
	cmdServe := cmd.NewServeCommand(app, true)
	os.Args = os.Args[:1]

	var wg sync.WaitGroup

	wg.Add(1)

	// wait for interrupt signal to gracefully shutdown the application
	go func() {
		defer wg.Done()
		quit := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
	}()

	go func() {
		defer wg.Done()
		if err := cmdServe.Execute(); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	//TODO: cleanup
	return nil
}
