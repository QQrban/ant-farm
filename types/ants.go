package types

import (
	"fmt"
	"strconv"
)

type Ant struct {
	Name    string
	Current int
	Queue   int
}

type ants struct {
	Number int
	All    map[string]*Ant
	Queues [][]*Ant
}

var AllMoves [][]string

func (ants *ants) Distribute() {
	// Distribute ants into queues
	// There are as many queues as there are paths
	// Each path has its own queue
	ants.Queues = make([][]*Ant, len(Graph.Paths))
	for i := 0; i < ants.Number; i++ {
		// Ants are named 1, 2, 3, ...
		name := strconv.Itoa(i + 1)
		ant := ants.All[name]
		for j := range ants.Queues {
			pathLen := len(Graph.Paths[j]) - 1
			queueLen := len(ants.Queues[j])
			// If the queue is full, it is no longer available
			// (Queue is full when the number of ants in the queue
			// combined with the length of the path
			// is greater than the number of minimal turns)
			if pathLen+queueLen <= Graph.Turns {
				// If there is room in queue, add ant to it
				ants.Queues[j] = append(ants.Queues[j], ant)
				ant.Queue = j
				// Ants are initially in the start vertex
				ant.Current = 0
				break
			}
		}
	}
}

func (ants *ants) Step(webVisualisation bool) {
	var movesOnStep []string
NEXT_QUEUE:
	// In turn for each queue...
	for _, queue := range ants.Queues {
		if len(queue) > 0 {
			for _, ant := range queue {
				// If ant is on start vertex, move it to the next vertex...
				if ant.Current == 0 {
					ant.Current++
					if !Visual.Animate {
						ants.Print(ant, webVisualisation)
						concat := fmt.Sprintf("%s-%s", ant.Name, Graph.Paths[ant.Queue][ant.Current].Name)
						movesOnStep = append(movesOnStep, concat)
					}
					// ...and continue with next queue
					continue NEXT_QUEUE
				} else if ant.Current < len(Graph.Paths[ant.Queue])-1 {
					// If ant is not on end vertex, move it to the next vertex
					ant.Current++
					if !Visual.Animate {
						ants.Print(ant, webVisualisation)
						concat := fmt.Sprintf("%s-%s", ant.Name, Graph.Paths[ant.Queue][ant.Current].Name)
						movesOnStep = append(movesOnStep, concat)
					}
				}
			}
		}
	}
	AllMoves = append(AllMoves, movesOnStep)
}

func (ants *ants) Move(webVisualisation bool) {
	for i := 0; i < Graph.Turns; i++ {
		ants.Step(webVisualisation)
		if !webVisualisation {
			if !Visual.Animate {
				fmt.Println()
			}
		}
	}
}

func (ants *ants) Print(ant *Ant, webVisualization bool) {
	if !webVisualization {
		fmt.Printf("L%s-%s ", ant.Name, Graph.Paths[ant.Queue][ant.Current].Name)
	}

}

var Ants = &ants{}
