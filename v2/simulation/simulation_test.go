package simulation

import (
	"testing"
)

func TestSimulation_InitializeGrid(t *testing.T) {
	sim := New()
	sim.Grid = nil
	sim.InitializeGrid()
	if sim.Grid == nil {
		t.Fatal("InitializeGrid should set Grid to a new grid instance")
	}
}

func TestSimulation_InitializeFirstGeneration(t *testing.T) {
	sim := New()
	sim.Population = nil
	sim.InitializeFirstGeneration()
	if sim.Population == nil {
		t.Fatal("InitializeFirstGeneration should set Population")
	}
	if len(sim.Population.Creatures) == 0 {
		t.Fatal("Population should have creatures after InitializeFirstGeneration")
	}
}

func TestSimulation_Update_PanicsOnMaxGenerations(t *testing.T) {
	sim := New()
	sim.Generation = Params.MaxGenerations
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when Generation >= MaxGenerations, but did not panic")
		}
	}()
	sim.Update()
}

func TestSimulation_Print(t *testing.T) {
	sim := New()
	// Should not panic or error
	sim.Print()
}

func TestSimulation_StepCreature(t *testing.T) {
	sim := New()
	c := sim.Population.Creatures[0]
	ageBefore := c.Age
	sim.StepCreature(c)
	if c.Age != ageBefore+1 {
		t.Errorf("StepCreature should increment Age by 1, got %d, want %d", c.Age, ageBefore+1)
	}
}
