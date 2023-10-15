package types

type explorer struct {
	visited map[*Vertex]bool
	current *Vertex
	path    []*Vertex
}

type PathData struct {
	Paths []string `json:"paths"`
}

func (explorer *explorer) Explore() {
	var prev *Vertex
	// If new search...
	if explorer.current == nil {
		// Put explorer on start vertex
		explorer.current = Graph.Start
	}
	// Explore each link from current vertex
	for v, _ := range explorer.current.Edges {
		if explorer.current == Graph.Start {
			// If on start vertex, reset path, with start vertex as first element
			explorer.path = []*Vertex{Graph.Start}
			explorer.visited[Graph.Start] = true
		}
		// If this vertex has not been visited on this path before...
		if !explorer.visited[v] {
			// Add vertex to path
			explorer.path = append(explorer.path, v)
			// Mark vertex as visited
			explorer.visited[v] = true
			// Record current vertex as previous
			// prev has to be local to this loop, otherwise it will be overwritten
			prev = explorer.current
			// Set current vertex to this vertex
			explorer.current = v
			if v == Graph.End {
				// Make a copy of the path
				tmp := make([]*Vertex, len(explorer.path))
				copy(tmp, explorer.path)
				// Add path to all paths
				Paths.All = append(Paths.All, tmp)
				AllPaths = append(AllPaths, tmp)
			} else {
				// Explore recursively from this vertex
				explorer.Explore()
			}
			// If we are here, we have reached the end of current path,
			// nothing more to explore: either end or dead end.
			// We are now stepping back in the recursive tree
			// Mark vertex as not visited (for further paths to be able to visit it)
			explorer.visited[v] = false
			// Set current vertex to previous
			explorer.current = prev
			// Remove last vertex from path
			if len(explorer.path) > 1 {
				explorer.path = explorer.path[:len(explorer.path)-1]
			}
		}
	}
}

var Explorer = &explorer{current: Graph.Start, visited: make(map[*Vertex]bool)}
