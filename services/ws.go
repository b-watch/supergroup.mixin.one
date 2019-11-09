package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
)

type WsBroadcastMessage struct {
	MessageId     string    `json:"id"`
	SpeakerName   string    `json:"speaker_name"`
	SpeakerAvatar string    `json:"speaker_avatar"`
	SpeakerId     string    `json:"speaker_id"`
	Category      string    `json:"category"`
	Data          string    `json:"data"`
	CreatedAt     time.Time `json:"created_at"`
}

func StartWebsocketService(name string, db *durable.Database, broadcastChan chan WsBroadcastMessage) {
	log.Println("Init websocket service")
	http.ListenAndServe(":7023", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		log.Println("handle conn")

		go func() {
			defer conn.Close()

			for {
				log.Println("check chan")
				select {
				case msg := <-broadcastChan:
					log.Println("new message:", msg)
					bts, err := json.Marshal(msg)
					if err != nil {
						fmt.Printf("StartWebsockService: %s\n", err)
						return
					}
					err = wsutil.WriteServerMessage(conn, ws.OpText, []byte(bts))
					if err != nil {
						// handle error
					}
					// default:
					// log.Println("no activity")
				}
				// time.Sleep(3000 * time.Millisecond)
			}
		}()
	}))
}
