package interceptors

import (
	"github.com/sirupsen/logrus"

	filter "github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
)

type textInterceptor struct {
	words  []string
	filter filter.DirtyFilter
}

func (t *textInterceptor) Enabled() bool {
	return len(t.words) > 0
}

func (t *textInterceptor) IsSensitive(message string) (sensitive bool) {
	result, err := t.filter.Filter(message, '*', '@', '$', '%', '&')
	if err != nil {
		logrus.Errorf("Text intercept failed: %s", err)
		return
	}

	if result != nil {
		sensitive = true
	}
	return
}

func LoadSensitiveWords(words []string) {
	TextInterceptor = textInterceptor{words: words}

	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: words,
	})
	if err != nil {
		panic(err)
	}

	TextInterceptor.filter = filter.NewDirtyManager(memStore).Filter()
}
