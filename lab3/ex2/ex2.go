
package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type LamportClock struct {
	counter int
}

func (l *LamportClock) Tick() {
	l.counter++
}

func (l *LamportClock) Send() int {
	l.Tick()
	return l.counter
}

func (l *LamportClock) Receive(ts int) {
	if ts > l.counter {
		l.counter = ts
	}
	l.Tick()
}

type Message struct {
	sender int
	ts     int
}

type Event struct {
	process          int
	internal_counter int
	kind             string
}

func less(a, b Event) bool {
	if a.internal_counter != b.internal_counter {
		return a.internal_counter < b.internal_counter
	}
	return a.process < b.process
}

func process(id int, in <-chan Message, out chan<- Message, events chan<- Event, wg *sync.WaitGroup) {
	defer wg.Done()
	clock := LamportClock{counter: 0}

	randomInternalEvents := rand.Intn(3)
	for i := 0; i < randomInternalEvents; i++ {
		clock.Tick()
		events <- Event{process: id, internal_counter: clock.counter, kind: "internal"}
	}

	ts := clock.Send()
	events <- Event{process: id, internal_counter: ts, kind: "send"}
	out <- Message{sender: id, ts: ts}

	randomInternalEvents = rand.Intn(3)
	for i := 0; i < randomInternalEvents; i++ {
		clock.Tick()
		events <- Event{process: id, internal_counter: clock.counter, kind: "internal"}
	}

	msg := <-in
	clock.Receive(msg.ts)
	events <- Event{process: id, internal_counter: clock.counter, kind: "recv"}
}

func main() {
	ch12 := make(chan Message, 1)
	ch23 := make(chan Message, 1)
	ch31 := make(chan Message, 1)
	eventsChan := make(chan Event)

	var wg sync.WaitGroup
	wg.Add(3)
	go process(1, ch31, ch12, eventsChan, &wg)
	go process(2, ch12, ch23, eventsChan, &wg)
	go process(3, ch23, ch31, eventsChan, &wg)

	go func() {
		wg.Wait()
		close(eventsChan)
	}()

	var eventList []Event
	for e := range eventsChan {
		eventList = append(eventList, e)
	}

	sort.Slice(eventList, func(i, j int) bool {
		return less(eventList[i], eventList[j])
	})

	fmt.Println("Event. ordonate conform Lamport:")
	for _, e := range eventList {
		fmt.Printf("Proces %d | T. Lamport: %d | Tip: %s\n", e.process, e.internal_counter, e.kind)
	}
}
