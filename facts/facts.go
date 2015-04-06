package facts

import "sync"

type Facts struct {
	Facts map[string]interface{}
	mu    sync.Mutex
}

func New() *Facts {
	m := make(map[string]interface{})
	return &Facts{Facts: m}
}

func (f *Facts) Add(key string, value interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Facts[key] = value
}
