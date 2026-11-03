package main

import (
	"fmt"
	"sync"
)

type Msg struct {
	Value int
	Done  bool
}

func send(send_id int, dreapta chan<- Msg) {
	if send_id != -1 {
		dreapta <- Msg{Value: send_id, Done: false}
	}
}

func recv(stanga <-chan Msg) Msg {
	return <-stanga
}

func node(id int, stanga <-chan Msg, dreapta chan<- Msg, status *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	send_id := id
	*status = false
	sent_done := false

	for {
		if !sent_done {
			send(send_id, dreapta)
		}

		msg := recv(stanga)

		if msg.Done {
			if !*status {
				dreapta <- msg
			}
			return
		}

		recv_id := msg.Value

		if recv_id > id {
			send_id = recv_id
		} else if recv_id == id {
			*status = true
			dreapta <- Msg{Value: id, Done: true}
			sent_done = true
			send_id = -1
		} else if recv_id < id {
			send_id = -1
		}
	}
}

func RingFlooding(ids []int) int {
	n := len(ids)
	if n == 0 {
		return -1
	}

	chans := make([]chan Msg, n)
	for i := 0; i < n; i++ {
		chans[i] = make(chan Msg, n)
	}

	status := make([]bool, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go node(ids[i], chans[(i-1+n)%n], chans[i], &status[i], &wg)
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
	ids := []int{7, 4, 9, 2, 5, 33, 10, 1}
	fmt.Println(RingFlooding(ids))
}
