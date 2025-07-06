package simulation

import (
	"biogo/v2/grid"
	"fmt"
	"log"
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
	log.Println("Initializing grid...")
	s.Grid = grid.NewGrid(Params.GridWidth, Params.GridHeight, 0)
}

func (s *Simulation) InitializeFirstGeneration() {
	log.Println("Initializing first generation...")
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
	log.Printf("Simulation update: Generation %d, Tick %d", s.Generation, s.Tick)
	if s.Tick < Params.MaxAge {
		s.Step()
	} else {
		log.Println("Max age reached, initializing new generation...")
		s.InitializeNewGeneration()
	}
	if s.Generation >= Params.MaxGenerations {
		log.Println("Simulation ended: reached max generations.")
		panic("Simulation ended")
	}
}

func (s *Simulation) InitializeNewGeneration() {
	log.Printf("Initializing new generation: %d", s.Generation+1)
	s.GeneticDiversity = s.Population.GeneticDiversity()
	s.Generation += 1
	s.Tick = 0
	log.Printf("Population before survival: %d", len(s.Population.Creatures))
	childrenGenomes := []*Genome{}
	for _, creature := range s.Population.Creatures {
		if PassedSurvivalCriteria(creature, s) {
			newGenome := AsexualReproduction(creature.Genome)
			childrenGenomes = append(childrenGenomes, newGenome)
		}
	}
	log.Printf("Children genomes after survival: %d", len(childrenGenomes))

	if len(childrenGenomes) == 0 {
		log.Println("All creatures have gone extinct.")
		panic("The creatures have gone extinct.")
	}
	survivalPercentage := float64(len(childrenGenomes)) / float64(len(s.Population.Creatures)) * 100
	log.Printf("Generation: %d\t%.2f%% Survived", s.Generation, survivalPercentage)
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
		if (i-grid.RESERVED_CELL_TYPES)%100 == 0 { // Log every 100th child for progress
			log.Printf("Created child %d at %v", i-grid.RESERVED_CELL_TYPES, loc)
		}
	}
	log.Printf("Total children created: %d", len(children))

	s.Population = &Population{
		Creatures:  children,
		DeathQueue: []DeathInstruction{},
		MoveQueue:  []MoveInstruction{},
	}
	log.Println("Zero-filling grid...")
	s.Grid.ZeroFill()
	log.Println("Creating wall...")
	s.Grid.CreateWall()
	log.Println("New generation initialization complete.")
}

func (s *Simulation) Step() {
	log.Printf("Simulation step: Generation %d, Tick %d", s.Generation, s.Tick)
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

// StepCreature advances the state of a single creature within the simulation by one step.
// It logs the creature's ID and age, increments the creature's age, processes its neural
// network to determine action levels, and then executes the resulting actions.
//
// Parameters:
//
//	c - Pointer to the Creature to be stepped.
func (s *Simulation) StepCreature(c *Creature) {
	log.Printf("Stepping creature ID %d, Age %d", c.Id, c.Age)
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
