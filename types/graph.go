package types

import (
	"fmt"
	"os"
)

type graph struct {
	Start    *Vertex
	End      *Vertex
	Vertices []*Vertex
	Edges    map[string][]string
	Paths    [][]*Vertex
	Turns    int
}

func (graph *graph) FindVertex(name string) *Vertex {
	var found *Vertex
	for _, v := range graph.Vertices {
		if v.Name == name {
			found = v
			break
		}
	}
	if found == nil {
		FaultyData("vertex not found: " + name)
	}
	return found
}

func (graph *graph) Check() {
	if len(graph.Vertices) == 0 {
		FaultyData("no rooms specified")
	}
	if numberOfEdges() == 0 {
		FaultyData("no edges specified")
	}
	if graph.Start == nil {
		FaultyData("no start room specified")
	}
	if graph.End == nil {
		FaultyData("no end room specified")
	}
}

func numberOfEdges() int {
	var count int
	for _, edges := range Graph.Edges {
		count += len(edges)
	}
	return count
}

func (g *graph) SortDiagonally() {
	for i := 0; i < len(g.Vertices)-1; i++ {
		for j := i + 1; j < len(g.Vertices); j++ {
			pos1 := g.Vertices[i].Position
			pos2 := g.Vertices[j].Position
			if pos1.Y > pos2.Y || pos1.Y == pos2.Y && pos1.X > pos2.X {
				g.Vertices[i], g.Vertices[j] = g.Vertices[j], g.Vertices[i]
			}
		}
	}
}

func (g *graph) SortByDensest() {
	// Sort by number of edges; with more edges in front
	for i := 0; i < len(g.Vertices)-1; i++ {
		for j := i + 1; j < len(g.Vertices); j++ {
			if len(g.Vertices[i].Edges) > len(g.Vertices[j].Edges) {
				g.Vertices[i], g.Vertices[j] = g.Vertices[j], g.Vertices[i]
			}
		}
	}
}

var Graph = &graph{
	Edges:    make(map[string][]string),
	Paths:    [][]*Vertex{},
	Vertices: []*Vertex{},
}

func FaultyData(msg string) {
	fmt.Printf("ERROR: invalid data format, %s\n", msg)
	os.Exit(0)
}
