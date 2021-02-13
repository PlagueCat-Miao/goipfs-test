package strategy

import (
	"github.com/PlagueCat-Miao/goipfs-test/storagebalance/filesys"
	"math/rand"

)
type Random struct {
	Base
	Redundant int  // 这里 k 是为平衡系统性能和网络负载而设置的一个常数，但必须是偶数，比如 k = 20。在 BitTorrent 的实现中，取值为 k = 8
	seed int64
}
func NewRandom(redundant int)*Random{
	var random Random
	random.Redundant =redundant
	random.NodeList = make([]*Node,0)
	random.seed  = 1613131121
	return &random
}

func (r *Random) AddFile (file filesys.TestFile){
	rand.Seed(r.seed)
	n := len(r.NodeList)
	list:=make([]*Node,n)
	copy(list,r.NodeList)

	savenum:=0
	for n>0 && savenum < r.Redundant {
		rnd := rand.Intn(n)
		if list[rnd].Remain >= file.Size {
			list[rnd].Remain -=  file.Size
			savenum++
		}
		if rnd+1 < n {
			list =append(list[:rnd],list[rnd+1:]...)
		}else{
			list =append(list[:rnd])
		}
		n = len(list)
	}
	if savenum!= r.Redundant{
		r.FailFiles = append(r.FailFiles,file)
	}
}
