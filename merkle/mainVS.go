package main

import (
	"bufio"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math"
	"os"
	"crypto/sha256"
	"fmt"
)
const block_byte_size  = int64 (256*1024)
const sha256Ans_size = 32
const Mulithash_size = 47

const XuChang_Bayi_Road =  "E:\\lab511-ipfs实验\\八一路与劳动路口反卡北\\"

var merkleTreeByte int64
var merkleDagByte int64
var reuse int64
var dagMap map[string] bool

func processBlock(block []byte) {
	hash:=GetSHA256HashCode(block)
	if _,ok := dagMap[hash];ok!=true {
		merkleDagByte+=Mulithash_size
		dagMap[hash] = true
	}else{
		reuse++
	}

}

func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	progress:=len(rd)
	for p, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			// 处理单个文件
			fmt.Printf("rate of advance : %v %% :%v \n",p*100/progress,fi.Name())
			cnt,_ := ReadAFileByBlock(pathname + fi.Name(), processBlock)
			i :=math.Ceil(math.Log2(float64(cnt)))
			merkleTreeNodeNum := int64( math.Exp2(i+1) - 1)
			fmt.Printf("cnt: %+v i: %+v merkleTreeNodeNum:%+v  \n",cnt,i,merkleTreeNodeNum )
			merkleTreeByte +=  merkleTreeNodeNum*(sha256Ans_size)

		}
	}
	return err
}

func main() {
	merkleTreeByte = 0
	merkleDagByte =0
	dagMap = make(map[string]bool)

	err:=GetAllFile(XuChang_Bayi_Road)
	if err !=nil{
		fmt.Printf("\n\nErr:   %+v\n" ,err)
	}else{
		fmt.Printf("\n\nmax  int64 have 9223372036854775807 B\n")
		fmt.Printf("merkleTree have %+v B\n",merkleTreeByte)
		fmt.Printf("merkleDAG  have %+v B ,reuse:%+v Byte\n",merkleDagByte,reuse)
	}

}
//处理文件 对块进行计数
func ReadAFileByBlock(filePth string, hookfn func([]byte)) (int64,error) {
	f, err := os.Open(filePth)
	if err != nil {
		return 0,err
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
				return cnt,nil
			}
			return cnt,err
		}

	}
	return cnt,nil
}

//SHA256生成哈希值
func GetSHA256HashCode(message []byte)string{
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