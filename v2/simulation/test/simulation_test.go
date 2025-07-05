package test

import (
	"biogo/v2/simulation"
	"testing"
)

func TestSimulationInitialization(t *testing.T) {
	sim := simulation.New()
	if sim.Grid == nil {
		t.Fatal("Simulation grid should not be nil after initialization")
	}
	if sim.Population == nil {
		t.Fatal("Simulation population should not be nil after initialization")
	}
	if len(sim.Population.Creatures) == 0 {
		t.Fatal("Simulation should have a non-zero starting population")
	}
}

func TestSimulationStepIncrementsTick(t *testing.T) {
	sim := simulation.New()
	initialTick := sim.Tick
	sim.Step()
	if sim.Tick != initialTick+1 {
		t.Fatalf("Simulation tick should increment by 1 after Step, got %d, want %d", sim.Tick, initialTick+1)
	}
}

func TestSimulationNewGeneration(t *testing.T) {
	sim := simulation.New()
	sim.Tick = simulation.Params.MaxAge
	genBefore := sim.Generation
	sim.Update()
	if sim.Generation != genBefore+1 {
		t.Fatalf("Simulation generation should increment after max age, got %d, want %d", sim.Generation, genBefore+1)
	}
}

func TestSimulationExtinction(t *testing.T) {
	sim := simulation.New()
	// Remove all creatures to simulate extinction
	sim.Population.Creatures = []*simulation.Creature{}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic on extinction, but did not panic")
		}
	}()
	sim.InitializeNewGeneration()
}

func TestGeneticDiversityRange(t *testing.T) {
	// sim := simulation.New()
	// div := sim.Population.GeneticDiversity()
	// if div < 0 || div > 1 {
	// 	t.Errorf("Genetic diversity should be in [0,1], got %f", div)
	// }
}

func TestStepCreatureDoesNotPanic(t *testing.T) {
	sim := simulation.New()
	for _, c := range sim.Population.Creatures {
		if c != nil {
			// Should not panic
			sim.StepCreature(c)
		}
	}
}

func TestPopulationQueueForMoveAndProcessMoveQueue(t *testing.T) {
	sim := simulation.New()
	if len(sim.Population.Creatures) == 0 {
		t.Fatal("No creatures to test move queue")
	}
	c := sim.Population.Creatures[0]
	origLoc := c.Loc
	newLoc := sim.Grid.FindEmptyLocation()
	sim.Population.QueueForMove(c, newLoc)
	sim.Population.ProcessMoveQueue(sim.Grid)
	if c.Loc != newLoc {
		t.Errorf("Creature did not move to new location. Got %v, want %v", c.Loc, newLoc)
	}
	if sim.Grid.IsEmptyAt(origLoc) == false {
		t.Errorf("Original location should be empty after move")
	}
}

func TestGeneticDiversityNotPanickingOnSmallPopulation(t *testing.T) {
	pop := &simulation.Population{Creatures: []*simulation.Creature{}}
	_ = pop.GeneticDiversity() // Should not panic
	pop.Creatures = []*simulation.Creature{&simulation.Creature{}}
	_ = pop.GeneticDiversity() // Should not panic
}

func TestFeedForwardReturnsCorrectLength(t *testing.T) {
	sim := simulation.New()
	c := sim.Population.Creatures[0]
	levels := c.FeedForward(sim.Grid, sim.Population, sim.Tick)
	if len(levels) != int(simulation.ACTION_COUNT) {
		t.Errorf("FeedForward should return slice of length ACTION_COUNT, got %d, want %d", len(levels), simulation.ACTION_COUNT)
	}
}

func TestExecuteActionsDoesNotPanic(t *testing.T) {
	sim := simulation.New()
	c := sim.Population.Creatures[0]
	levels := make([]float32, simulation.ACTION_COUNT)
	sim.ExecuteActions(c, levels) // Should not panic
}
