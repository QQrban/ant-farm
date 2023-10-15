package types

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type visual struct {
	Lines      [][]rune
	Zoom       Pos
	MaxDim     Pos
	Animate    bool
	VIndex     map[int]map[int][]rune // Vertex positions index
	Index      map[int]map[int]rune   // y -> x -> path char: |-/\`´.:
	Auto       bool
	Sorted     map[int]map[int]string
	ColorIndex map[int]map[int]string
}

func (v *visual) SetDimensions() {
	maxDim := Pos{X: 0, Y: 0}
	for _, vertex := range Graph.Vertices {
		pos := vertex.Position
		if pos.Y > maxDim.Y {
			maxDim.Y = pos.Y
		}
		if pos.X > maxDim.X {
			maxDim.X = pos.X
		}
		vertex.Position.X *= v.Zoom.X
		vertex.Position.Y *= v.Zoom.Y
	}
	v.MaxDim = Pos{X: (maxDim.X + 1) * v.Zoom.X, Y: maxDim.Y*v.Zoom.Y + 1}
	if v.Lines == nil {
		v.Lines = make([][]rune, v.MaxDim.Y)
	}
	for i := 0; i < v.MaxDim.Y; i++ {
		v.Lines = append(v.Lines, make([]rune, v.MaxDim.X))
		v.Lines[i] = []rune(strings.Repeat(" ", v.MaxDim.X))
	}
}

func (v *visual) MakeVIndex() {
	v.VIndex = make(map[int]map[int][]rune)
	for _, vertex := range Graph.Vertices {
		pos := vertex.Position
		l := min(len(vertex.Name), v.Zoom.X)
		if _, exists := v.VIndex[pos.Y]; !exists {
			v.VIndex[pos.Y] = map[int][]rune{pos.X: []rune(vertex.Name[:l])}
		} else {
			v.VIndex[pos.Y][pos.X] = []rune(vertex.Name[:l])
		}
	}
}

func (v *visual) PutVertices() {
	for _, vertex := range Graph.Vertices {
		pos := vertex.Position
		name := v.VIndex[pos.Y][pos.X]
		l := len(name)
		lineStart := v.Lines[pos.Y][:pos.X]
		lineEnd := v.Lines[pos.Y][pos.X+l:]
		v.Lines[pos.Y] = append(lineStart, name...)
		v.Lines[pos.Y] = append(v.Lines[pos.Y], lineEnd...)
	}
}

func (v *visual) PutEdges() {
	Graph.SortDiagonally()
	done := map[*Vertex][]*Vertex{}
	for _, v1 := range Graph.Vertices {
		v1.SortEdgesCross()
		pos1 := v1.Position
		len1 := len(v.VIndex[pos1.Y][pos1.X])
		for _, v2 := range v1.Sorted {
			if _, exists := done[v2]; !exists || !slices.Contains(done[v2], v1) {
				pos2 := v2.Position
				if pos1.Y == pos2.Y { // Horizontally aligned vertices
					start := pos1.X + len1
					str := strings.Repeat("-", pos2.X-start)
					ending := v.Lines[pos1.Y][start+len(str):]
					v.Lines[pos1.Y] = append(v.Lines[pos1.Y][:start], []rune(str)...)
					v.Lines[pos1.Y] = append(v.Lines[pos1.Y], ending...)
					v1.Edges[v2] = &Path{Positions: []Pos{}, Marks: []rune(str)}
					for x := 0; x < len(str); x++ {
						if _, exists := v.Index[pos1.Y]; !exists {
							v.Index[pos1.Y] = map[int]rune{start + x: '-'}
						} else {
							v.Index[pos1.Y][start+x] = '-'
						}
						v1.Edges[v2].Positions = append(v1.Edges[v2].Positions, Pos{X: start + x, Y: pos1.Y})
					}
					v2Positions := slices.Clone(v1.Edges[v2].Positions)
					v2.Edges[v1] = &Path{Positions: v2Positions, Marks: []rune(str)}
					slices.Reverse(v2.Edges[v1].Positions)
				} else if pos2.X == pos1.X { // Vertically aligned vertices
					v1.Edges[v2] = &Path{Positions: []Pos{}, Marks: []rune{}}
					for y := pos1.Y + 1; y < pos2.Y; y++ {
						v.Lines[y][pos1.X] = '|'
						if _, exists := v.Index[y]; !exists {
							v.Index[y] = map[int]rune{pos1.X: '|'}
						} else {
							v.Index[y][pos1.X] = '|'
						}
						v1.Edges[v2].Positions = append(v1.Edges[v2].Positions, Pos{X: pos1.X, Y: y})
						v1.Edges[v2].Marks = append(v1.Edges[v2].Marks, '|')
					}
					v2Positions := slices.Clone(v1.Edges[v2].Positions)
					v2Marks := slices.Clone(v1.Edges[v2].Marks)
					v2.Edges[v1] = &Path{Positions: v2Positions, Marks: v2Marks}
					slices.Reverse(v2.Edges[v1].Positions)
				}
				done[v1] = append(done[v1], v2)
			}
		}
	}
	for _, v1 := range Graph.Vertices { // All other edges
		v1.SortEdgesByDegrees()
		for _, v2 := range v1.Sorted {
			if _, exists := done[v2]; !exists || !slices.Contains(done[v2], v1) {
				path, found := v.FindFreePath(v1, v2)
				if found {
					v.PutPath(path)
					v1.Edges[v2] = path
					v2Positions := slices.Clone(v1.Edges[v2].Positions)
					v2Marks := slices.Clone(v1.Edges[v2].Marks)
					v2.Edges[v1] = &Path{Positions: v2Positions, Marks: v2Marks}
					slices.Reverse(v2.Edges[v1].Positions)
					slices.Reverse(v2.Edges[v1].Marks)
				}
				done[v1] = append(done[v1], v2)
			}
		}
	}
}

func (v *visual) AddToIndex(y, x int, str rune) {
	if _, exists := v.Index[y]; !exists {
		v.Index[y] = map[int]rune{x: str}
	} else {
		v.Index[y][x] = str
	}
}

func (v *visual) PutPath(path *Path) {
	for i, pos := range path.Positions {
		if _, exists := v.Index[pos.Y]; !exists {
			v.Index[pos.Y] = map[int]rune{pos.X: path.Marks[i]}
		} else {
			v.Index[pos.Y][pos.X] = path.Marks[i]
		}
		v.Lines[pos.Y][pos.X] = path.Marks[i]
	}
}

func (v *visual) FindFreePath(v1, v2 *Vertex) (*Path, bool) {
	pos1 := v1.Position
	pos2 := v2.Position
	len1 := min(len(v1.Name), v.Zoom.X)
	len2 := min(len(v2.Name), v.Zoom.X)
	diff := Pos{X: pos2.X - pos1.X, Y: pos2.Y - pos1.Y}
	path := &Path{Positions: []Pos{}}
	if diff.Y > 0 && diff.X >= diff.Y {
		if pos2.X > pos1.X+len1+diff.Y-1 { // far right
			if v.IsFree2(pos1.Y, pos1.X+len1) && v.IsFree2(pos2.Y-1, pos2.X-1) {
				y := pos1.Y
				mark := '-'
				for x := pos1.X + len1; x < pos2.X; x++ {
					if x == pos2.X-diff.Y {
						mark = '.'
					} else if x > pos2.X-diff.Y {
						y++
						mark = '\\'
					}
					if _, isY := v.VIndex[y]; isY {
						if _, isX := v.VIndex[y][x]; isX {
							goto NEXT
						}
					}
					path.Add(Pos{x, y}, mark)
				}
				return path, true
			}
		NEXT:
			path.Positions = nil
			path.Marks = nil
			y := pos1.Y
			if v.IsFree2(pos2.Y, pos2.X-1) { // Check target position first
				var mark rune
				for x0 := pos1.X + len1; x0 > pos1.X; x0-- { // Check possible source position
					if v.IsFree2(pos1.Y+1, x0) {
						for x := x0; x < pos2.X; x++ {
							if y < pos2.Y {
								y++
								mark = '\\'
								if pos2.Y == y {
									mark = '`'
								}
							} else {
								mark = '-'
							}
							if _, isY := v.VIndex[y]; isY {
								if _, isX := v.VIndex[y][x]; isX {
									return path, false
								}
							}
							path.Add(Pos{x, y}, mark)
						}
						return path, true
					}
				}
			}
		} else {
			if v.IsFree2(pos2.Y-1, pos2.X-1) { // Check target port
				y := pos1.Y
				if pos1.X <= pos2.X-diff.Y && pos2.X-diff.Y < pos1.X+len1 && v.IsFree2(pos1.Y+1, pos2.X-diff.Y) {
					for x := pos2.X - diff.Y + 1; x < pos2.X; x++ {
						y++
						if _, isY := v.VIndex[y]; isY {
							if _, isX := v.VIndex[y][x]; isX {
								return path, false
							}
						}
						path.Add(Pos{x, y}, '\\')
					}
					return path, true
				}
			}
		}
	} else if diff.Y > 0 && diff.X > 0 && diff.Y > diff.X {
		var mark rune
		var x int
		for x = pos1.X; x < pos1.X+len1; x++ {
			mark := v.Lines[pos1.Y+1][x]
			if mark == ' ' && x >= pos1.X {
				break
			}
		}
		for y := pos1.Y + 1; y < pos2.Y; y++ {
			if y < pos2.Y-diff.X {
				mark = '|'
			} else {
				mark = '\\'
				if y == pos2.Y-diff.X {
					mark = ':'
				} else {
					x++
				}
			}
			if _, isY := v.VIndex[y]; isY {
				if _, isX := v.VIndex[y][x]; isX {
					return path, false
				}
			}
			path.Add(Pos{x, y}, mark)
		}
		return path, true
	} else if diff.X < 0 && -diff.X >= diff.Y {
		if pos2.X+len2+diff.Y <= pos1.X { // far left
			if v.IsFree2(pos1.Y, pos1.X-1) && v.IsFree2(pos2.Y-1, pos2.X+len2) {
				y := pos1.Y
				mark := '-'
				for x := pos1.X - 1; x >= pos2.X+len2; x-- {
					if x == pos2.X+len2+diff.Y-1 {
						mark = '.'
					}
					if x < pos2.X+len2+diff.Y-1 {
						mark = '/'
						y++
					}
					if _, isY := v.VIndex[y]; isY {
						if _, isX := v.VIndex[y][x]; isX {
							goto NEXT2
						}
					}
					path.Add(Pos{x, y}, mark)
				}
				return path, true
			}
		NEXT2:
			path.Positions = nil
			path.Marks = nil
			y := pos1.Y
			if v.IsFree2(pos2.Y, pos2.X+len2) { // Check target position first
				var mark rune
				for x0 := pos1.X - 1; x0 < pos1.X+len1; x0++ { // Check possible source position
					if v.IsFree2(pos1.Y+1, x0) {
						for x := x0; x >= pos2.X+len2; x-- {
							if y < pos2.Y {
								y++
								mark = '/'
								if pos2.Y == y {
									mark = '´' //'\'' //8217
								}
							} else {
								mark = '-'
							}
							if _, isY := v.VIndex[y]; isY {
								if _, isX := v.VIndex[y][x]; isX {
									return path, false
								}
							}
							path.Add(Pos{x, y}, mark)
						}
						return path, true
					}
				}
			}
		} else { // Diagonals
			if v.IsFree2(pos2.Y-1, pos2.X+len2) { // Check target port
				var x int
				for x = pos1.X - 1; x < pos2.X+len2+diff.Y; x++ {
					if v.IsFree2(pos1.Y+1, x) {
						break
					}
				}
				for y := pos1.Y + 1; y < pos2.Y; y++ {
					if _, isY := v.VIndex[y]; isY {
						if _, isX := v.VIndex[y][x]; isX {
							return path, false
						}
					}
					path.Add(Pos{x, y}, '/')
					x--
				}
				return path, true
			}
		}
	} else if diff.Y > 0 && diff.X < 0 && diff.Y > abs(diff.X) {
		var mark rune
		var X int
		for X = pos2.X; X < min(pos2.X+len2, pos1.X-1); X++ {
			mark := v.Lines[pos2.Y-1][X]
			if mark == ' ' && X >= pos2.X {
				break
			}
		}
		x := pos1.X
		for y := pos1.Y + 1; y < pos2.Y; y++ {
			if x == X {
				mark = '|'
			} else {
				mark = '/'
				x--
				if x == X {
					mark = ':'
				}
			}
			if _, isY := v.VIndex[y]; isY {
				if _, isX := v.VIndex[y][x]; isX {
					return path, false
				}
			}
			path.Add(Pos{x, y}, mark)
		}
		return path, true
	}
	return &Path{}, false
}

func (v *visual) IsFree2(y, x int) bool {
	if y < 0 || y >= len(v.Lines) || x < 0 || x >= len(v.Lines[y]) {
		return false
	}
	return v.Lines[y][x] == ' '
}

func (v *visual) MakeBase() {
	v.PutVertices()
	v.PutEdges()
	v.MakeColorIndex()
	//fmt.Println(v.ColorIndex)
}

func (v *visual) MakeColorIndex() {
	v.ColorIndex = make(map[int]map[int]string)
	for i, path := range Graph.Paths {
		for j := 0; j < len(path)-1; j++ {
			v1 := path[j]
			v2 := path[j+1]
			edge := v1.Edges[v2]
			if edge != nil {
				//fmt.Println(v1.Name, v2.Name, "ERROR: Edge not found (visual.go, 387)")
				//} else {
				for _, pos := range edge.Positions {
					if _, exists := v.ColorIndex[pos.Y]; !exists {
						v.ColorIndex[pos.Y] = map[int]string{pos.X: colors[i%len(colors)]}
					} else {
						v.ColorIndex[pos.Y][pos.X] = colors[i%len(colors)]
					}
				}
			}
		}
	}
}

/*
	func (v *visual) PrintBase() {
		for _, line := range v.Lines {
			//fmt.Println("|" + string(line) + "|")
			fmt.Println(string(line))
		}
	}
*/
func (v *visual) PrintBase() {
	for y, line := range v.Lines {
		str := ""
		c := 0
		for x, r := range line {
			if color, exists := v.ColorIndex[y][c]; exists {
				str += color + string(line[x]) + Reset
			} else {
				str += string(line[x])
			}
			x++
			c++
			if r == '\U0001F41C' {
				c++
			}
		}
		fmt.Println(str)
	}
}
func (v *visual) Show() {
	v.SetDimensions()
	v.MakeVIndex()
	v.MakeBase()
	v.PrintBase()

	fmt.Println("Press ENTER to run all, SPACE+ENTER to step through, Ctrl-C to exit...")
	b, err := bufio.NewReader(os.Stdout).ReadByte()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var action string
	switch b {
	case 32:
		action = "STEP"
	case 10:
		action = "RUN"
	case 27:
		os.Exit(1)
	default:
		fmt.Printf("Unknown action: %v\n", b)
		os.Exit(1)
	}
	if action != "" {
		time.Sleep(1 * time.Second)

		fmt.Println("\033[2J")
		fmt.Println("\033[H")
		switch action {
		case "STEP":
			v.Step()
		case "RUN":
			v.Run(0)
		}
	}
	time.Sleep(1 * time.Second)
}

func (v *visual) MoveAnts(turn *Turn, lines [][]rune) {
	// Restore clean maze
	frame := 0
	for {
		v.Lines = deepCopy(lines)
		found := false
		for ant, vertex := range turn.EnRoute {
			bigPath := Graph.Paths[ant.Queue]       // Path to follow
			currentVertex := bigPath[ant.Current]   // Current vertex
			stepPath := currentVertex.Edges[vertex] // Path to next vertex
			if stepPath == nil {
				fmt.Println("ERROR: Edge not found (visual.go, 433)", currentVertex.Name, vertex.Name)
				os.Exit(1)
			}
			var pos Pos
			if frame >= len(stepPath.Positions) {
				pos = vertex.Position
			} else {
				pos = stepPath.Positions[frame]
				found = true
			}
			if pos.Y < len(v.Lines) && pos.X < len(v.Lines[pos.Y]) && pos != Graph.End.Position {
				v.Lines[pos.Y][pos.X] = '\U0001F41C' //ant
			}
		}
		for y, l := range v.Lines {
			line := slices.Clone(l)
			v.Lines[y] = nil
			skip := false
			for _, c := range line {
				//fmt.Print(c)
				if !skip {
					v.Lines[y] = append(v.Lines[y], c)
				}
				skip = false
				if c == '\U0001F41C' {
					skip = true
				}
			}
		}

		fmt.Println("\033[2J") // Clear screen
		fmt.Println("\033[H")  // Move cursor to top left corner
		v.PrintBase()
		time.Sleep(200 * time.Millisecond)
		frame++
		if !found {
			break
		}
	}
	for ant, _ := range turn.EnRoute {
		ant.Current++
	}
}

func (v *visual) Step() {
	lines := deepCopy(v.Lines)
	for i, turn := range Turns.Data {
		v.MoveAnts(turn, lines)
		fmt.Println("Press ENTER to run all, SPACE+ENTER to step through, Ctrl-C to exit...")
		b, err := bufio.NewReader(os.Stdout).ReadByte()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		switch b {
		case 32:
			continue
		case 10:
			v.Lines = deepCopy(lines)
			v.Run(i + 1)
		case 27:
			os.Exit(1)
		default:
			fmt.Printf("Unknown action: %v\n", b)
			os.Exit(1)
		}
	}
}

func (v *visual) Run(t int) {
	lines := deepCopy(v.Lines)
	for i := t; i < len(Turns.Data); i++ {
		turn := Turns.Data[i]
		v.MoveAnts(turn, lines)
		time.Sleep(1 * time.Second)
	}
}

var Visual = &visual{
	Index:  map[int]map[int]rune{},
	Lines:  [][]rune{},
	Zoom:   Pos{X: 4, Y: 2},
	Sorted: map[int]map[int]string{},
	MaxDim: Pos{X: 0, Y: 0},
}

func deepCopy(slice [][]rune) [][]rune {
	lines := make([][]rune, len(slice))
	for i := 0; i < len(slice); i++ {
		line := make([]rune, len(slice[0]))
		copy(line, slice[i])
		lines[i] = line
	}
	return lines
}

func deepCopyVertices(slice [][]*Vertex) [][]*Vertex {
	lines := make([][]*Vertex, len(slice))
	for i := 0; i < len(slice); i++ {
		line := make([]*Vertex, len(slice[0]))
		copy(line, slice[i])
		lines[i] = line
	}
	return lines
}
