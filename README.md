
Example usage:

```golang
package main

import (
        "fmt"
        "math/rand"
        "time"

        "github.com/t0mk/ratehistogram"
)

func main() {
        y := `
h1:
  edges: [1.0, 2, 3, 4, 5]
  secs: 5
h2:
  edges: [5, 6, 7, 8, 9, 10]
  secs: 2
`
        fmt.Println("vim-go")
        hm, err := ratehistogram.NewHMapFromYAML([]byte(y))
        if err != nil {
                panic(err)
        }
        feedTicker := time.NewTicker(time.Millisecond * 5)
        go func() {
                for range feedTicker.C {
                        hm["h1"].Record(rand.NormFloat64()*1 + 2.5)
                        hm["h2"].Record(rand.NormFloat64()*1 + 7.5)
                }
        }()
        ticker := time.NewTicker(time.Millisecond * 500)
        for t := range ticker.C {
                fmt.Println("Tick at", t)
                fmt.Println("h1:", hm["h1"].Observe())
                fmt.Println("h2:", hm["h2"].Observe())
        }
}

func oldmain() {
        fmt.Println("vim-go")
        nhc := ratehistogram.Conf{Edges: []float64{1, 2, 3, 4, 5}, Secs: 5}
        rh, err := ratehistogram.NewRateHistogram(nhc)
        if err != nil {
                panic(err)
        }
        feedTicker := time.NewTicker(time.Millisecond * 5)
        go func() {
                for range feedTicker.C {
                        rh.Record(rand.NormFloat64()*1 + 2.5)
                }
        }()
        ticker := time.NewTicker(time.Millisecond * 500)
        for t := range ticker.C {
                fmt.Println("Tick at", t)
                fmt.Println(rh.Observe())
        }
}

```
