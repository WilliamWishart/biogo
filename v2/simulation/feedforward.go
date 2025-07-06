package simulation

import (
	"biogo/v2/grid"
	"math"
)

// The "Brains"

func (c *Creature) FeedForward(g *grid.Grid, p *Population, step int) []float32 {
	// Zero buffers
	for i := range c.actionLevelsBuf {
		c.actionLevelsBuf[i] = 0
	}
	for i := range c.neuronAccumulatorsBuf {
		c.neuronAccumulatorsBuf[i] = 0
	}
	neuronOutputsEvaluated := false

	for _, gene := range c.Nnet.Edges {

		// First we evaluate the outputs to ACTIONS
		if gene.SinkType == ACTION && !neuronOutputsEvaluated {
			for key, neuron := range c.Nnet.HiddenNeurons {
				if neuron.Driven {
					neuron.Output = float32(math.Tanh(float64(c.neuronAccumulatorsBuf[int(key)])))
				}
			}
			neuronOutputsEvaluated = true
		}

		var inputVal float32
		if gene.SourceType == SENSOR {
			inputVal = c.GetSensor(gene.SourceID, g, p, step)
		} else {
			inputVal = c.Nnet.HiddenNeurons[gene.SourceID].Output
		}

		if gene.SinkType == ACTION {
			c.actionLevelsBuf[gene.SinkID] += inputVal * gene.WeightAsFloat32()
		} else {
			c.neuronAccumulatorsBuf[gene.SinkID] += inputVal * gene.WeightAsFloat32()
		}
	}
	return c.actionLevelsBuf
}
