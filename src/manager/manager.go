package manager

import (
	"sync"
	"veatla/simulator/src/world"
)

type WorldManager struct {
	worlds map[world.WorldID]*world.World
	mu     sync.Mutex
}
