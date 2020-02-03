package interceptors

import (
	"strings"
	"sync"

	"github.com/MixinNetwork/supergroup.mixin.one/config"
)

var (
	loadOnce        sync.Once
	TextInterceptor textInterceptor
	LinkInterceptor linkInterceptor
)

func LoadInterceptors() {
	loadOnce.Do(func() {
		sensitiveWords := strings.Split(config.AppConfig.System.SensitiveWords, "|")
		LoadSensitiveWords(sensitiveWords)
		whitelistedHosts := config.AppConfig.System.DetectLinkWhitelist
		LoadWhitelistedHosts(whitelistedHosts)
	})
}
