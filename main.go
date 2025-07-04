package main

import (
	"biogo/v2/simulation"
	"biogo/v2/ui"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.Println("Seeding random number generator...")
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	_ = r // Use r for random numbers instead of global rand

	log.Println("Initializing simulation...")
	sim := simulation.New()

	log.Println("Running simulation steps...")
	for i := 0; i < 50*simulation.Params.MaxAge; i++ {
		start := time.Now()
		sim.Update()
		end := time.Now()
		if sim.Tick != 0 && sim.Tick%simulation.Params.MaxAge == 0 {
			fmt.Printf("\tStep took : %s\n", end.Sub(start))
		}
	}

	log.Println("Creating game UI...")
	game, err := ui.NewGame(sim)
	if err != nil {
		log.Fatalf("failed to create game: %v", err)
	}

	log.Println("Configuring window...")
	ebiten.SetWindowSize(simulation.Params.GridWidth*2, simulation.Params.GridHeight*2)
	ebiten.SetWindowTitle("Genetic Simulation")

	log.Println("Starting game loop...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
