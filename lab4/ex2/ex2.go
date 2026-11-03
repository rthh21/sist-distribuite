
package main

import (
	"fmt"
	"sync"
)

type Msg struct {
	Value int
}

func send(send_id int, dreapta chan<- Msg) {
	dreapta <- Msg{Value: send_id}
}

func recv(stanga <-chan Msg) Msg {
	return <-stanga
}

func node(id int, stanga []<-chan Msg, dreapta []chan<- Msg, status *bool, runde int, wg *sync.WaitGroup) {
	defer wg.Done()

	max_id := id
	*status = false

	for r := 0; r < runde; r++ {
		for _, out := range dreapta {
			send(max_id, out)
		}

		for _, in := range stanga {
			msg := recv(in)
			if msg.Value > max_id {
				max_id = msg.Value
			}
		}
	}

	if max_id == id {
		*status = true
	}
}

func GraphFlooding(ids []int, adjMatrix [][]int) int {
	n := len(ids)
	if n == 0 {
		return -1
	}

	chans := make([][]chan Msg, n)
	for i := 0; i < n; i++ {
		chans[i] = make([]chan Msg, n)
		for j := 0; j < n; j++ {
			if adjMatrix[i][j] == 1 {
				chans[i][j] = make(chan Msg, 1)
			}
		}
	}

	status := make([]bool, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		var stanga []<-chan Msg
		var dreapta []chan<- Msg

		for j := 0; j < n; j++ {
			if adjMatrix[j][i] == 1 {
				stanga = append(stanga, chans[j][i])
			}
			if adjMatrix[i][j] == 1 {
				dreapta = append(dreapta, chans[i][j])
			}
		}

		wg.Add(1)
		go node(ids[i], stanga, dreapta, &status[i], n-1, &wg)
	}

	wg.Wait()

	for i := 0; i < n; i++ {
		if status[i] {
			return ids[i]
		}
	}
	return -1
}

func main() {
	ids := []int{10, 42, 5}
	adjMatrix := [][]int{
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}
	fmt.Println("Lider:", GraphFlooding(ids, adjMatrix))
}
