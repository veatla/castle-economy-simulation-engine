package agents

import (
	"math"
	"math/rand"
	"time"
	"veatla/simulator/src/utils"
	worldQuery "veatla/simulator/src/world-query"

	"github.com/google/uuid"
)

func CreateSimpleAgent(q worldQuery.WorldQuery) Agent {
	id := uuid.New()
	r := rand.New(rand.NewSource(q.GetWorldSeed() + utils.UUIDToInt64(id)))
	angle := r.Float64() * 2 * math.Pi
	worldWidth, worldHeight := q.GetBoundaries()

	tx := r.Float64() * worldWidth
	tz := r.Float64() * worldHeight

	for q.IsPointBlocked(tx, tz) {
		tx = r.Float64() * worldWidth
		tz = r.Float64() * worldHeight
	}

	agent := Agent{
		ID:        id,
		X:         tx,
		Z:         tz,
		Width:     1.0,
		Height:    1.0,
		VX:        math.Cos(angle),
		VZ:        math.Sin(angle),
		baseSpeed: r.Float64()*0.02 + 0.01,
		changeDirIn: r.Intn(200) + 50,
		rng:       r,
		stuck: stuckState{
			threshold: 100,
			lastX:    tx,
			lastZ:    tz,
		},
		log: wanderingLog{
			events:         make([]WanderingEvent, 0),
			lastWanderTime: time.Now(),
		},
	}
	agent.Wandering = agent.SetWanderingTarget(q)
	return agent
}
