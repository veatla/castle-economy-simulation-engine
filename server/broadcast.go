package server

import (
	"math"

	"veatla/simulator/src/agents"
	"veatla/simulator/src/constructions"
)

const broadcastBatchSize = 1000

// BroadcastWorld sends updated agents and obstacles to connected WebSocket clients in batches.
func BroadcastWorld(tick int, updated []agents.Agent, obstacles []constructions.Obstacle) {
	obsSnap := make([]ObstacleSnapshot, 0, len(obstacles))
	for _, o := range obstacles {
		obsSnap = append(obsSnap, ObstacleSnapshot{
			ID:   o.ID,
			MinX: o.MinX,
			MinZ: o.MinZ,
			MaxX: o.MaxX,
			MaxZ: o.MaxZ,
			Type: "obstacle",
		})
	}

	total := len(updated)
	if total == 0 {
		msg := BroadcastMessage{Tick: tick, Updated: []AgentSnapshot{}, Obstacles: obsSnap}
		hub.broadcast(msg)
		return
	}

	for start := 0; start < total; start += broadcastBatchSize {
		end := start + broadcastBatchSize
		if end > total {
			end = total
		}
		chunk := updated[start:end]
		snap := make([]AgentSnapshot, 0, len(chunk))
		for _, a := range chunk {
			as := AgentSnapshot{
				ID:       a.ID,
				X:        a.X,
				Z:        a.Z,
				Type:     "agent",
				Rotation: math.Atan2(a.VZ, a.VX) + math.Pi/2,
				NoPath:   a.NoPath,
			}
			if len(a.GetPath()) > 0 {
				as.Path = make([]struct {
					X float64 `json:"x"`
					Z float64 `json:"z"`
				}, len(a.GetPath()))
				for i, p := range a.GetPath() {
					as.Path[i].X = p.X
					as.Path[i].Z = p.Z
				}
			}
			snap = append(snap, as)
		}
		msg := BroadcastMessage{Tick: tick, Updated: snap, Obstacles: obsSnap}
		hub.broadcast(msg)
	}
}
