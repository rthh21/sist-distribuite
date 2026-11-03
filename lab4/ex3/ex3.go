
package main

import (
	"fmt"
	"sync"
)

type Msg struct {
	Value float64
}

func send(send_val float64, dreapta chan<- Msg) {
	dreapta <- Msg{Value: send_val}
}

func recv(stanga <-chan Msg) Msg {
	return <-stanga
}

func node(vot float64, stanga []<-chan Msg, dreapta []chan<- Msg, status *float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, out := range dreapta {
		send(vot, out)
	}

	suma := vot
	for _, in := range stanga {
		msg := recv(in)
		suma += msg.Value
	}

	*status = suma / float64(len(stanga)+1)
}

func MajorityFlooding(votes []int, adjMatrix [][]int) bool {
	n := len(votes)
	if n == 0 {
		return false
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

	valori := make([]float64, n)
	for i := 0; i < n; i++ {
		valori[i] = float64(votes[i])
	}

	for {
		status := make([]float64, n)
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
			go node(valori[i], stanga, dreapta, &status[i], &wg)
		}

		wg.Wait()

		minVal, maxVal := status[0], status[0]
		suma := 0.0
		for i := 0; i < n; i++ {
			valori[i] = status[i]
			suma += status[i]
			if status[i] < minVal {
				minVal = status[i]
			}
			if status[i] > maxVal {
				maxVal = status[i]
			}
		}

		if maxVal-minVal < 0.000001 {
			return suma/float64(n) > 0.5
		}
	}
}

func main() {
	voturi := []int{1, 1, 0}
	adjMatrix := [][]int{
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}
	fmt.Println("Majoritate obtinuta:", MajorityFlooding(voturi, adjMatrix))
}
