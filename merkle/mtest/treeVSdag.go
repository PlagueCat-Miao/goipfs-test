package mtest

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mr-tron/base58"
	"io"
	"io/ioutil"
	"math"
	"os"
	"time"
)

const block_byte_size = int64(256 * 1024)
const sha256Ans_size = 32
const Mulithash_size = 47

//const XuChang_Bayi_Road = "E:\\lab511-ipfs实验\\八一路与劳动路口反卡北侧\\"
//const XuChang_Huatuo_Road = "E:\\lab511-ipfs实验\\华佗路与五一路口反卡东侧\\"
//const XuChang_TiBao_Road = "E:\\lab511-ipfs实验\\天宝路与文峰路口反卡北侧\\"

const XuChang_Bayi_Road = "/home/hellcat/plague/data-play/lab511-vidoe-test/八一路与劳动路口反卡北侧/"
const XuChang_Huatuo_Road = "/home/hellcat/plague/data-play/lab511-vidoe-test/华佗路与五一路口反卡东侧/"
const XuChang_TiBao_Road = "/home/hellcat/plague/data-play/lab511-vidoe-test/天宝路与文峰路口反卡北侧/"

var pathMap = map[string][]string{
	"A组": {XuChang_Bayi_Road},
	"B组": {XuChang_Huatuo_Road},
	"C组": {XuChang_TiBao_Road},

	"A+B组": {XuChang_Bayi_Road, XuChang_Huatuo_Road},
	"B+C组": {XuChang_Huatuo_Road, XuChang_TiBao_Road},
	"A+C组": {XuChang_Bayi_Road, XuChang_TiBao_Road},

	"A+B+C组": {XuChang_Bayi_Road, XuChang_Huatuo_Road, XuChang_TiBao_Road},
}

var testBlock []byte

var merkleTreeByte int64
var merkleDagByte int64
var treeSumTime time.Duration
var dagSumTime time.Duration

var TreeTimeList []time.Duration

var dagMap map[string]bool
var reuse int64

func InitTest() {
	testBlock = []byte("asdfghhklqpwoeirutyvbxnz,zxcnm")
	merkleTreeByte = 0
	merkleDagByte = 0
	treeSumTime = 0
	dagSumTime = 0
	TreeTimeList = nil

	dagMap = make(map[string]bool)
	reuse = 0
}

func TestSpeedAndSpace(label string) (int64, int64, int64, int64, int64, error) {
	//初始化
	InitTest()

	//处理
	for i, path := range pathMap[label] {
		err := AddDirByFile(path, processBlock, i, len(pathMap[label]))
		if err != nil {
			return -1, -1, -1, -1, -1, fmt.Errorf("testSpeedandSpace Err:%v, Path:%v, label:%v", err, path, label)
		}
	}

	//打印结果
	fmt.Printf("\n\n**test-%s report**\n", label)
	fmt.Printf("merkleTree have %+v B, need %+v ms\n", merkleTreeByte, treeSumTime.Milliseconds())
	fmt.Printf("merkleDAG  have %+v B, need %+v ms, reuse:%+v block\n\n", merkleDagByte, dagSumTime.Milliseconds(), reuse)

	return merkleTreeByte / 1024, treeSumTime.Milliseconds(), merkleDagByte / 1024, dagSumTime.Milliseconds(), reuse, nil
}

//处理文件夹
func AddDirByFile(pathname string, hookfn func([]byte), listNum, listLen int) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return fmt.Errorf("ReadDir err:%v,pathname:%v", err, pathname)
	}

	progress := len(rd)
	for p, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("*[%s] Skip ! *\n ", pathname+"\\"+fi.Name())
			continue
		} else {
			// 处理单个文件
			cnt, err := AddFileByBlock(pathname+fi.Name(), hookfn)
			if err != nil {
				return fmt.Errorf("ReadAFileByBlock err:%v", err)
			}
			i := math.Ceil(math.Log2(float64(cnt)))
			merkleTreeNodeNum := int64(math.Exp2(i+1) - 1)
			merkleTreeByte += merkleTreeNodeNum * (sha256Ans_size)
			treeTime, err := AverageTreeTime(cnt)
			if err != nil {
				return fmt.Errorf("AverageTreeTime err:%v", err)
			}
			treeSumTime += treeTime * time.Duration(merkleTreeNodeNum)
			fmt.Printf("rate of advance : %v %% :%v \r", p*100/(progress*listLen)+listNum*100/listLen, fi.Name())
			//fmt.Printf("merkleTree -> 叶子数:%+v 树高:%+v 节点数:%+v 平局速度:%v 微秒/块\n\n", cnt, i, merkleTreeNodeNum,treeTime.Microseconds())
		}
	}
	return nil
}

//处理文件
func AddFileByBlock(filePth string, hookfn func([]byte)) (int64, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return 0, fmt.Errorf("Open err:%v ,filePath: %v ", err, filePth)
	}
	defer f.Close()

	buf := make([]byte, block_byte_size) //一次读取多少个字节
	bfRd := bufio.NewReader(f)
	cnt := int64(0)
	for {
		n, err := bfRd.Read(buf)
		cnt++
		hookfn(buf[:n]) // n 是成功读取字节数

		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return cnt, nil
			}
			return cnt, fmt.Errorf("Read err:%v ,filePath: %v ", err, filePth)
		}

	}
	return cnt, nil
}

//处理一个数据块
func processBlock(block []byte) {
	startTime := time.Now()
	hash := GetSHA256HashCode(block)
	TreeTimeList = append(TreeTimeList, time.Now().Sub(startTime))
	base58.Encode([]byte(hash))
	dagTime := time.Now().Sub(startTime)
	dagTime += 4

	if _, ok := dagMap[hash]; ok != true {
		merkleDagByte += Mulithash_size
		dagSumTime += dagTime
		dagMap[hash] = true
	} else {
		reuse++
	}
}

//SHA256生成哈希值
func GetSHA256HashCode(message []byte) string {
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New()
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode

}

//merkle Tree 计算平均时间
func AverageTreeTime(cnt int64) (time.Duration, error) { //Milliseconds() //ms
	ans := time.Duration(0)
	if cnt == 0 || TreeTimeList == nil {
		return -1, fmt.Errorf("TreeTimeList is empty")
	}
	for _, t := range TreeTimeList {
		ans += t
	}
	ans = ans / time.Duration(cnt)
	TreeTimeList = nil
	return ans, nil
}
