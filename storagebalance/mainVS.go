package main

import (
	"fmt"
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/filesys"
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/show"
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/strategy"

	"log"
	"math"
	"net/http"
)

const max_nodes = 100
const file_nodes = 100

const redundant = 5
const space = 45 * 1024 //30*1024 有奇效
const test_file_set = "C"
const is_generate_init = false

var strategySequence = []string{"asce", "kad", "mrss", "randm"}

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

func main() {
	if is_generate_init == true {
		filesys.GenerateFileInfo(file_nodes)
	}
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
		if test_file_set == "A" {
			file.Size = filesys.UnifySize
		} else if test_file_set == "C" {
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
		} //data B team is read FileInfo
		sum += file.Size * redundant
		for _, op := range opStrategy {
			op.AddFile(file)
		}
	}

	fmt.Printf("\n\nSet %s print report:\n", test_file_set)
	fmt.Printf("数据可用性: %+vMB\n", sum)
	strategyHeat := make(map[string][]int64)
	for _, stgName := range strategySequence {
		if op, ok := opStrategy[stgName]; ok {
			remSD, rem, use, maxRem := op.PerformanceEvaluation()
			strategyHeat[fmt.Sprintf("%s  (𝜎: %.4f )",stgName,remSD)] = op.PrintNodesUse()
			fmt.Printf("%5v - 负载标准差: %.4f remain: %+vMB use: %+vMB\n", stgName, remSD, rem, use)
			fmt.Printf("        最大节点剩余空间:%v 失败率: %+.4f\n", maxRem, op.FailReport(n))
		} else {
			fmt.Printf("%5v - 未找到\n", stgName)
		}
	}
	showHeatMap(strategyHeat)
}

func showHeatMap(strategyHeat map[string][]int64) {
	e := show.HeatmapExamples{}
	e.Experiments(strategyHeat, test_file_set)

	fs := http.FileServer(http.Dir("show/html"))
	log.Println("running server at http://localhost:8848")
	log.Fatal(http.ListenAndServe("localhost:8848", logRequest(fs)))
}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
