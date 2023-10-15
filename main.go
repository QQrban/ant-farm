package main

import (
	"bufio"
	"flag"
	"fmt"
	"lem-in/types"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// >>>>> MAIN >>>>>
func main() {
	var (
		zoomX int
		zoomY int
		//visual bool
		ants             int
		webVisualisation bool
	)
	flag.BoolVar(&webVisualisation, "web", false, "Enable web visualisation")
	flag.IntVar(&zoomX, "zoomX", 4, "Zoom x dimension to unclutter view")
	flag.IntVar(&zoomX, "x", 4, "Zoom x dimension to unclutter view")
	flag.IntVar(&zoomY, "zoomY", 2, "Zoom y dimension to unclutter view")
	flag.IntVar(&zoomY, "y", 2, "Zoom y dimension to unclutter view")
	//flag.BoolVar(&visual, "visual", false, "Visualize ants movements through graph")
	//flag.BoolVar(&visual, "v", false, "Visualize ants movements through graph")
	flag.IntVar(&ants, "ants", 0, "Number of ants (default 0)")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}
		if fi.Mode()&os.ModeNamedPipe == 0 {
			fmt.Println("Error: missing filename")
			fmt.Println("USAGE: go run . [-ants <int>] <filename|path>")
			fmt.Println("VISUALISATION: go run . [-ants <int>] <filename|path> | go run . [-zoomX <int>|-x <int>] [-zoomY <int>|-y <int>]")
			os.Exit(1)
		}
	}

	types.Visual.Zoom = types.Pos{X: zoomX, Y: zoomY}
	//types.Visual.Animate = visual
	types.Ants.Number = ants

	start := time.Now()
	fileName := parseInput(args)
	if types.Visual.Auto {
		types.Visual.Show()
	} else {
		types.Paths.Find()
		//printBestPaths()
		types.Ants.Distribute()
		//printDistribution()
		if types.Visual.Animate {
			types.Visual.PrintBase()
		} else {
			printInput(fileName, webVisualisation)
		}
	}
	types.Ants.Move(webVisualisation)

	if webVisualisation {
		fmt.Println("Server running on port 8080")
		types.PrintAllPaths()
		http.Handle("/", http.FileServer(http.Dir("./static")))
		http.HandleFunc("/paths", homeHandler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("Server failed to start: ", err)
		}
	} else {
		if !types.Visual.Animate {
			fmt.Println()
			fmt.Printf("Moved %d ants along %v disjoint paths in %v turns.\n", types.Ants.Number, len(types.Graph.Paths), types.Graph.Turns)
			fmt.Printf("Found altogether %v paths, %v best paths in %v.\n", len(types.Paths.All), len(types.Graph.Paths), time.Since(start))
			fmt.Println("To see visualisation, follow instructions in readme.\n")
		}
	}
}

// >>>>> SENDING ANTS >>>>>
func printDistribution() {
	fmt.Println("Distribution:")
	for i, queue := range types.Ants.Queues {
		fmt.Printf("Queue %d (%d): ", i, len(queue))
		for _, ant := range queue {
			fmt.Printf("%s, ", ant.Name)
		}
		fmt.Println()
	}
}

// >>>>> PARSING >>>>>
func parseInput(args []string) string {
	// Parse input
	var file *os.File
	fileName := ""
	var scanner *bufio.Scanner //= bufio.NewScanner(os.Stdin) //
	if len(args) > 0 {
		file, fileName = openFile(args)
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}
	// Parse ants
	scanner.Scan()
	line := scanner.Text()
	getAnts(line)
SCANNING:
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Set start and end vertices
			if line == "##start" || line == "##end" {
				scanner.Scan()
				nextLine := scanner.Text()
				vertex := setVertex(nextLine)
				vertex.Capacity = types.Ants.Number
				if line == "##start" {
					types.Graph.Start = vertex
					types.StartJSON = vertex.Name
				} else {
					types.Graph.End = vertex
				}
			}
		} else if len(line) > 0 && line[0] == 'L' {
			// turns are established already, make visualizer
			types.Visual.Animate = true
			types.Visual.Auto = true
			types.Turns.InStart = types.Ants.Number
			types.Turns.Parse(line)
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) > 0 && line[0] == 'L' {
					types.Turns.Parse(line)
				} else {
					types.Turns.ExtractPaths()
					break SCANNING
				}
			}
		} else if strings.Contains(line, " ") {
			setVertex(line)
		} else if strings.Contains(line, "-") {
			setEdge(line)
		}
	}
	types.Graph.Check()
	//fmt.Println(types.Graph.Edges)
	return fileName
}

func openFile(args []string) (*os.File, string) {
	if len(args) < 1 {
		fmt.Println("No file specified.")
		os.Exit(1)
	}
	fileName := args[0]

	if !strings.Contains(fileName, "/") {
		fileName = "examples/" + fileName
	}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file.")
		os.Exit(1)
	}
	return file, fileName
}

func getAnts(line string) {
	number, err := strconv.Atoi(line)
	if err != nil {
		types.FaultyData(fmt.Sprintf("invalid number of ants: %s", line))
	}

	if !FlagIsPassed("ants") {
		types.Ants.Number = number
	}
	number = types.Ants.Number
	if number < 1 || number > math.MaxInt32 {
		types.FaultyData(fmt.Sprintf("invalid number of ants: %d", number))
	}
	// To be able to get references to ants if turns are already calculated, i.e. automatic visualisation.
	types.Ants.All = make(map[string]*types.Ant)
	for i := 0; i < number; i++ {
		ant := &types.Ant{}
		ant.Name = strconv.Itoa(i + 1)
		types.Ants.All[ant.Name] = ant
	}
}

func setVertex(line string) *types.Vertex {
	vertex := &types.Vertex{
		Edges:  make(map[*types.Vertex]*types.Path),
		Sorted: make([]*types.Vertex, 0),
	}
	fields := strings.Split(line, " ")
	vertex.Name = fields[0]
	vertex.Position.X, _ = strconv.Atoi(fields[1])
	vertex.Position.Y, _ = strconv.Atoi(fields[2])
	vertex.Capacity = 1
	types.Graph.Vertices = append(types.Graph.Vertices, vertex)
	types.Graph.Edges[vertex.Name] = []string{}
	return vertex
}

func setEdge(line string) {
	fields := strings.Split(line, "-")
	v1Name := fields[0]
	v2Name := fields[1]
	if v1Name == v2Name {
		types.FaultyData(fmt.Sprintf("Invalid edge: %s", line))
	}
	vertex1 := types.Graph.FindVertex(v1Name)
	vertex2 := types.Graph.FindVertex(v2Name)
	if _, exists := vertex1.Edges[vertex2]; !exists {
		vertex1.Edges[vertex2] = &types.Path{}
		vertex2.Edges[vertex1] = &types.Path{}
		vertex1.Sorted = append(vertex1.Sorted, vertex2)
		vertex2.Sorted = append(vertex2.Sorted, vertex1)
		types.Graph.Edges[v1Name] = append(types.Graph.Edges[v1Name], v2Name)
		types.Graph.Edges[v2Name] = append(types.Graph.Edges[v2Name], v1Name)
		pair := types.EdgePair{
			V1Name: v1Name,
			V2Name: v2Name,
		}
		types.EdgePairs = append(types.EdgePairs, pair)
	}
}

// >>>>> PRINTING >>>>>
func printInput(fileName string, webVisualisation bool) {
	if !webVisualisation {
		input, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println("Error reading file")
			os.Exit(0)
		}
		spec := string(input)
		i := strings.Index(spec, "\n")
		fmt.Printf("%d%s\n\n", types.Ants.Number, spec[i:])
		//fmt.Printf("%s\n\n", string(input))
	}
}

func printBestPaths() {
	fmt.Println("Best paths:")
	fmt.Println("Ants:", types.Ants.Number, "Turns:", types.Graph.Turns)
	for i, path := range types.Graph.Paths {
		fmt.Printf("Path %d (%d): ", i, len(path)-1)
		for _, vertex := range path {
			fmt.Printf("%s, ", vertex.Name)
		}
		fmt.Println()
	}
}

func pyramid(n int) []int {
	// Returns lengths for each line in pyramid
	// depending on number of elements n
	var lines []int
	levels := int(math.Sqrt(float64(n)))
	canonical := levels * levels
	baas := ((levels-1)*levels*2 + levels) / levels
	extra := n - canonical
	var add int
	for i := 1; i <= levels; i++ {
		if extra >= (levels + i) {
			add = 2
		} else if extra >= i {
			add = 1
		}
		lines = append(lines, baas+add)
		baas = baas - 2
	}
	return lines
}

func printTurnData() {
	for i, turn := range types.Turns.Data {
		fmt.Printf("Turn %d: in start: %d, in end %d\n", i+1, turn.InStart, turn.InEnd)
		for ant, vertex := range turn.EnRoute {
			fmt.Printf("Ant %s: %s\n", ant.Name, vertex.Name)
		}
	}
	for i, path := range types.Graph.Paths {
		fmt.Printf("Path %d (%d): ", i, len(path))
		for _, vertex := range path {
			fmt.Printf("%s, ", vertex.Name)
		}
		fmt.Println()
	}
	for i, queue := range types.Ants.Queues {
		fmt.Printf("Queue %d (%d): ", i, len(queue))
		for _, ant := range queue {
			fmt.Printf("%s, ", ant.Name)
		}
		fmt.Println()
	}
}

func FlagIsPassed(name string) bool {
	// Check if a flag is passed
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
