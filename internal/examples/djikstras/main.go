package main

import (
	"fmt"

	"github.com/dynago/dg/dict"
	"github.com/dynago/dg/set"
	"github.com/dynago/dg/list"
)

/* This is an example using dictionaries to implement Djikstra's Algorithm */

type Edge struct {
	node   string
	weight uint
}

const UNVISITED = ^uint(0) - 1

func getUInt(value interface{}) uint {
	switch value.(type) {
	case int:
		return uint(value.(int))
	case uint:
		return value.(uint)
	default:
		panic("interface{} should convert to either int or uint")
	}
	return 0
}

func minNode(unvisited dict.DictInterface) (string, uint) {
	minN := ""
	minD := UNVISITED + 1
	for node := range unvisited.Iterate() {
		d, _ := unvisited.Get(node)
		distance := getUInt(d)
		if distance < minD {
			minN = node.(string)
			minD = distance
		}
	}
	return minN, minD
}

func getPath(start string, end string, parents dict.DictInterface) list.ListInterface {
	curr := end
	path, _ := list.MakeListFromValues(end)
	for curr != start {
		c, _ := parents.Get(curr)
		if c == nil {
			return nil
		}
		curr = c.(string)
		path.Append(curr)
	}
	path, _ = path.Reverse()
	return path
}

/*
Input:
      3
  a ----- b
  |     /
1 |   / 1
  | /
  c       d

Expected Output:
a -> a: distance = 0,   path = [a]
a -> b: distance = 2,   path = [a c b]
a -> c: distance = 1,   path = [a c]
a -> d: distance = inf, path = []
*/


func main() {
	// Edges
	a, _ := set.MakeSetFromValues(Edge{"b", 3}, Edge{"c", 1})
	b, _ := set.MakeSetFromValues(Edge{"a", 3}, Edge{"c", 1})
	c, _ := set.MakeSetFromValues(Edge{"a", 1}, Edge{"b", 1})
	d, _ := set.MakeSetFromValues()

	nodes := []interface{}{"a", "b", "c", "d"}
	edges := []interface{}{a, b, c, d}

	graph, _ := dict.MakeDictFromKeyValues(nodes, edges)

	start := "a"

	// unvisited nodes -> distances
	unvisited, _ := dict.MakeDictFromKeyValues(nodes, []interface{}{UNVISITED, UNVISITED, UNVISITED, UNVISITED})
	unvisited.Set(start, 0) // set vertex to distance = 0

	visited, _ := dict.MakeDict()
	parents, _ := dict.MakeDict()

	for unvisited.Length() > 0 {
		minN, minD := minNode(unvisited)

		neighbours, _ := graph.Get(minN)
		neighbourSet := neighbours.(*set.Set)
		for n := range neighbourSet.Iterate() {
			neighbour := n.(Edge)
			if contains, _ := visited.Contains(neighbour.node); !contains {
				newD := minD + neighbour.weight
				if neighD, _ := unvisited.Get(neighbour.node); newD < getUInt(neighD) {
					unvisited.Set(neighbour.node, newD)
                    parents.Set(neighbour.node, minN)
				}
			}
		}

		visited.Set(minN, minD)
        unvisited.Remove(minN)
	}

	for _, end := range nodes {
		fmt.Printf("%s -> %s: ", start, end)
		distance, _ := visited.Get(end.(string)) 
		if distance == UNVISITED {
			fmt.Println("distance = inf,\tpath = []")
		} else {
			path := getPath(start, end.(string), parents)
			fmt.Printf("distance = %d,\tpath = %s\n", getUInt(distance), path)
		}
	}
}
