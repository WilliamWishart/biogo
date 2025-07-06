package simulation

import (
	"biogo/v2/grid"
	"fmt"
	"log"
	"math"
)

// The "Brains"

// FeedForward performs a single feedforward pass through the creature's neural network.
// It processes the network's edges, propagating sensor inputs and hidden neuron outputs
// to compute the activation levels for each possible action. The function updates hidden
// neuron outputs as needed and accumulates weighted inputs for both hidden neurons and
// action outputs. It returns a slice of float32 values representing the activation levels
// for each action at the current simulation step.
//
// Parameters:
//
//	g    - Pointer to the simulation grid.
//	p    - Pointer to the population.
//	step - The current simulation step.
//
// Returns:
//
//	[]float32 - Slice containing the activation levels for each action.
func (c *Creature) FeedForward(g *grid.Grid, p *Population, step int) []float32 {
	log.Printf("FeedForward: Creature Id %d, Step %d", c.Id, step)
	actionLevels := make([]float32, ACTION_COUNT)
	neuronAccumulators := map[byte]float32{}
	neuronOutputsEvaluated := false

	for _, gene := range c.Nnet.Edges {

		// First we evaluate the outputs to ACTIONS
		if gene.SinkType == ACTION && !neuronOutputsEvaluated {
			log.Printf("Evaluating hidden neuron outputs for Creature Id %d", c.Id)
			for key, neuron := range c.Nnet.HiddenNeurons {
				if neuron.Driven {
					neuron.Output = float32(math.Tanh(float64(neuronAccumulators[key])))
				}
			}
			neuronOutputsEvaluated = true
		}

		var inputVal float32
		if gene.SourceType == SENSOR {
			log.Printf("Getting sensor value: SourceID %d, Creature Id %d", gene.SourceID, c.Id)
			inputVal = c.GetSensor(gene.SourceID, g, p, step)
		} else {
			if _, ok := c.Nnet.HiddenNeurons[gene.SourceID]; !ok {
				log.Printf("Hidden neuron not found: SourceID %d, Type %d, Creature Id %d", gene.SourceID, gene.SourceType, c.Id)
				fmt.Printf("\n\nNot okay, trying to see %d of type %d, %s", gene.SourceID, gene.SourceType, c.Nnet.String())
				for _, gene := range c.Nnet.Edges {
					fmt.Printf("\n%s", gene.PrettyString())
				}
				fmt.Printf("\nC.Nnet.HiddenNeurons: %v\n", c.Nnet.HiddenNeurons)
			}
			inputVal = c.Nnet.HiddenNeurons[gene.SourceID].Output
		}

		if gene.SinkType == ACTION {
			log.Printf("Accumulating action: SinkID %d, Creature Id %d", gene.SinkID, c.Id)
			actionLevels[gene.SinkID] += inputVal * gene.WeightAsFloat32()
		} else {
			log.Printf("Accumulating hidden neuron: SinkID %d, Creature Id %d", gene.SinkID, c.Id)
			neuronAccumulators[gene.SinkID] += inputVal * gene.WeightAsFloat32()
		}
	}
	log.Printf("FeedForward complete: Creature Id %d, ActionLevels: %v", c.Id, actionLevels)
	return actionLevels
}
