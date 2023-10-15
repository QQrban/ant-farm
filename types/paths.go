package types

import (
	"fmt"
	"math"
)

type Path struct {
	Positions []Pos
	Marks     []rune
}

func (path *Path) Add(pos Pos, mark rune) {
	path.Positions = append(path.Positions, pos)
	path.Marks = append(path.Marks, mark)
}

type paths struct {
	All [][]*Vertex
}

func (paths *paths) String() string {
	str := ""
	for _, path := range paths.All {
		for _, v := range path {
			str += v.Name + " "
		}
		str += "\n"
	}
	return str
}

func (paths *paths) Find() {
	// Find all paths
	Explorer.Explore()
	paths.Sort()
	paths.Disjoin()
}

func (paths *paths) Sort() {
	for i := 0; i < len(paths.All)-1; i++ {
		for j := i; j < len(paths.All); j++ {
			if len(paths.All[i]) > len(paths.All[j]) {
				paths.All[i], paths.All[j] = paths.All[j], paths.All[i]
			}
		}
	}
}

func (paths *paths) Disjoin() {
	//fmt.Println(paths)
	var candidates = make(map[int][]Combination)
	for i, path := range paths.All {
		var disjoint [][]*Vertex
		disjoint = append(disjoint, path)
		if i == 0 {
			candidates[1] = []Combination{Combination{len: paths.length(disjoint), last: 0, paths: disjoint}}
		}
		paths.AddCandidates(i, disjoint, &candidates)
	}
	//printCombination(candidates)
	paths.SelectBest(candidates)
}

func (paths *paths) AddCandidates(i int, disjoint [][]*Vertex, candidates *map[int][]Combination) {
	for j := i + 1; j < len(paths.All); j++ {
		if paths.areDisjoint(disjoint, paths.All[j]) {
			disjoint1 := append(disjoint, paths.All[j])
			if _, exists := (*candidates)[len(disjoint1)]; !exists {
				(*candidates)[len(disjoint1)] = []Combination{Combination{len: paths.length(disjoint1), last: j, paths: disjoint1}}
			} else {
				(*candidates)[len(disjoint1)] = append((*candidates)[len(disjoint1)], Combination{len: paths.length(disjoint1), last: j, paths: disjoint1})
			}
			paths.AddCandidates(j, disjoint1, candidates)
		}
	}
}

func printCombination(candidates map[int]Combination) {
	fmt.Println("Candidates:")
	for i, c := range candidates {
		fmt.Printf("%d:", i)
		fmt.Printf(" %d", c.len)
		for _, path := range c.paths {
			fmt.Printf(" %d(", len(path)-1)
			for _, v := range path {
				fmt.Printf("%s ", v.Name)
			}
			fmt.Printf(")")
		}
		fmt.Println()
	}
}

func (paths *paths) SelectBest(candidates map[int][]Combination) {
	Graph.Turns = math.MaxInt32
	for i, combinations := range candidates {
		for _, combination := range combinations {
			elements := Ants.Number + combination.len - i
			steps := elements / i
			if elements%i > 0 {
				steps++
			}
			if steps < Graph.Turns {
				Graph.Paths = combination.paths
				Graph.Turns = steps
			}
		}
	}
}

func (paths *paths) length(disjoint [][]*Vertex) int {
	length := 0
	for _, path := range disjoint {
		length += len(path) - 1
	}
	return length
}

func (paths *paths) areDisjoint(disjointPaths [][]*Vertex, path2 []*Vertex) bool {
	for _, path1 := range disjointPaths {
		for _, v1 := range path1[1 : len(path1)-1] {
			for _, v2 := range path2[1 : len(path2)-1] {
				if v1 == v2 {
					return false
				}
			}
		}
	}
	return true
}

var Paths = &paths{}
