package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
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
	r := gin.Default()
	m := melody.New()

	r.GET("/messages", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
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
	})

	r.Run(":7023")
}
