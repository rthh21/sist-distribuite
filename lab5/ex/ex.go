package main

import (
    "fmt"
    "math"
    "os"
)

type Robot struct {
    X,Y float64
}

type graf struct {
    r []Robot
}

func init_graf() ([]Robot,error){
    file, err := os.Open("in.txt")
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var n int 
    fmt.Fscan(file, &n)
    Robots := make([]Robot, n)
    for i := 0; i < n; i++ {
        var r Robot
        fmt.Fscan(file, &r.X, &r.Y)
        Robots[i] = r
    }
    return Robots, nil
}   

func dist(a, b Robot) float64 {
    return math.Hypot(a.X-b.X, a.Y-b.Y)
}

// func parcurge(n int, r Robot) graf{
//  // for fiecare nod
//  // calculam new_position = (pos + sum(vecini)) / (1 + len(vecini))
//  // mutam nod la new_position in graf
//  // return graf nou
// }

func parcurge(n int, g graf) graf {
    graf_nou := graf{r: make([]Robot, n)}
    raza := 20.0

    for i := 0; i < n; i++ {
        sumX := g.r[i].X
        sumY := g.r[i].Y
        nr := 1

        for j := 0; j < n; j++ {
            if i != j && dist(g.r[i], g.r[j]) < raza {
                sumX += g.r[j].X
                sumY += g.r[j].Y
                nr++
            }
        }

        graf_nou.r[i].X = sumX / float64(nr)
        graf_nou.r[i].Y = sumY / float64(nr)
    }
    return graf_nou
}

func comprima(g graf) graf {
    var active []Robot
    for _, r := range g.r {
        isDuplicate := false
        for _, a := range active {
            if dist(r, a) < 1e-5 { 
                isDuplicate = true
                break
            }
        }
        if !isDuplicate {
            active = append(active, r)
        }
    }
    return graf{r: active}
}

func main() {
    Roboti, err := init_graf()
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Poz initiale:")
    for i, r := range Roboti {
        fmt.Printf("%d: x = %.2f, y = %.2f\n", i+1, r.X, r.Y)
    }

    g := graf{r: Roboti}
    max_iter := 1000 

    for iter := 0; iter < max_iter; iter++ {
        g = comprima(g)
        n := len(g.r)
        
        if n <= 1 {
            fmt.Printf("\nRendezvous in %d iteratii.\n", iter)
            break
        }
        
        g = parcurge(n, g)
    }

    fmt.Println("\nPoz finale:")
    for i, r := range g.r {
        fmt.Printf("x = %.4f, y = %.4f\n", r.X, r.Y)
    }
    // ------------------------------------

    // test_graf := graf{r: make([]Robot, n)}
    
    // var r1, r2 Robot
    // r1.X = 1.0
    // r2.X = 2.0
    // r1.Y = 2.0
    // r2.Y = -5.0
    // fmt.Println(dist(r1,r2))
}