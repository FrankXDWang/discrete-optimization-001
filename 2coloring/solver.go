package main

import "sort"
import "fmt"
import "os"
//import "container/list"
//import "io"
//import "bytes"
//import "container/heap"
//import "encoding/binary"
//import "compress/gzip"

// functions which Go developers should have implemented but happened
// to be too lazy and religious to do so

func max(a int32, b int32) (r int32) {
    if a > b {
        return a
    } else {
        return b
    }
}

type Edge struct {
    u int32 // first vertex id
    v int32 // second vertex id
}

type Vertex struct {
    index int32 // original index
    color int32 // vexter color
    E []int32   // list of connected edges
}

type Edges []Edge
type Vertices []Vertex

// graph contains of edges and vertices
type Graph struct {
    E Edges
    V Vertices
}

type VarHeuristic int
type ValHeuristic int

const (
    VAR_BRUTE VarHeuristic = iota
    VAR_MRV
    VAR_MCV
)

const (
    VAL_BRUTE ValHeuristic = iota
    VAL_LCV
)

// additional information (besides the graph), required for CSP
type CSPContext struct {
    g *Graph
    //domain [][]int32 // possible values (colors) for each variable (vertex)
    domains []map[int32]bool // possible values (colors) for each variable (vertex)
    numColors int32  // target number of colors (does not change)
    currentUnassignedVertex int // current vertex in recursive solution calls
    varHeuristic VarHeuristic
    valHeuristic ValHeuristic
}

// save vertex order without reordering graph vertices
type VertexOrder struct {
    g *Graph
    order []int32
}

// v -- global index of the vertex
// e -- local index of the edge in vertex edge list
// return global index of the other vertex
func (self *Graph) otherVertex(v int32, e int32) int32 {
    V := self.V[v]
    E := self.E[V.E[e]]
    if E.v == v {
        return E.u
    } else {
        return E.v
    }
}

type ByInt32 []int32
func (self ByInt32) Len() int { return len(self) }
func (self ByInt32) Less(i, j int) bool { return self[i] < self[j] }
func (self ByInt32) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

type ByIndex VertexOrder
func (self ByIndex) Len() int { return len(self.order) }
func (self ByIndex) Less(i, j int) bool { return self.order[i] < self.order[j] }
func (self ByIndex) Swap(i, j int) { self.order[i], self.order[j] = self.order[j], self.order[i] }

type ByDegree VertexOrder
func (self ByDegree) Len() int { return len(self.order) }
func (self ByDegree) Less(i, j int) bool { return len(self.g.V[self.order[i]].E) < len(self.g.V[self.order[j]].E) }
func (self ByDegree) Swap(i, j int) { self.order[i], self.order[j] = self.order[j], self.order[i] }

// type ByDegree VertexOrder
// func (self ByDegree) Len() int { return len(self) }
// func (self ByDegree) Less(i, j int) bool { return len(self[i].E) < len(self[j].E) }
// func (self ByDegree) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

func (g *Graph) NV() int { return len(g.V) }
func (g *Graph) NE() int { return len(g.E) }

func (g *Graph) degree() int32 {
    var maxDegree int32 = 0
    for i := 0; i < len(g.V); i++ {
        maxDegree = max(maxDegree, int32(len(g.V[i].E)))
    }
    return maxDegree
}

func (g *Graph) chromaticNumber() int32 {
    var maxChNum int32 = 0
    for i := 0; i < len(g.V); i++ {
        maxChNum = max(maxChNum, int32(g.V[i].color))
    }
    return maxChNum
}

func (g *Graph) vertexNeighborColors(i int32) []int32 {
    neibColors := make([]int32, 0)
    // get colors of all neighbors
    for j := 0; j < len(g.V[i].E); j++ {
        neibVertex := g.V[g.otherVertex(i, int32(j))]
        neibColors = append(neibColors, neibVertex.color)
    }
    return neibColors
}

func minUnusedColor(colors *[]int32) int32 {
    sort.Sort(ByInt32(*colors))
    //fmt.Println(*colors)

    if len(*colors) == 1 {
        return (*colors)[0] + 1
    }

    for i := 0; i < len(*colors) - 1; i++ {
        if (*colors)[i+1] - (*colors)[i] > 1 {
            return (*colors)[i] + 1
        }
    }

    return (*colors)[len(*colors) - 1] + 1
}

func (g *Graph) assignVertexColor(i int32) {
    neibColors := g.vertexNeighborColors(i)
    // find min unused color
    min_color := minUnusedColor(&neibColors)
    g.V[i].color = min_color
    //fmt.Println(min_color)
}

func (g *Graph) printColors() {
    for i := 0; i < len(g.V); i++ {
        if (i != len(g.V) - 1) {
            //fmt.Printf("%d (%d) ", g.V[i].color - 1, g.V[i].index)
            fmt.Printf("%d ", g.V[i].color - 1)
        } else {
            fmt.Printf("%d", g.V[i].color - 1)
        }
    }
    fmt.Printf("\n")
}

func (g *Graph) printSolution() {
    fmt.Println(g.chromaticNumber(), 0)
    g.printColors()
}

// greedy approach
func (g *Graph) solveGreedySimple() {
    //NE := len(g.E)
    NV := len(g.V)
    //D := degree(&g)

    ord := make([]int32, NV)
    for i := 0; i < NV; i++ {
        ord[i] = int32(i)
    }
    vertexOrder := VertexOrder{g, ord}

    sort.Sort(sort.Reverse(ByDegree(vertexOrder)))

    //sort.Sort(sort.Reverse(ByDegree(g.V)))
    //sort.Sort(ByDegree(g.V))
    //sort.Reverse(ByDegree(g.V))

    //fmt.Println(D)

    for i := 0; i < NV; i++ {
        g.assignVertexColor(vertexOrder.order[i])
    }

    //sort.Sort(ByIndex(g.V))

    //fmt.Println(g.chromaticNumber(), 0)
    //g.printColors()
    g.printSolution()
}

//
// CSP
//

func (c *CSPContext) init(nColors int) {
    c.domains = make([]map[int32]bool, c.g.NV())
    for i := 0; i < c.g.NV(); i++ {
        c.domains[i] = make(map[int32]bool)
        for j := 0; j < nColors; j++ {
            c.domains[i][int32(j + 1)] = true
        }
    }
    c.currentUnassignedVertex = 0
}

func (v Vertex) numSameColorNeighbors(g *Graph) int {
    num := 0
    // check all neighbor vertices
    for i := 0; i < len(v.E); i++ {
        otherVertexIndex := g.otherVertex(v.index, int32(i))
        if g.V[otherVertexIndex].color == v.color {
            num += 1
        }
    }
    return num
}

// check if vertex is valid
func (v Vertex) valid(g *Graph) bool {
    // unassigned color?
    if v.color == 0 {
        return false
    }

    return v.numSameColorNeighbors(g) == 0
}

// check if graph is valid
func (g *Graph) valid() bool {
    // check if all vertices are valid
    for i := 0; i < g.NV(); i++ {
        if !g.V[i].valid(g) {
            return false
        }
    }

    return true
}

// func (c *CSPContext) tryNextVal(vertex int) bool {
//     //fmt.Println("Trying vertex", vertex)
//     // try all colors for current vertex
//     for color := 0; color < int(c.numColors); color++ {
//         c.g.V[vertex].color = c.domain[vertex][color]
//         //fmt.Println(c.g)
//         //c.g.printSolution()
//         if c.solve() {
//             return true
//         }
//         c.currentUnassignedVertex -= 1
//     }
// 
//     return false
// }
// 
// func (c *CSPContext) tryNextVar() bool {
//     switch {
//     case c.varHeuristic == VAR_BRUTE:
//         vertex := c.currentUnassignedVertex
//         c.currentUnassignedVertex += 1
//         if c.tryNextVal(vertex) {
//             return true
//         }
//         c.currentUnassignedVertex -= 1
// 
//     case c.varHeuristic == VAR_MRV:
//         fmt.Println("VAR_MRV not implemented")
// 
//     case c.varHeuristic == VAR_MCV:
//         fmt.Println("VAR_MCV not implemented")
//     }
// 
//     return false
// }

func (c *CSPContext) forwardCheckVertexColor(vertex int32, color int32) {
        for j := 0; j < len(c.g.V[vertex].E); j++ {
        neibVertexIndex := c.g.otherVertex(vertex, int32(j))
        delete(c.domains[neibVertexIndex], color)
        //neibColors = append(neibColors, neibVertex.color)
    }
}

func (c *CSPContext) getMRVVertex() int32 {
    // there must be unset vertices
    if c.currentUnassignedVertex >= c.g.NV() {
        panic("Call to getMRVVertex with no unset variables")
        return -1
    }

    var vertex int32 = -1

    // scan all domains, find the smallest one
    // (number of vertices == number of domains)
    for i := 0; i < c.g.NV(); i++ {
        // only check unset vertices
        if c.g.V[i].color != 0 {
            continue
        }

        // assign vertex if not yet assigned
        if vertex == -1 {
            vertex = int32(i)
        }

        if len(c.domains[i]) < len(c.domains[vertex]) {
            vertex = int32(i)
        }
    }

    if vertex == -1 {
        panic("getMRVVertex could not find the vertex")
        return -1
    }
    return vertex;
}

func (c *CSPContext) solve() bool {
    // all vars assigned?
    if c.currentUnassignedVertex >= c.g.NV() {
        return c.g.valid()
    }

    // TODO: support MRV (minimum remaining values):
    // 1.+use maps instead of lists for domains (faster deletion)
    // 2.+forward check domain changes to neighbors after assigning
    //    color to the vertex
    // 3.+select MRV vertex (scan all vertices and select min)
    // 4. try LCV (for values)
    // 5. try constraint propagation (stronger version of forward checking)

    // select var
    vertex := c.getMRVVertex() //c.currentUnassignedVertex
    c.currentUnassignedVertex += 1

    // save a copy of state of all current domains
    savedDomains := pushDomains(c.domains)

    for color, _ := range c.domains[vertex] {
        // set another color
        c.g.V[vertex].color = color

        // propagate color. this will change current domains state
        c.forwardCheckVertexColor(int32(vertex), color)

        //fmt.Println(c.g)
        //c.g.printSolution()
        if c.solve() {
            return true
        }

        // restore domains state to previous
        popDomains(&c.domains, &savedDomains)
    }

    /*
    for color := 0; color < int(c.numColors); color++ {
        c.g.V[vertex].color = c.domain[vertex][int32(color)]
        //fmt.Println(c.g)
        //c.g.printSolution()
        if c.solve() {
            return true
        }
    }
    */

    // unselect var
    c.currentUnassignedVertex -= 1
    c.g.V[vertex].color = 0

    return false
}

// contraint-satisfaction approach
func (g *Graph) solveCSP() {
    //nColors := g.degree() + 1
    nColors := int32(78)

    //fmt.Println("Solving for", nColors, "colors")

    csp := CSPContext{g, nil, nColors, 0, VAR_BRUTE, VAL_BRUTE}
    csp.init(int(nColors))
    //fmt.Println(csp)
    if csp.solve() {
        //fmt.Println("SOLUTION")
        g.printSolution()
    }
}

func solveFile(filename string, alg string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Cannot open file:", filename, err)
        return
    }
    defer file.Close()

    var NV, NE int32
    var i, v, u int32

    fmt.Fscanf(file, "%d %d", &NV, &NE)

    //v := make([]int32, n)
    E := make([]Edge, NE)
    V := make([]Vertex, NV)

    for i = 0; i < NV; i++ {
        V[i] = Vertex{int32(i), 0, make([]int32, 0)}
    }

    for i = 0; i < NE; i++ {
        fmt.Fscanf(file, "%d %d", &v, &u)
        E[i] = Edge{v, u}
        V[v].E = append(V[v].E, i)
        V[u].E = append(V[u].E, i)
    }

    g := Graph{E, V}

    //g.solveGreedySimple()
    //g.solveCSP()

    //fmt.Println(E)
    //fmt.Println(V)

    switch {
    case alg == "estimate":
        fmt.Println("Estimation is not implemented in this assignment")
        //fmt.Println("DP estimated memory usage, MB:",
        //            (int(K+1) * int(n+1) * 4 + int(n)) / 1024 / 1024)
    case alg == "greedy":
        g.solveGreedySimple()
    case alg == "csp":
        g.solveCSP()
    default:
        g.solveCSP()
    }
}

func pushDomains(domains []map[int32]bool) []map[int32]bool {
    newDomains := make([]map[int32]bool, len(domains))
    for i := 0; i < len(domains); i++ {
        newDomains[i] = make(map[int32]bool)
        for k, v := range domains[i] {
            newDomains[i][k] = v
        }
    }
    return newDomains
}

func popDomains(dst *[]map[int32]bool, src *[]map[int32]bool) {
    newDomains := pushDomains(*src)
    *dst = newDomains
}

func test(alg string) {
    N := 2
    d := make([]map[int32]bool, N)

    d[0] = make(map[int32]bool)
    d[1] = make(map[int32]bool)

    d[0][1] = true
    d[0][2] = true
    d[0][3] = true

    d[1][2] = true

    saved := pushDomains(d)

    d[1][1] = true
    d[1][3] = true

    fmt.Println(d)
    fmt.Println(saved)

    popDomains(&d, &saved)

    fmt.Println(d)
    fmt.Println(len(d[0]))
}

func main() {
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    solveFile(os.Args[1], alg)
    //test(alg)

    /*
    c1 := []int32{1, 2, 3}
    fmt.Println(minUnusedColor(&c1))
    c2 := []int32{0, 2, 3}
    fmt.Println(minUnusedColor(&c2))
    c3 := []int32{1, 2, 3, 4, 6, 7, 8, 9}
    fmt.Println(minUnusedColor(&c3))
    */
}


// type Node struct {
//     index int32   // index in the input data
//     value int32
//     weight int32
//     bound float32 // this is used as priority
//     selected byte
//     sel []byte
// }
// 
// // Priority queue -------------------------------------------------------------
// 
// type Items []Node
// 
// func (self Items) Len() int { return len(self) }
// func (self Items) Less(i, j int) bool { return self[i].bound < self[j].bound }
// func (self Items) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
// func (self *Items) Push(x interface{}) { *self = append(*self, x.(Node)) }
// func (self *Items) Pop() (popped interface{}) {
//     popped = (*self)[len(*self)-1]
//     *self = (*self)[:len(*self)-1]
//     return
// }
// 
// // Sorting --------------------------------------------------------------------
// 
// type ByValuePerWeight Items
// func (self ByValuePerWeight) Len() int { return len(self) }
// func (self ByValuePerWeight) Less(i, j int) bool {
//     a := float32(self[i].value) / float32(self[i].weight)
//     b := float32(self[j].value) / float32(self[j].weight)
//     return a > b
// }
// func (self ByValuePerWeight) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
// 
// type ByIndex Items
// func (self ByIndex) Len() int { return len(self) }
// func (self ByIndex) Less(i, j int) bool { return self[i].index < self[j].index }
// func (self ByIndex) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
// 
// // Branch and Bound -----------------------------------------------------------
// 
// func (node *Node) estimate(K int32, N int32, items Items) float32 {
//     var j, k int32
//     var totweight int32
//     var result float32
// 
//     if node.weight >= K {
//         return 0
//     }
// 
//     result = float32(node.value)
//     totweight = node.weight
//     j = node.index + 1
// 
//     for j < N && totweight + items[j].weight <= K {
//         totweight += items[j].weight
//         result += float32(items[j].value)
//         j++
//     }
// 
//     k = j
//     if k < N {
//         result += float32((K - totweight) * items[k].value / items[k].weight)
//     }
// 
//     return result
// }
// 
// // see
// // http://books.google.ru/books?id=QrvsNy9paOYC&pg=PA235&lpg=PA235&dq=knapsack+problem+branch+and+bound+C%2B%2B&source=bl&ots=e6ok2kODMN&sig=Yh5__d3iAFa5rEkaCoBJ2JAWybk&hl=en&sa=X&ei=k1EDULDrHIfKqgHqtYyxDA&redir_esc=y#v=onepage&q&f=true
// 
// func knapsackBranchAndBound(K int32, items Items, maxvalue *int32) []byte {
//     var N int32 = int32(len(items))
//     var u, v Node
//     //var x = make([]byte, N) // currently selected items
//     var bestset = make([]byte, N) // best selected items
//     pq := &Items{}
// 
//     heap.Init(pq)
//     *maxvalue = 0
// 
//     // initialize root
//     u = Node{0, 0, 0, 0, 0, make([]byte, N)}
//     // index = -1, start with fake root node
//     v = Node{-1, 0, 0, 0, 0, make([]byte, N)}
//     v.bound = v.estimate(K, N, items)
//     heap.Push(pq, v)
// 
//     for pq.Len() != 0 {
//         v = heap.Pop(pq).(Node)
//         if v.bound > float32(*maxvalue) {
//             // make child that includes the item
//             u = Node{v.index+1,
//                      v.value + items[v.index+1].value,
//                      v.weight + items[v.index+1].weight,
//                      0,
//                      0,
//                      make([]byte, N)}
// 
//             copy(u.sel, v.sel)
//             u.sel[u.index] = 1
// 
//             if u.weight <= K && u.value > *maxvalue {
//                 *maxvalue = u.value
//                 copy(bestset, u.sel)
//             }
//             u.bound = u.estimate(K, N, items)
//             if u.bound > float32(*maxvalue) {
//                 heap.Push(pq, u)
//             }
// 
//             // make child that does not include the item
//             u = Node{v.index+1,
//                      v.value,
//                      v.weight,
//                      0,
//                      0,
//                      make([]byte, N)}
//             u.bound = u.estimate(K, N, items)
//             copy(u.sel, v.sel)
//             u.sel[u.index] = 0
// 
//             if u.bound > float32(*maxvalue) {
//                 heap.Push(pq, u)
//             }
//         }
//     }
// 
//     return bestset
// }
// 
// func solveBranchAndBound(K int32, v []int32, w []int32) {
//     N := len(v)
//     items := make([]Node, N)
//     for i := 0; i < N; i++ {
//         items[i] = Node{int32(i), v[i], w[i], -1, 0, make([]byte, N)}
//     }
//     sort.Sort(ByValuePerWeight(items))
// 
//     var maxvalue int32 = -1
//     bestset := knapsackBranchAndBound(K, items, &maxvalue)
//     fmt.Println(maxvalue, 1) // not always optimal (1), actually
// 
//     // restore indexes
//     for i := 0; i < N; i++ {
//         items[i].selected = bestset[i]
//     }
//     sort.Sort(ByIndex(items))
//     for i := 0; i < N; i++ {
//         if items[i].selected == 1 {
//             fmt.Printf("1")
//         } else {
//             fmt.Printf("0")
//         }
//         if i != N-1 { // last?
//             fmt.Printf(" ")
//         }
//     }
//     fmt.Printf("\n")
// }
// 
// // Dynamic Programming --------------------------------------------------------
// 
// func dumpToFile(file *os.File, data []int32) int64 {
//     // write data into bytes Buffer
//     var buf bytes.Buffer
//     binary.Write(&buf, binary.LittleEndian, data)
// 
//     // prepare packed buffer
//     var packedBuf bytes.Buffer
//     z := gzip.NewWriter(&packedBuf)
// 
//     // write unpacked to packed through gzip
//     buf.WriteTo(z)
//     z.Flush()
//     z.Close()
// 
//     size := int64(packedBuf.Len())
//     packedBuf.WriteTo(file)
//     return size
// }
// 
// func loadFromFile(file *os.File, size int, data *[]int32) {
//     var packedIn bytes.Buffer
//     packedIn.ReadFrom(io.LimitReader(file, int64(size)))
// 
//     unz, _ := gzip.NewReader(&packedIn)
//     var unpacked bytes.Buffer
//     unpacked.ReadFrom(unz)
// 
//     //var dataIn = make([]int32, N)
//     binary.Read(&unpacked, binary.LittleEndian, data)
// }
// 
// func solveDynamicProgramming(K int32, v []int32, w []int32) {
//     var N int32 = int32(len(v))
// 
//     file, _ := os.Create("dptable.bin")
//     var offsets = make([]int64, N+1) // store offsets of dumped columns in the file
//     var sizes = make([]int, N+1)
//     var position int64
//     var lastPackedSize int64
// 
//     var O = make([][]int32, 2)
//     var x = make([]byte, N)
//     // create O lookup table
// 
//     var k int32
//     var j int32
//     var i int32
// 
//     for i = 0; i <= 1; i++ {
//         O[i] = make([]int32, K+1)
//         O[i][0] = 0
//     }
// 
//     // reset item-in-use table
//     for j = 0; j < N; j++ {
//         x[j] = 0
//     }
//     O[0][0] = 0
//     O[1][0] = 0
//     // O(k,j) denotes the optimal solution to the
//     // knapsack problem with capacity k and
//     // items [1..j]
// 
//     lastPackedSize = dumpToFile(file, O[0])
//     offsets[0] = position
//     sizes[0] = int(lastPackedSize)
//     position += lastPackedSize
// 
//     // for all items
//     for j = 1; j <= N; j++ {
//         // for all capacities
//         for k = 1; k <= K; k++ {
//             if w[j-1] <= k {
//                 O[1][k] = max(O[0][k], v[j-1] + O[0][k-w[j-1]])
//             } else {
//                 O[1][k] = O[0][k]
//             }
//         }
// 
//         // dump to disk and save offset
//         lastPackedSize = dumpToFile(file, O[1])
//         offsets[j] = position
//         sizes[j] = int(lastPackedSize)
//         position += lastPackedSize
// 
//         for i := 0; i <= int(K); i++ {
//             O[0][i] = O[1][i]
//         }
//     }
// 
//     file.Sync()
//     fmt.Println(O[1][K], 1)
// 
//     // restore best set of items
//     k = K
//     file.Seek(offsets[N], 0)
//     loadFromFile(file, sizes[N], &O[1])
//     for i = N; i > 0; i-- {
//         // preload first (previous) column
//         file.Seek(offsets[i-1], 0)
//         loadFromFile(file, sizes[i-1], &O[0])
// 
//         if O[1][k] != O[0][k] {
//             x[i-1] = 1
//             k -= w[i-1]
//         }
//     }
// 
//     // print best set
//     for i = 0; i < N; i++ {
//         if i == N-1 {
//             fmt.Printf("%d", x[i])
//         } else {
//             fmt.Printf("%d ", x[i])
//         }
//     }
//     fmt.Printf("\n")
// 
//     file.Close()
// }
