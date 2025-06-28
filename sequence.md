# Genetic Simulation - Diagrams and Flow

This document provides PlantUML diagrams to illustrate the structure and operation of the main components in the genetic simulation program.

---

## 1. Component/Class Diagram

```plantuml
@startuml
package main {
  class Main {
    +main()
  }
}

package simulation {
  class Simulation {
    +New()
    +Update()
    +Tick : int
    +Params : Params
  }
  class Params {
    +MaxAge : int
    +GridWidth : int
    +GridHeight : int
  }
}

package ui {
  class Game {
    +NewGame(sim: Simulation)
  }
}

package ebiten {
  class Ebiten {
    +SetWindowSize(w, h)
    +SetWindowTitle(title)
    +RunGame(game)
  }
}

Main --> Simulation : uses
Main --> Game : uses
Game --> Simulation : references
Main --> Ebiten : uses
Simulation --> Params : uses
@enduml
```

---

## 2. Sequence Diagram

```plantuml
@startuml
actor User
User -> Main: start program
activate Main
Main -> Simulation: New()
Main -> Simulation: Update() [loop 50*MaxAge]
Main -> UI.Game: NewGame(sim)
Main -> Ebiten: SetWindowSize()
Main -> Ebiten: SetWindowTitle()
Main -> Ebiten: RunGame(game)
Ebiten -> Game: game loop
Game -> Simulation: access state/render
@enduml
```

---

## 3. Process Flow Chart

```plantuml
@startuml
start
:Seed random;
:Create Simulation;
repeat
  :Update Simulation;
  if (Tick % MaxAge == 0?) then (yes)
    :Print step duration;
  endif
repeat while (i < 50*MaxAge?)
:Create Game with Simulation;
:Set window size and title;
:Run game loop with Ebiten;
stop
@enduml
```

---

**Summary:**
- The program seeds randomness, creates a simulation, runs updates for a set number of steps, then creates a UI game and starts the Ebiten game loop for visualization.
- The diagrams above illustrate the relationships, runtime sequence, and process flow.
