package types

var AllPaths [][]*Vertex

type EdgePair struct {
	V1Name string
	V2Name string
}
type Position struct {
	X int
	Y int
}
type VertexInfo struct {
	Name     string
	Position Position
}
type DataObject struct {
	Start    string
	Paths    []string
	Vertices []VertexInfo
	Edges    []string
	AntMoves [][]string
}

var StartJSON string
var AllData []DataObject
var EdgePairs []EdgePair

func PrintAllPaths() {
	var paths []string
	var vertexInfos []VertexInfo
	var edges []string
	seen := make(map[string]bool)
	for _, path := range AllPaths {
		var s string
		for _, vertex := range path {
			s += vertex.Name + " "
			name := vertex.Name
			if !seen[name] {
				seen[name] = true
				vertexInfo := VertexInfo{
					Name: vertex.Name,
					Position: Position{
						X: vertex.Position.X,
						Y: vertex.Position.Y,
					},
				}
				vertexInfos = append(vertexInfos, vertexInfo)
			}
		}
		paths = append(paths, s)
	}
	for _, e := range EdgePairs {
		edges = append(edges, e.V1Name+"-"+e.V2Name)
	}
	data := DataObject{
		Edges:    edges,
		Paths:    paths,
		Vertices: vertexInfos,
		AntMoves: AllMoves,
		Start:    StartJSON,
	}
	AllData = append(AllData, data)
}
