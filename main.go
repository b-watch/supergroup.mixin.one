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

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/interceptors"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
	"github.com/MixinNetwork/supergroup.mixin.one/services"
)

func resetBotPreferences() {
	data := `{"receive_message_source":"EVERYBODY", "accept_conversation_source":"CONTACTS"}`
	token, err := bot.SignAuthenticationToken(config.AppConfig.Mixin.ClientId, config.AppConfig.Mixin.SessionId, config.AppConfig.Mixin.SessionKey, "POST", "/me/preferences", data)
	if err != nil {
		log.Panicln("resetBotPreferences:", err)
	}
	_, err = bot.Request(context.Background(), "POST", "/me/preferences", []byte(data), token)
	if err != nil {
		log.Panicln("resetBotPreferences:", err)
	} else {
		log.Println("resetBotPreferences ... done")
	}
}

func main() {
	service := flag.String("service", "http", "run a service")
	dir := flag.String("dir", "./", "config.yaml dir")
	slug := flag.String("slug", "", "process identity")
	flag.Parse()
	log.Printf("process slug: %s\n", slug)

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

	resetBotPreferences()

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
		WsBroadcastChan := make(chan models.WsBroadcastMessage, 3)
		go services.StartWebsocketService(*service, database, WsBroadcastChan)
		go func() {
			hub := services.NewHub(database, WsBroadcastChan)
			err := hub.StartService(*service)
			if err != nil {
				log.Println(err)
			}
		}()
		http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.Service.HTTPListenPort+2000), http.DefaultServeMux)
	}
}
