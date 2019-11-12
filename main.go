package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/interceptors"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/services"
)

func main() {
	service := flag.String("service", "http", "run a service")
	dir := flag.String("dir", "./", "config.yaml dir")
	flag.Parse()

	config.LoadConfig(*dir)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.AppConfig.Database.DatebaseUser,
		config.AppConfig.Database.DatabasePassword,
		config.AppConfig.Database.DatabaseHost,
		config.AppConfig.Database.DatabasePort,
		config.AppConfig.Database.DatabaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panicln(err)
	}
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(128)
	db.SetMaxIdleConns(4)

	defer db.Close()
	database, err := durable.NewDatabase(context.Background(), db)
	if err != nil {
		log.Panicln(err)
	}

	interceptors.LoadInterceptors()

	plugin.LoadPlugins(database)

	switch *service {
	case "http":
		if config.AppConfig.System.AccpetWeChatPayment {
			go services.StartWxPaymentWatch(*service, database)
		}
		if config.AppConfig.System.AutoEstimate {
			go services.StartCurrencySync(*service, database)
		}
		if config.AppConfig.System.RewardsEnable {
			go services.StartRank(*service, database, connStr)
		}
		err := StartServer(database)
		if err != nil {
			log.Println(err)
		}
	default:
		go func() {
			hub := services.NewHub(database)
			err := hub.StartService(*service)
			if err != nil {
				log.Println(err)
			}
		}()
		http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.Service.HTTPListenPort+2000), http.DefaultServeMux)
	}
}
