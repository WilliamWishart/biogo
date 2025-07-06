// actions.go: Contains logic for handling all simulation actions (movement, responsiveness, oscillator period, etc.) for creatures.

package simulation

// Action logic for Simulation.ExecuteActions split out for clarity and maintainability.
import (
	"biogo/v2/grid"
	"biogo/v2/utils"
	"math"
)

func (s *Simulation) handleResponsiveness(c *Creature, actionLevels []float32) {
	if IsActionEnabled(SET_RESPONSIVENESS) {
		level := actionLevels[SET_RESPONSIVENESS]
		level = (float32(math.Tanh(float64(level/float32(utils.ClampByteAsFloat32(0, 1, c.Genome.Responsiveness))))) + 1) / 2
		c.Responsiveness = level
	}
}

func (s *Simulation) handleOscillatorPeriod(c *Creature, actionLevels []float32) {
	if IsActionEnabled(SET_OSCILLATOR_PERIOD) {
		periodf := actionLevels[SET_OSCILLATOR_PERIOD]
		newPeriodf := float32(math.Tanh(float64(periodf)+1) / 2)
		newPeriod := 1 + int(1.5+math.Exp(7*float64(newPeriodf)))
		if newPeriod >= 2 && newPeriod <= math.MaxUint8 {
			c.Clock = newPeriod
		}
	}
}

func (s *Simulation) handleMovement(c *Creature, actionLevels []float32, responseAdjust float32) (float32, float32) {
	moveX, moveY := float32(0), float32(0)
	if IsActionEnabled(MOVE_X) {
		moveX = actionLevels[MOVE_X]
	}
	if IsActionEnabled(MOVE_Y) {
		moveX = actionLevels[MOVE_Y]
	}
	if IsActionEnabled(MOVE_EAST) {
		moveX += actionLevels[MOVE_EAST]
	}
	if IsActionEnabled(MOVE_WEST) {
		moveX -= actionLevels[MOVE_WEST]
	}
	if IsActionEnabled(MOVE_NORTH) {
		moveY += actionLevels[MOVE_NORTH]
	}
	if IsActionEnabled(MOVE_SOUTH) {
		moveY += actionLevels[MOVE_SOUTH]
	}
	if IsActionEnabled(MOVE_FWD) {
		level := actionLevels[MOVE_FWD]
		moveX += float32(c.LastMoveDir.X) * level
		moveY += float32(c.LastMoveDir.Y) * level
	}
	if IsActionEnabled(MOVE_LEFT) {
		level := actionLevels[MOVE_LEFT]
		offset := c.LastMoveDir.Rotate90CCW()
		moveX += float32(offset.X) * level
		moveY += float32(offset.Y) * level
	}
	if IsActionEnabled(MOVE_RIGHT) {
		level := actionLevels[MOVE_RIGHT]
		offset := c.LastMoveDir.Rotate90CW()
		moveX += float32(offset.X) * level
		moveY += float32(offset.Y) * level
	}
	if IsActionEnabled(MOVE_RL) {
		level := actionLevels[MOVE_RL]
		offset := grid.CENTER
		if level < 0 {
			offset = c.LastMoveDir.Rotate90CCW()
		} else if level > 0 {
			offset = c.LastMoveDir.Rotate90CW()
		}
		moveX += float32(offset.X) * level
		moveY += float32(offset.Y) * level
	}
	if IsActionEnabled(MOVE_RANDOM) {
		level := actionLevels[MOVE_RANDOM]
		offset := grid.RandomDir()
		moveX += float32(offset.X) * level
		moveY += float32(offset.Y) * level
	}
	moveX = float32(math.Tanh(float64(moveX))) * responseAdjust
	moveY = float32(math.Tanh(float64(moveY))) * responseAdjust
	return moveX, moveY
}
