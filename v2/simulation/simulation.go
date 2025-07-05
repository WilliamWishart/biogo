package simulation

import (
	"biogo/v2/grid"
	"fmt"
	"math"
	"math/rand"
)

type Simulation struct {
	Grid             *grid.Grid
	Population       *Population
	Tick             int
	Generation       int // Might be useless?
	GeneticDiversity float32
	Challenge        ChallengeType
}

func New() *Simulation {
	sim := Simulation{
		Challenge: Params.Challenge,
	}
	sim.InitializeGrid()
	sim.InitializeFirstGeneration()
	return &sim
}

func (s *Simulation) InitializeGrid() {
	s.Grid = grid.NewGrid(Params.GridWidth, Params.GridHeight, 0)
}

func (s *Simulation) InitializeFirstGeneration() {
	pop := NewPopulation()
	emptyLocs := s.Grid.ShuffledEmptyLocations()
	if len(emptyLocs) < Params.StartingPopulation {
		panic("Not enough empty locations for starting population")
	}
	for i := grid.RESERVED_CELL_TYPES; i < Params.StartingPopulation+grid.RESERVED_CELL_TYPES; i++ {
		loc := emptyLocs[i-grid.RESERVED_CELL_TYPES]
		pop.Creatures[i-grid.RESERVED_CELL_TYPES] = NewCreature(i, loc, MakeRandomGenome())
		s.Grid.Set(loc, i)
	}
	s.Population = pop
}

func (s *Simulation) Update() {
	if s.Tick < Params.MaxAge {
		s.Step()
	} else {
		s.InitializeNewGeneration()
	}
	if s.Generation >= Params.MaxGenerations {
		panic("Simulation ended")
	}
}

func (s *Simulation) InitializeNewGeneration() {
	// s.GeneticDiversity = s.Population.GeneticDiversity()
	s.Generation += 1
	s.Tick = 0
	childrenGenomes := []*Genome{}
	for _, creature := range s.Population.Creatures {
		if PassedSurvivalCriteria(creature, s) {
			newGenome := AsexualReproduction(creature.Genome)
			childrenGenomes = append(childrenGenomes, newGenome)
		}
	}

	if len(childrenGenomes) == 0 {
		panic("The creatures have gone extinct.")
	}
	survivalPercentage := float64(len(childrenGenomes)) / float64(len(s.Population.Creatures)) * 100
	fmt.Printf("Generation: %d\t%.2f%% Survived\n", s.Generation, survivalPercentage)

	children := []*Creature{}
	emptyLocs := s.Grid.ShuffledEmptyLocations()
	if len(emptyLocs) < Params.MaxPopulation {
		panic("Not enough empty locations for new generation")
	}
	for i := grid.RESERVED_CELL_TYPES; i < Params.MaxPopulation+grid.RESERVED_CELL_TYPES; i++ {
		loc := emptyLocs[i-grid.RESERVED_CELL_TYPES]
		child := NewCreature(i-grid.RESERVED_CELL_TYPES, loc, childrenGenomes[(i-grid.RESERVED_CELL_TYPES)%len(childrenGenomes)])
		children = append(children, child)
		s.Grid.Set(loc, i)
	}

	s.Population = &Population{
		Creatures:  children,
		DeathQueue: []DeathInstruction{},
		MoveQueue:  []MoveInstruction{},
	}
	s.Grid.ZeroFill()
	s.Grid.CreateWall()
}

func (s *Simulation) Step() {
	for _, creature := range s.Population.Creatures {
		if creature.Alive {
			s.StepCreature(creature)
		}
	}
	s.Population.ProcessMoveQueue(s.Grid)
	// TODO()
	// s.Population.ProcessReproductionQueue(s.Grid)
	// s.Population.ProcessDeathQueue()
	s.Tick++
}

func (s *Simulation) StepCreature(c *Creature) {
	c.Age++
	actionLevels := c.FeedForward(s.Grid, s.Population, s.Tick)
	s.ExecuteActions(c, actionLevels)
}

func (s *Simulation) Print() {
	s.Grid.Print()
	fmt.Printf("Population Size: %d", len(s.Population.Creatures))
}

func (s *Simulation) ExecuteActions(c *Creature, actionLevels []float32) {
	s.handleResponsiveness(c, actionLevels)
	s.handleOscillatorPeriod(c, actionLevels)

	responseAdjust := responseCurve(c.Responsiveness)
	moveX, moveY := s.handleMovement(c, actionLevels, responseAdjust)

	moveXSign, moveYSign := 1, 1
	if moveX < 0 {
		moveXSign = -1
	}
	if moveY < 0 {
		moveYSign = -1
	}

	moveXBool := prob2Bool(math.Abs(float64(moveX)))
	moveYBool := prob2Bool(math.Abs(float64(moveY)))
	movementOffset := grid.Dir{X: moveXBool * moveXSign, Y: moveYBool * moveYSign}
	newCoord := c.GetNextLoc(movementOffset)
	if s.Grid.IsInBounds(newCoord) && s.Grid.IsEmptyAt(newCoord) {
		s.Population.QueueForMove(c, newCoord)
	}
}

// Range in 0...1
func prob2Bool(val float64) int {
	if rand.Float64() < val {
		return 1
	} else {
		return 0
	}
}

func responseCurve(resp float32) float32 {
	k := float64(Params.ResponseCurveKFactor)
	return float32(math.Pow(float64(resp)-2.0, -2*k)) - float32(math.Pow(2.0, -2.0*k))*(1-resp)
}
