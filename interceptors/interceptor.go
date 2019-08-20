package interceptors

import (
	"strings"
	"sync"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
)

var (
	loadOnce        sync.Once
	TextInterceptor textInterceptor
)

func LoadInterceptors() {
	loadOnce.Do(func() {
		sensitiveWords := strings.Split(config.AppConfig.System.SensitiveWords, "|")
		LoadSensitiveWords(sensitiveWords)
	})
}
