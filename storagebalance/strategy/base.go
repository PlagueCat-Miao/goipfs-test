package strategy

import "math"

type Strategy interface {
	PerformanceEvaluation ()(remStandardDeviation float64,rem int64)
	AddNode(node Node)
	AddFile(file *File)
}
//Kademlia
type Node struct {
	ID [20]byte  //160 bit=20Byte  19 is big
	Remain int64 //KB
	Capacity int64
}
type File struct{
	ID [20]byte //对于长度小于2^64位的消息,SHA1会产生一个160位的消息摘要
	Size int64 // KB
}


type Base struct {
	NodeList []*Node
}
//返回两个指标 ：1剩余容量的标准差 代表负载均衡 2系统剩余容量
func (b *Base) PerformanceEvaluation ()(remStandardDeviation float64,rem int64){
	n := float64(len(b.NodeList))
	cap := int64(0)
	rem = int64(0)
	for _,node := range b.NodeList{
		cap+= node.Capacity
		rem+= node.Remain
	}
	remAvge := float64(rem) / n
	s := float64(0)
	for _,node := range b.NodeList{
		s += (float64(node.Remain) - remAvge) * (float64(node.Remain) - remAvge)
	}
	remStandardDeviation = math.Sqrt(s/n)
	return
}
func (b *Base) AddNode (node Node){
	 myNode :=Node{
		 ID:       node.ID,
		 Remain:   node.Remain,
		 Capacity: node.Capacity,
	 }
	b.NodeList = append(b.NodeList,&myNode)
}
