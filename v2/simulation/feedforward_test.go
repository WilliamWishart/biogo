package simulation

import (
	"testing"
)

func BenchmarkFeedForward(b *testing.B) {
	sim := New()
	c := sim.Population.Creatures[0]
	grid := sim.Grid
	pop := sim.Population
	tick := sim.Tick

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.FeedForward(grid, pop, tick)
	}
}
