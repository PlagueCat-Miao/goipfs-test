package strategy

import "sort"

type MRSS struct {
	Base
	Redundant int  // 这里 k 是为平衡系统性能和网络负载而设置的一个常数，但必须是偶数，比如 k = 20。在 BitTorrent 的实现中，取值为 k = 8
}

func NewMRSS()*MRSS{
	var mrss MRSS
	mrss.Redundant= 4
	mrss.NodeList = make([]*Node,0)
	return &mrss
}

func (m *MRSS) AddFile (file *File){
	sort.Slice(m.NodeList, func(i, j int) bool {
		return m.NodeList[i].Remain > m.NodeList[j].Remain
	})
	n:= len(m.NodeList)
	savenum:=0
	for i:=0 ; i< n && savenum < m.Redundant;i++ {
		if m.NodeList[i].Remain > file.Size {
			m.NodeList[i].Remain -=  file.Size
			savenum++
		}
	}
}