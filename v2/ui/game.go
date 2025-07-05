package ui

import (
	"biogo/v2/simulation"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// BlockSize determines the size of each grid block in pixels.
const BlockSize = 2

// StatLineXOffset is the horizontal offset for stat lines from the right edge.
const StatLineXOffset = 200

// StatLineYSpacing is the vertical spacing between stat lines.
const StatLineYSpacing = 20

// StatLineYBase is the base Y offset for stat lines.
const StatLineYBase = 3

// CenterLineWidth is the width of the center line in grid units.
const CenterLineWidth = 5

type Game struct {
	Simulation *simulation.Simulation
	Grid       *Grid

	lastGeneration int
	fontFace       font.Face // Add fontFace to reuse
}

// NewGame creates a new Game instance and initializes the font face.
// Returns an error if font parsing or face creation fails.
func NewGame(sim *simulation.Simulation) (*Game, error) {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}
	g := Game{
		Simulation: sim,
		Grid:       NewGrid(0, 0, BlockSize),
		fontFace:   face,
	}
	for _, creature := range g.Simulation.Population.Creatures {
		red, green, blue, alpha := creature.Genome.ToColor()
		c := color.RGBA{
			R: red,
			G: green,
			B: blue,
			A: alpha,
		}
		img := g.Grid.AddBlob(BlockSize, c)
		img.Translate(float64(creature.Loc.X*int(BlockSize)), float64(creature.Loc.Y*int(BlockSize)))
	}

	center := g.Simulation.Grid.SizeX() / 2
	minX := center - CenterLineWidth/2
	maxX := center + CenterLineWidth/2
	minY := g.Simulation.Grid.SizeY() / 4
	maxY := minY + g.Simulation.Grid.SizeY()/2
	g.Grid.AddLine(float64(minX*BlockSize), float64(minY*BlockSize), float64(maxX*BlockSize), float64(maxY*BlockSize))
	return &g, nil
}

func (g *Game) Update() error {
	// Simulation update should be handled outside the UI layer.
	// This method now only updates the UI representation if the simulation state has changed.
	if g.lastGeneration != g.Simulation.Generation {
		g.Grid.blobs = []*Blob{}
		for _, creature := range g.Simulation.Population.Creatures {
			red, green, blue, alpha := creature.Genome.ToColor()
			c := color.RGBA{
				R: red,
				G: green,
				B: blue,
				A: alpha,
			}
			img := g.Grid.AddBlob(BlockSize, c)
			img.Translate(float64(creature.Loc.X*int(BlockSize)), float64(creature.Loc.Y*int(BlockSize)))
		}
		g.lastGeneration = g.Simulation.Generation
	}
	for i, creature := range g.Simulation.Population.Creatures {
		img := g.Grid.blobs[i]
		img.Move(float64(creature.Loc.X*int(BlockSize)), float64(creature.Loc.Y*int(BlockSize)))
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 15, 255})
	g.Grid.DrawGrid(screen)
	g.AddStatLine(screen, "Population", len(g.Simulation.Population.Creatures), 1)
	g.AddStatLine(screen, "Generation", g.Simulation.Generation, 2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) AddStatLine(img *ebiten.Image, description string, statLine int, count int) {
	text.Draw(
		img,
		fmt.Sprintf("%s: %d", description, statLine),
		g.fontFace,
		g.Simulation.Grid.SizeX()*BlockSize-StatLineXOffset,
		StatLineYSpacing*count+StatLineYBase,
		color.White,
	)
}
