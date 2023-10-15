package types

import (
	"fmt"
	"slices"
	"strings"
)

type Turn struct {
	InStart int
	InEnd   int
	EnRoute map[*Ant]*Vertex
}

func (t *Turn) String() string {
	return fmt.Sprintf("in start: %d, en route: %d, in end: %d", t.InStart, len(t.EnRoute), t.InEnd)
}

type turns struct {
	InStart int
	InEnd   int
	Data    []*Turn
}

func (t *turns) String() string {
	for i, turn := range t.Data {
		fmt.Printf("Turn %d: %s\n", i, turn)
	}
	return ""
}

func (t *turns) Parse(line string) {
	// Initialize turn
	turn := &Turn{
		InStart: t.InStart,
		InEnd:   t.InEnd,
		EnRoute: make(map[*Ant]*Vertex),
	}
	// Parse line
	fields := strings.Split(strings.Trim(line, " "), " ")
	for _, field := range fields {
		parts := strings.Split(field, "-")
		antName := parts[0][1:]
		ant := Ants.All[antName]
		vertexName := parts[1]
		vertex := Graph.FindVertex(vertexName)
		// Track ant's movement
		if ant.Current == 0 {
			// Just stepping out of start
			turn.InStart--
		}
		if vertex == Graph.End {
			// Reaches end
			turn.InEnd++
			turn.EnRoute[ant] = vertex
		} else {
			// Still en route
			turn.EnRoute[ant] = vertex
		}
		// Increase step counter
		ant.Current++
	}
	// Update turns data
	t.Data = append(t.Data, turn)
	t.InStart = turn.InStart
	t.InEnd = turn.InEnd
}

func (t *turns) ExtractPaths() {
	// Initialize paths
	paths := make(map[*Ant][]*Vertex)
	for _, ant := range Ants.All {
		paths[ant] = []*Vertex{}
	}
	// Extract paths
	for _, turn := range t.Data {
		for ant, vertex := range turn.EnRoute {
			paths[ant] = append(paths[ant], vertex)
		}
	}
	for ant, path := range paths {
		// Amend paths
		path = append([]*Vertex{Graph.Start}, path...)
		path = append(path, Graph.End)
		// Check if path is already in paths
		var found bool
		for i, p := range Graph.Paths {
			if slices.Equal(p, path) {
				found = true
				Ants.Queues[i] = append(Ants.Queues[i], ant)
				ant.Queue = i
				break
			}
		}
		if !found {
			Graph.Paths = append(Graph.Paths, path)
			Ants.Queues = append(Ants.Queues, []*Ant{ant})
			ant.Queue = len(Graph.Paths) - 1
		}
	}
	// Reset current position of ants to start vertex
	for _, ant := range Ants.All {
		ant.Current = 0
	}
}

var Turns = turns{
	InStart: Ants.Number,
	InEnd:   0,
	Data:    []*Turn{},
}
