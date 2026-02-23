package server

import "github.com/google/uuid"

// AgentSnapshot is the JSON shape for one agent sent to clients.
type AgentSnapshot struct {
	ID       uuid.UUID `json:"id"`
	X        float64   `json:"x"`
	Z        float64   `json:"z"`
	Rotation float64   `json:"rotation"`
	Type     string    `json:"type"`
	NoPath   bool      `json:"noPath,omitempty"`
	Path     []struct {
		X float64 `json:"x"`
		Z float64 `json:"z"`
	} `json:"path,omitempty"`
}

// ObstacleSnapshot is the JSON shape for one obstacle sent to clients.
type ObstacleSnapshot struct {
	ID   uuid.UUID `json:"id"`
	MinX float64   `json:"minX"`
	MinZ float64   `json:"minZ"`
	MaxX float64   `json:"maxX"`
	MaxZ float64   `json:"maxZ"`
	Type string    `json:"type"`
}

// BroadcastMessage is the message sent to WebSocket clients each tick.
type BroadcastMessage struct {
	Tick      int                `json:"tick"`
	Updated   []AgentSnapshot    `json:"updated"`
	Obstacles []ObstacleSnapshot `json:"obstacles"`
}
