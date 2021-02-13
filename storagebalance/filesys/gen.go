package filesys

import (
	"encoding/json"
	"fmt"

	"math/rand"
	"os"
)

const ID_max_byte = 20
const argeSize = 720
const UnifySize = 800
const R_smallSize = 240
const R_midSize = 460
const R_unifySize = 800
const R_bigSize = 1770

type TestFile struct {
	ID   [ID_max_byte]byte `json:"id"`   //对于长度小于2^64位的消息,SHA1会产生一个160位的消息摘要
	Size int64             `json:"size"` // MB
}

//等大小的num个800+-80文件
func GenerateFileInfo(num int) error {
	seed := int64(1613131121) // time.Now().Unix()
	rand.Seed(seed)
	var FilesInfo []TestFile
	bigfilenum := 0
	leftTreenum := 0
	size :=int64(0)
	for i := 0; i < num; i++ {
		var newID [20]byte
		for j := 0; j < ID_max_byte; j++ {
			rnd := rand.Intn(256)
			newID[j] = byte(rnd)
			if j == ID_max_byte-1 && newID[j] < 128 {
				leftTreenum++
			}
		}
		up := rand.Intn(160)
		newSize := (int64)(argeSize + up) //800 +-80
		size +=newSize
		if up >= 80 {
			bigfilenum++
		}
		FilesInfo = append(FilesInfo, TestFile{newID, newSize})
	}
	err := GenerateJson(&FilesInfo)
	if err != nil {
		return fmt.Errorf("GenerateJson Error : %v ", err)
	}
	fmt.Printf("seed-%v report:\v", seed)
	fmt.Printf("filenum: %v bigfilenum : %v %% leftTreenum: %v %%", num, 100*bigfilenum/num, 100*leftTreenum/num)
	fmt.Printf("size: %v ",size)

	return nil
}

func GenerateJson(FilesInfo *[]TestFile) error {
	jsonFile, err := os.Create("FileInfo.json")
	if err != nil {
		return fmt.Errorf("Error creating JSON file: %v ", err)
	}
	//创建编码器
	encoder := json.NewEncoder(jsonFile)
	//结构编码至JSON文件
	err = encoder.Encode(FilesInfo)
	if err != nil {
		return fmt.Errorf("Error encoding JSON to file: %v ", err)
	}
	return nil
}

//读取
func ReadFileInfo(filename string) ([]TestFile, error) {
	filePtr, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Open file failed [Err:%s] ", err.Error())
	}
	defer filePtr.Close()

	var FilesInfo []TestFile
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&FilesInfo)
	if err != nil {
		return nil, fmt.Errorf("Decoder failed: %v ", err.Error())
	}
	return FilesInfo, nil
}
