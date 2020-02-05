package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
	"github.com/MixinNetwork/supergroup.mixin.one/models"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"gopkg.in/olahol/melody.v1"
)

func StartWebsocketService(name string, db *durable.Database, broadcastChan chan models.WsBroadcastMessage) {
	log.Println("Init websocket service")
	m := melody.New()

	go func() {
		for {
			select {
			case msg := <-broadcastChan:
				bts, err := json.Marshal(msg)
				if err != nil {
					fmt.Printf("StartWebsockService: %s\n", err)
					return
				}

				m.Broadcast(bts)
			}
		}
	}()

	http.ListenAndServe(":"+strconv.Itoa(config.AppConfig.Service.HTTPWebsocketPort), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	}))
}
