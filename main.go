package main

import (
	"biogo/v2/simulation"
	"biogo/v2/ui"
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	enableProfile := false
	// Check environment variable
	if os.Getenv("BIOGO_PROFILE") == "1" {
		enableProfile = true
	}
	// Check command-line flag
	profileFlag := flag.Bool("profile", false, "Enable CPU and memory profiling")
	flag.Parse()
	if *profileFlag {
		enableProfile = true
	}

	rand.Seed(time.Now().UnixNano())

	var f, mf *os.File
	if enableProfile {
		// Start CPU profiling
		var err error
		f, err = os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()

		// Start memory profiling
		mf, err = os.Create("mem.prof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer func() {
			runtime.GC() // get up-to-date statistics
			pprof.WriteHeapProfile(mf)
			mf.Close()
		}()
	}

	sim := simulation.New()
	// for i := 0; i < 2*simulation.Params.MaxAge; i++ {
	// 	start := time.Now()
	// 	sim.Update()
	// 	end := time.Now()
	// 	if sim.Tick != 0 && sim.Tick%simulation.Params.MaxAge == 0 {
	// 		fmt.Printf("\tStep took : %s\n", end.Sub(start))
	// 	}
	// }

	game := ui.NewGame(sim)

	ebiten.SetWindowSize(simulation.Params.GridWidth*2, simulation.Params.GridHeight*2)
	ebiten.SetWindowTitle("Genetic Simulation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal()
	}
}
