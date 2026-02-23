package agents

import (
	"fmt"
	"time"
	worldQuery "veatla/simulator/src/world-query"
)

func (agent *Agent) logWanderingEvent(q worldQuery.WorldQuery) {
	now := time.Now()
	duration := now.Sub(agent.log.lastWanderTime)

	event := WanderingEvent{
		Timestamp: now,
		X:         agent.X,
		Z:         agent.Z,
		TargetX:   agent.Wandering.X,
		TargetZ:   agent.Wandering.Z,
		Duration:  duration,
	}

	agent.log.events = append(agent.log.events, event)
	agent.log.lastWanderTime = now
}

func (agent *Agent) GetWanderingEvents() []WanderingEvent {
	return agent.log.events
}

func (agent *Agent) ClearWanderingEvents() {
	agent.log.events = make([]WanderingEvent, 0)
}

func (agent *Agent) PrintWanderingLog() {
	fmt.Printf("\n=== Wandering Log for Agent %s ===\n", agent.ID.String()[:8])
	fmt.Printf("Total wandering events: %d\n", len(agent.log.events))
	for i, event := range agent.log.events {
		fmt.Printf("  Event %d: [%s] From (%.2f, %.2f) to (%.2f, %.2f) - Duration: %v\n",
			i+1, event.Timestamp.Format("15:04:05"), event.X, event.Z, event.TargetX, event.TargetZ, event.Duration)
	}
}
