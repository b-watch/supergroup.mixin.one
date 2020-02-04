package models

import (
	"github.com/MixinNetwork/supergroup.mixin.one/config"
)

const (
	ResIconBroadcast = "broadcast"
	ResIconDefault   = "default"
)

func ResGetIcon(iconName string) string {
	host := config.AppConfig.Service.HTTPResourceHost
	switch iconName {
	case ResIconBroadcast:
		return host + "/icons/broadcast.png"
	default:
		return host + "/icons/default.png"
	}
}
