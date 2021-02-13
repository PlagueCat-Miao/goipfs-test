package strategy

import (
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/filesys"
)

type Asce struct {
	Base
	Redundant int  // 这里 k 是为平衡系统性能和网络负载而设置的一个常数，但必须是偶数，比如 k = 20。在 BitTorrent 的实现中，取值为 k = 8
	Index int
}

func NewAsce(redundant int)*Asce{
	var asce Asce
	asce.Redundant =redundant
	asce.NodeList = make([]*Node,0)
	asce.Index=0
	return &asce
}

func (a *Asce) AddFile (file filesys.TestFile){
	n:= len(a.NodeList)
	savenum:=0
	for i:=0 ; i< n && savenum < a.Redundant;i++ {
		if a.NodeList[a.Index].Remain >= file.Size {
			a.NodeList[a.Index].Remain -=  file.Size
			savenum++
		}
		a.Index=(a.Index+1)%n
	}
	if savenum!= a.Redundant{
		a.FailFiles = append(a.FailFiles,file)
	}
}