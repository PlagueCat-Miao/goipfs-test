package main

import (
	"fmt"
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/strategy"
	"math"
)

const max_nodes = 20

func InitNodes(strategies ...strategy.Strategy) {

	topByte := byte(1)
	step := byte(math.Floor(math.Exp2(8) / float64(max_nodes)))

	for i := 0; i < max_nodes; i++ {
		node := strategy.Node{
			ID:       [20]byte{19: topByte},
			Remain:   1024,
			Capacity: 1024,
		}
		topByte += step
		for _,stg:= range strategies{
			stg.AddNode(node)
		}

	}

}

func main() {
	kad := strategy.NewKademlia()
	mrss := strategy.NewMRSS()
	InitNodes(kad,mrss)
	file := &strategy.File{
		ID:   [20]byte{19: 34 },
		Size: 30,
	}
	kad.AddFile(file)
	mrss.AddFile(file)
	remSD,rem:=kad.PerformanceEvaluation()
	fmt.Printf("kad : %.4f %+v\n", remSD,rem)
	remSD,rem=mrss.PerformanceEvaluation()
	fmt.Printf("mrss : %.4f %+v\n", remSD,rem)
}
