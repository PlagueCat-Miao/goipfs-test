package strategy

import (
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/filesys"
	"sort"
)

type Kademlia struct {
	Base
	Redundant int // 这里 k 是为平衡系统性能和网络负载而设置的一个常数，但必须是偶数，比如 k = 20。在 BitTorrent 的实现中，取值为 k = 8
}

func NewKademlia(redundant int) *Kademlia {
	var kad Kademlia
	kad.Redundant = redundant
	kad.NodeList = make([]*Node, 0)
	return &kad
}

func (k *Kademlia) AddFile(file filesys.TestFile) {
	sort.Slice(k.NodeList, func(i, j int) bool { //从大到小
		return XORJudge(k.NodeList[i].ID, k.NodeList[j].ID, file.ID) //从近至远
	})
	n := len(k.NodeList)
	savenum := 0
	for i := 0; i < n && savenum < k.Redundant; i++ {
		if k.NodeList[i].Remain >= file.Size {
			k.NodeList[i].Remain -= file.Size
			savenum++
		}
	}
	if savenum != k.Redundant {
		k.FailFiles = append(k.FailFiles, file)
	}

}

func XORJudge(Ibyte [20]byte, Jbyte [20]byte, fileHash [20]byte) bool { //I更近 ture
	for i := 19; i >= 0; i-- {
		if Ibyte[i]^fileHash[i] < Jbyte[i]^fileHash[i] { // I更近
			return true
		} else if Ibyte[i]^fileHash[i] > Jbyte[i]^fileHash[i] { // J更近
			return false
		}
	}
	return false // 不太可能平局
}
