package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type Info struct {
	ID    int
	Value int
}

type Memory map[Info]struct{}

func diametru(adj [][]int) int {
	n := len(adj)
	diam := 0
	for src := 0; src < n; src++ {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		dist[src] = 0
		queue := []int{src}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
		for _, d := range dist {
			if d == -1 {
				return -1
			}
			if d > diam {
				diam = d
			}
		}
	}
	return diam
}

func ruleazaNod(id int, token int, in <-chan []Memory, out chan<- Memory) {
	mem := make(Memory)
	mem[Info{ID: id, Value: token}] = struct{}{}

	for {
		out <- mem

		mesaje, ok := <-in
		if !ok {
			return
		}

		for i := 0; i < len(mesaje); i++ {
			for info := range mesaje[i] {
				mem[info] = struct{}{}
			}
		}
	}
}

func floodSet(n int, s int, diam int, tokens []int, adj [][]int) (int, Memory) {
	in := make([]chan []Memory, n)
	out := make([]chan Memory, n)
	crashed := make([]bool, n)

	for i := 0; i < n; i++ {
		in[i] = make(chan []Memory)
		out[i] = make(chan Memory)
		crashed[i] = false

		go ruleazaNod(i, tokens[i], in[i], out[i])
	}

	t := 0
	crashedCount := 0

	for {
		if crashedCount < s {
			if rand.Intn(2) == 1 {
				var active []int
				for i := 0; i < n; i++ {
					if !crashed[i] {
						active = append(active, i)
					}
				}
				if len(active) > 1 {
					victima := active[rand.Intn(len(active))]
					crashed[victima] = true
					crashedCount++
					<-out[victima]
					close(in[victima])
				}
			}
		}

		memorii := make([]Memory, n)
		for i := 0; i < n; i++ {
			if !crashed[i] {
				memorii[i] = <-out[i]
			}
		}

		inbox := make([][]Memory, n)
		for i := 0; i < n; i++ {
			if !crashed[i] {
				for _, v := range adj[i] {
					if !crashed[v] {
						inbox[v] = append(inbox[v], memorii[i])
					}
				}
			}
		}

		for i := 0; i < n; i++ {
			if !crashed[i] {
				in[i] <- inbox[i]
			}
		}

		if t > (s+1)*diam {
			break
		}
		t++
	}

	var memFinala Memory
	for i := 0; i < n; i++ {
		if !crashed[i] {
			if memFinala == nil {
				memFinala = <-out[i]
			} else {
				<-out[i]
			}
			close(in[i])
		}
	}

	return t, memFinala
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var n, s int
	fmt.Scan(&n, &s)

	tokens := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&tokens[i])
	}

	scanner := bufio.NewScanner(os.Stdin)
	adj := make([][]int, n)

	// citire liste de adiacenta
	for i := 0; i < n; {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		for _, w := range strings.Fields(line) {
			var v int
			fmt.Sscan(w, &v)
			adj[i] = append(adj[i], v)
		}
		i++
	}

	diam := diametru(adj)
	if diam == -1 {
		fmt.Println("Graful nu e tare conex") 
		return
	}

	t, memFinala := floodSet(n, s, diam, tokens, adj)

	var final []Info
	for info := range memFinala {
		final = append(final, info)
	}

	sort.Slice(final, func(i, j int) bool {
		return final[i].ID < final[j].ID
	})

	fmt.Printf("Diametrul grafului: %d\n", diam)
	fmt.Printf("Iteratia de oprire: %d\n", t)
	for i, info := range final {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Printf("(%d, %d)", info.ID, info.Value)
	}
	fmt.Println()
}
