
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Clock struct {
	offset time.Duration
}

func (c *Clock) Now() time.Time {
	return time.Now().Add(c.offset)
}

func estimate(t1, t2, t3, t4 time.Time) (theta, delta time.Duration) {
	theta = (t2.Sub(t1) + t3.Sub(t4)) / 2
	delta = t4.Sub(t1) - t3.Sub(t2)
	return theta, delta
}

func measure(c_clock, s_clock *Clock) (theta, delta time.Duration) {
	t1 := c_clock.Now()
	time.Sleep(3*time.Millisecond + time.Duration(rand.Intn(10))*time.Millisecond)

	t2 := s_clock.Now()
	time.Sleep(5*time.Millisecond + time.Duration(rand.Intn(20))*time.Millisecond)
	t3 := s_clock.Now()

	time.Sleep(3*time.Millisecond + time.Duration(rand.Intn(10))*time.Millisecond)
	t4 := c_clock.Now()

	return estimate(t1, t2, t3, t4)
}

func main() {
	client := &Clock{offset: 500 * time.Millisecond}
	server := &Clock{offset: 0}

	fmt.Println("Offset init:", client.offset)

	thetaA, deltaA := measure(client, server)
	fmt.Printf("a: Theta: %v, Delta: %v, Offset nou: %v\n", thetaA, deltaA, client.offset+thetaA)

	var minTheta, minDelta time.Duration
	for i := 0; i < 10; i++ {
		theta, delta := measure(client, server)
		if i == 0 || theta < minTheta {
			minTheta = theta
		}
		if i == 0 || delta < minDelta {
			minDelta = delta
		}
	}

	fmt.Printf("b: Theta minim: %v, Delta minim: %v\n", minTheta, minDelta)
	fmt.Printf("   Offset corectat (Theta minim): %v\n", client.offset+minTheta)
	fmt.Printf("(b - a): Theta: %v, Delta: %v\n", minTheta-thetaA, minDelta-deltaA)
}
