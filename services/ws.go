package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/config"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"gopkg.in/olahol/melody.v1"
)

type WsBroadcastMessage struct {
	MessageId     string                       `json:"id"`
	SpeakerName   string                       `json:"speaker_name"`
	SpeakerAvatar string                       `json:"speaker_avatar"`
	SpeakerId     string                       `json:"speaker_id"`
	Category      string                       `json:"category"`
	Data          string                       `json:"data"`
	Text          string                       `json:"text"`
	Attachment    WsBroadcastMessageAttachment `json:"attachment"`
	CreatedAt     time.Time                    `json:"created_at"`
}

type WsBroadcastMessageAttachment struct {
	ID        string `json:"id"`
	Size      int    `json:"size"`
	MimeType  string `json:"mime_type"`
	Persisted bool   `json:"-"`

	Name      *string `json:"name,omitempty"`
	Duration  *uint   `json:"duration,omitempty"`
	Waveform  []byte  `json:"waveform,omitempty"`
	Width     *uint   `json:"width,omitempty"`
	Height    *uint   `json:"height,omitempty"`
	Thumbnail []byte  `json:"thumbnail,omitempty"`

	ViewUrl string `json:"view_url"`
}

func StartWebsocketService(name string, db *durable.Database, broadcastChan chan WsBroadcastMessage) {
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
