package main

import (
	"fmt"

	"github.com/greenthepear/egriden"
)

func main() {
	g := egriden.NewEgridenGame()
	g.CreateGridLayerOnTop("Test", 16, 20, 20, egriden.Dense, 0, 0)
	fmt.Print(g.GridLayer(0))
}
