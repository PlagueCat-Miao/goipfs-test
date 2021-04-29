package strategy

import (
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/filesys"
	"math"
)

type Strategy interface {
	PerformanceEvaluation() (remStandardDeviation float64, rem int64, use int64, maxRem int64)
	AddNode(node Node)
	AddFile(file filesys.TestFile)
	FailReport(num int) (successRate float64)
	PrintNodesUse() (nodesUse []int64)
}

//Kademlia
type Node struct {
	ID       [20]byte //160 bit=20Byte  19 is big
	Remain   int64    //MB  40*1024 40GB
	Capacity int64
}

type Base struct {
	NodeList  []*Node
	FailFiles []filesys.TestFile
}

//返回两个指标 ：1剩余容量的标准差 代表负载均衡 2系统剩余容量
func (b *Base) PerformanceEvaluation() (remStandardDeviation float64, rem int64, use int64, maxRem int64) {
	n := float64(len(b.NodeList))
	cap := int64(0)
	rem = int64(0)
	maxRem = int64(0)
	for _, node := range b.NodeList {
		cap += node.Capacity
		rem += node.Remain
		if maxRem < node.Remain {
			maxRem = node.Remain
		}
	}
	remAvge := float64(rem) / n
	s := float64(0)
	for _, node := range b.NodeList {
		s += (float64(node.Remain) - remAvge) * (float64(node.Remain) - remAvge)
	}
	remStandardDeviation = math.Sqrt(s / n)
	use = cap - rem
	return
}
func (b *Base) AddNode(node Node) {
	myNode := Node{
		ID:       node.ID,
		Remain:   node.Remain,
		Capacity: node.Capacity,
	}
	b.NodeList = append(b.NodeList, &myNode)
}
func (b *Base) FailReport(num int) (successRate float64) {
	successRate = float64(len(b.FailFiles)) / float64(num)

	return
}

func (b *Base) PrintNodesUse() []int64 {
	var nodesUse []int64
	for _, node := range b.NodeList {
		nodesUse = append(nodesUse, node.Capacity-node.Remain)
	}
	return nodesUse
}
