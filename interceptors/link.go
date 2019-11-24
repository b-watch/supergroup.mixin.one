package interceptors

import (
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

type linkInterceptor struct {
	linkRe           *regexp.Regexp
	whitelistedHosts []string
}

func (t *linkInterceptor) Enabled() bool {
	return t.linkRe != nil
}

func (t *linkInterceptor) HasExternalLinks(message string) (ret bool) {
	links := t.linkRe.FindAllString(message, -1)
	if len(links) != 0 {
		for _, link := range links {
			isWhitelisted := false
			for _, wlHost := range t.whitelistedHosts {
				// @TODO need a better implement
				if strings.Contains(link, wlHost) {
					isWhitelisted = true
					break
				}
			}
			if !isWhitelisted {
				return true
			}
		}
	}
	return false
}

func LoadWhitelistedHosts(hosts []string) {
	LinkInterceptor = linkInterceptor{}
	LinkInterceptor.whitelistedHosts = hosts
	LinkInterceptor.linkRe = xurls.Relaxed()
}
