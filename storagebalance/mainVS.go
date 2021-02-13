package main

import (
	"fmt"
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/filesys"
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/strategy"
	"math"
)

const max_nodes = 100
const redundant = 5
const space = 45 * 1024 //30*1024 有奇效
func InitNodes(strategies map[string]strategy.Strategy) {

	topByte := byte(1)
	step := byte(math.Floor(math.Exp2(8) / float64(max_nodes)))

	for i := 0; i < max_nodes; i++ {
		node := strategy.Node{
			ID:       [20]byte{19: topByte},
			Remain:   space,
			Capacity: space,
		}
		topByte += step
		for _, stg := range strategies {
			stg.AddNode(node)
		}

	}

}

//func main(){
//	filesys.GenerateFileInfo(1000)
//}

func main() {
	kad := strategy.NewKademlia(redundant)
	mrss := strategy.NewMRSS(redundant)
	asce := strategy.NewAsce(redundant)
	randm := strategy.NewRandom(redundant)
	opStrategy := make(map[string]strategy.Strategy)
	opStrategy["kad"] = kad
	opStrategy["mrss"] = mrss
	opStrategy["asce"] = asce
	opStrategy["randm"] = randm
	InitNodes(opStrategy)
	files, _ := filesys.ReadFileInfo("FileInfo-1613131121.json")

	n := len(files)
	sum := int64(0)
	for i, file := range files {
		fmt.Printf("progress : %+v %%\r", 100*i/n)

		//file.Size = filesys.UnifySize
		No := i % 4
		switch No {
		case 0:
			file.Size = filesys.R_bigSize
			break
		case 1:
			file.Size = filesys.R_midSize
			break
		case 2:
			file.Size = filesys.R_smallSize
			break
		case 3:
			file.Size = filesys.R_unifySize
			break
		}

		sum += file.Size * redundant
		for _, op := range opStrategy {
			op.AddFile(file)
		}
	}
	fmt.Printf("\n\nprint report:\n")
	fmt.Printf("数据可用性: %+vMB\n", sum)
	for name, op := range opStrategy {
		remSD, rem, use, maxRem := op.PerformanceEvaluation()
		fmt.Printf("%5v - 负载标准差: %.4f remain: %+vMB use: %+vMB\n", name, remSD, rem, use)
		fmt.Printf("        最大节点剩余空间:%v 失败率: %+.4f\n", maxRem, op.FailReport(n))
	}

}
