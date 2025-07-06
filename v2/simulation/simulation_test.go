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

func TestSimulation_ExecuteActions_Movement(t *testing.T) {
	sim := New()
	c := sim.Population.Creatures[0]
	c.Alive = true

	// Place creature away from the east edge
	c.Loc.X = 1
	c.Loc.Y = 1
	sim.Grid.Set(c.Loc, c.Id)
	// Ensure east cell is empty
	east := c.Loc
	east.X++
	sim.Grid.Set(east, 0)

	oldLoc := c.Loc
	actionLevels := make([]float32, 16)
	actionLevels[MOVE_EAST] = 1.0

	moved := false
	for i := 0; i < 100; i++ { // Try up to 10 times to account for probabilistic movement
		sim.ExecuteActions(c, actionLevels)
		sim.Population.ProcessMoveQueue(sim.Grid)
		if c.Loc != oldLoc {
			moved = true
			break
		}
	}
	if !moved {
		t.Errorf("Expected creature to move east, but location did not change after multiple attempts")
	}
}

func TestSimulation_ExecuteActions_Responsiveness(t *testing.T) {
	sim := New()
	c := sim.Population.Creatures[0]
	c.Alive = true
	oldResp := c.Responsiveness
	actionLevels := make([]float32, 16)
	actionLevels[SET_RESPONSIVENESS] = 1.0

	sim.ExecuteActions(c, actionLevels)

	if c.Responsiveness == oldResp {
		t.Errorf("Expected responsiveness to change, but it did not")
	}
}

func TestSimulation_ExecuteActions_OscillatorPeriod(t *testing.T) {
	sim := New()
	c := sim.Population.Creatures[0]
	c.Alive = true
	oldClock := c.Clock
	actionLevels := make([]float32, 16)
	actionLevels[SET_OSCILLATOR_PERIOD] = 1.0

	sim.ExecuteActions(c, actionLevels)

	if c.Clock == oldClock {
		t.Errorf("Expected oscillator period to change, but it did not")
	}
}
