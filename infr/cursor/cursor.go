package cursor

import "sync"

var (
	cursorMap sync.Map
	mu        sync.Mutex
)

func NextCursor(key string, cnt int64) int64 {
	mu.Lock()
	defer mu.Unlock()

	v, ok := cursorMap.Load(key)
	if !ok {
		return 0
	}
	x := (v.(int64) + 1) % cnt
	cursorMap.Store(key, x)
	return x
}
