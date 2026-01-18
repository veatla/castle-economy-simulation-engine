package manager

import (
	"example/hello/src/world"
	"sync"
)

type WorldManager struct {
	worlds map[world.WorldID]*world.World
	mu     sync.Mutex
}
