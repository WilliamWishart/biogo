package ui

import (
	"biogo/v2/simulation"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewGameInitializesGridAndSimulation(t *testing.T) {
	sim := simulation.New()
	game := NewGame(sim)
	if game.Grid == nil {
		t.Fatal("Grid should not be nil after NewGame")
	}
	if game.Simulation != sim {
		t.Fatal("Game should reference the provided simulation")
	}
}

func TestGameLayoutReturnsInput(t *testing.T) {
	game := NewGame(simulation.New())
	w, h := 800, 600
	sw, sh := game.Layout(w, h)
	if sw != w || sh != h {
		t.Errorf("Layout should return input dimensions, got (%d,%d), want (%d,%d)", sw, sh, w, h)
	}
}

func TestGameDrawDoesNotPanic(t *testing.T) {
	game := NewGame(simulation.New())
	screen := ebiten.NewImage(800, 600)
	// Should not panic
	game.Draw(screen)
}

func TestAddStatLineDoesNotPanic(t *testing.T) {
	game := NewGame(simulation.New())
	img := ebiten.NewImage(800, 600)
	// Should not panic
	game.AddStatLine(img, "TestStat", 42, 1)
}
