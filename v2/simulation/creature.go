// creature.go: Defines the Creature struct and related methods for simulation entities (creatures), including neural net creation and movement logic.

package simulation

import (
	"biogo/v2/grid"
	"biogo/v2/utils"
	"fmt"
	"math"
)

type Creature struct {
	Id             int
	Energy         float32
	Responsiveness float32
	Age            int
	Alive          bool
	Clock          int
	Nnet           NeuralNet
	Loc            grid.Coord
	BirthLoc       grid.Coord
	LastMoveDir    grid.Dir
	Genome         *Genome

	actionLevelsBuf       []float32
	neuronAccumulatorsBuf []float32
}

func NewCreature(id int, loc grid.Coord, g *Genome) *Creature {
	c := Creature{
		Id:             id,
		Energy:         float32(g.MaxEnergy / math.MaxUint8),
		Age:            0,
		Alive:          true,
		Clock:          int(g.OscPeriod), // TODO() Maybe fix this?
		Nnet:           NeuralNet{},
		Loc:            loc,
		BirthLoc:       loc,
		Responsiveness: float32(utils.ClampByteAsFloat32(0, 1, g.Responsiveness)) / 2,
		Genome:         g,
	}
	c.CreateNeuralNet()
	return &c
}

// Takes a creature's genome and uses it to build a NeuralNetwork
func (c *Creature) CreateNeuralNet() {
	c.Nnet = *CreateNeuralNetworkFromGenome(c.Genome.Brain, c.Genome.NeuronCount)
	// Preallocate buffers for FeedForward
	c.actionLevelsBuf = make([]float32, ACTION_COUNT)
	c.neuronAccumulatorsBuf = make([]float32, len(c.Nnet.HiddenNeurons))
}

func (c Creature) String() string {
	return fmt.Sprintf("\nCREATURE| \nID: %d,\nEnergy: %f,\nResponsiveness: %f,\nAge: %d,\nAlive: %t,\nClock: %d,\nNnet: \n%s,\nLoc: %v,\nBirthLoc: %v,\nLastMoveDir%v",
		c.Id,
		c.Energy,
		c.Responsiveness,
		c.Age,
		c.Alive,
		c.Clock,
		c.Nnet.String(),
		c.Loc,
		c.BirthLoc,
		c.LastMoveDir)
}

func (c Creature) GetNextLoc(d grid.Dir) grid.Coord {
	x := c.Loc.X + d.X
	y := c.Loc.Y + d.Y
	return grid.Coord{
		X: x,
		Y: y,
	}
}
