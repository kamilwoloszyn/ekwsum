# ekwsum
Simple sum control generator for ekw

## Quick start

Importujemy bibiotekę

`go get github.com/kamilwoloszyn/ekwsum`

I możemy jej używać: 


```
package main

import (
	"fmt"
	"log"

	"github.com/kamilwoloszyn/ekwsum"
)

func main() {
	ekwNum, err := ekwsum.NewEkwNumber("<Your ekw number>")
	if err != nil {
		log.Fatal(err)
	}
	if err := ekwNum.Validate(); err != nil {
		log.Fatal(err)
	}
	sum := ekwNum.SumControl()
	fmt.Println(sum)
}


```

## Docs

Available [here](https://pkg.go.dev/github.com/kamilwoloszyn/ekwsum)