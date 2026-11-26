package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Message struct {
	sender_id int
	round     int
	value     int
}

type Node struct {
	id           int
	x            int
	s            int
	n            int
	is_byzantine bool
	inbox        chan Message
}

func bcast_state(nd *Node, nodes []*Node) {
	for i := 0; i < nd.n; i++ {
		nodes[i].inbox <- Message{
			sender_id: nd.id,
			round:     1,
			value:     nd.x,
		}
	}
}

func malicious_bcast_state(nd *Node, nodes []*Node) {
	for i := 0; i < nd.n; i++ {
		nodes[i].inbox <- Message{
			sender_id: nd.id,
			round:     1,
			value:     rand.Intn(2),
		}
	}
}

func bcast_tie(nd *Node, nodes []*Node, maj int) {
	for i := 0; i < nd.n; i++ {
		nodes[i].inbox <- Message{
			sender_id: nd.id,
			round:     2,
			value:     maj,
		}
	}
}

func malicious_bcast_tie(nd *Node, nodes []*Node) {
	for i := 0; i < nd.n; i++ {
		nodes[i].inbox <- Message{
			sender_id: nd.id,
			round:     2,
			value:     rand.Intn(2),
		}
	}
}

func run_node(nd *Node, nodes []*Node, done chan<- bool, next <-chan bool) {
	for t := 0; ; t++ {
		if nd.is_byzantine {
			malicious_bcast_state(nd, nodes)
		} else {
			bcast_state(nd, nodes)
		}

		c0, c1 := 0, 0
		for i := 0; i < nd.n; i++ {
			msg := <-nd.inbox
			if msg.value == 1 {
				c1++
			} else {
				c0++
			}
		}

		maj := 0
		mult := c0
		if c1 >= nd.n/2+1 {
			maj = 1
			mult = c1
		}

		done <- true
		<-next

		leader_id := t % nd.n
		if nd.id == leader_id {
			if nd.is_byzantine {
				malicious_bcast_tie(nd, nodes)
			} else {
				bcast_tie(nd, nodes, maj)
			}
		}

		tie_msg := <-nd.inbox
		if mult > nd.n/2+nd.s {
			nd.x = maj
		} else {
			nd.x = tie_msg.value
		}

		done <- true
		if t >= nd.s+1 {
			return
		}
		<-next
	}
}

func main() {
	var n, s int
	if _, err := fmt.Scan(&n, &s); err != nil {
		fmt.Fprintln(os.Stderr, "Err: citire n, s:", err)
		return
	}

	if s < 0 || n <= 4*s {
		fmt.Fprintln(os.Stderr, "Err: Trebuie n > 4s si s >= 0")
		return
	}

	vals := make([]int, 0, n)
	for len(vals) < n {
		var tok string
		if _, err := fmt.Scan(&tok); err != nil {
			fmt.Fprintln(os.Stderr, "Err: citire valori initiale:", err)
			return
		}
		for _, ch := range strings.TrimSpace(tok) {
			if ch == '0' || ch == '1' {
				vals = append(vals, int(ch-'0'))
				if len(vals) == n {
					break
				}
			}
		}
	}

	nodes := make([]*Node, n)
	next := make([]chan bool, n)
	done := make(chan bool, n)

	for i := 0; i < n; i++ {
		nodes[i] = &Node{
			id:           i,
			x:            vals[i],
			s:            s,
			n:            n,
			is_byzantine: false,
			inbox:        make(chan Message, n),
		}
		next[i] = make(chan bool)
	}

	byz_cnt := 0
	for byz_cnt < s {
		rand_id := rand.Intn(n)
		if !nodes[rand_id].is_byzantine {
			nodes[rand_id].is_byzantine = true
			byz_cnt++
		}
	}

	for i := 0; i < n; i++ {
		go run_node(nodes[i], nodes, done, next[i])
	}

	for t := 0; t <= s+1; t++ {
		for i := 0; i < n; i++ {
			<-done
		}
		for i := 0; i < n; i++ {
			next[i] <- true
		}

		for i := 0; i < n; i++ {
			<-done
		}
		if t < s+1 {
			for i := 0; i < n; i++ {
				next[i] <- true
			}
		}
	}

	for i := 0; i < n; i++ {
		if !nodes[i].is_byzantine {
			fmt.Println(nodes[i].x)
			break
		}
	}
}