package idg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Md52Uint64 md5 字符串转数字
func Md52Uint64(id string) (uint64, uint64) {
	// md5 结果使用数字表示
	md51, _ := strconv.ParseUint(id[:16], 16, 64)
	md52, _ := strconv.ParseUint(id[16:], 16, 64)
	return md51, md52
}

// Uint642Md5 md5 字符串转数字
func Uint642Md5(md51, md52 uint64) string {
	// md5 结果使用数字表示
	return fmt.Sprintf("%s%s", fmt.Sprintf("%x", md51), fmt.Sprintf("%x", md52))
}

// Mkdir 检查目录是否存在，如果不存在则创建
func Mkdir(dir string) {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("mkdir %s error: [%v]\n", dir, err)
	} else {
		log.Printf("mkdir %s success", dir)
	}
}

// Write2File 将数据以字符串形式写入
func Write2File(fileName string, idChan chan string, done chan bool) {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for id := range idChan {
		_, err = fmt.Fprintln(f, id)
		if err != nil {
			fmt.Println(err)
			done <- false
			return
		}
	}
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

// WriteBytes2File 写入bytes 数据到文件
func WriteBytes2File(fileName string, idChan chan string, done chan bool) {
	f, err := os.Create(fmt.Sprintf("%s.bytes", fileName))
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for id := range idChan {
		id := fmt.Sprintf("%s\n", id)
		idBytes := []byte(id)
		_, err = f.Write(idBytes)
		if err != nil {
			fmt.Println(err)
			done <- false
			return
		}
	}
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

func Write2Binary(f *os.File, data []interface{}) error {
	binBuf := new(bytes.Buffer)
	for _, v := range data {
		err := binary.Write(binBuf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}
	bytes := binBuf.Bytes()
	_, err := f.Write(bytes)
	return err
}

// WriteBinary2File 写入二进制 数据到文件
func WriteBinary2File(fileName string, idChan chan string, done chan bool) {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for id := range idChan {
		ids := strings.Split(id, " ")
		// 将数据转为int64 整数存储
		md51, md52 := Md52Uint64(ids[0])
		number, _ := strconv.ParseInt(ids[1], 10, 64)
		data := []interface{}{md51, md52, number}
		err = Write2Binary(f, data)
		if err != nil {
			fmt.Println(err)
			done <- false
			return
		}
	}
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

func ReadData(file *os.File, offset, limit int64) (int64, int, []byte) {
	// 用来计算offset的初始位置
	// 0 = 相对于文件开始位置 offset = 0+offset
	// 1 = 相对于当前位置    offset = offset+offset
	// 2 = 相对于文件结尾处  offset = end-offset
	whence := 0
	newPosition, err := file.Seek(offset, whence)
	if err != nil {
		log.Fatal(err)
	}

	// 从文件中读取len(b)字节的文件。
	// 返回0字节意味着读取到文件尾了
	// 读取到文件结束会返回io.EOF的error
	data := make([]byte, limit)
	n, err := file.Read(data)
	if err != nil {
		log.Printf("File: %s read completed\n", file.Name())
	}
	return newPosition, n, data
}

// ReadFromBinary 从二进制文件中读取指定字节的数据
func ReadFromBinary(fileName string, offset, limit int64) (int64, int, []byte) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	return ReadData(file, offset, limit)
}

func Bytes2Uint64(data []byte, n int) []uint64 {
	// bytes to int64
	var numbers []uint64
	for i := 0; i < n; i = i + 8 {
		_data := data[i : i+8]
		var number uint64
		err := binary.Read(bytes.NewBuffer(_data), binary.BigEndian, &number)
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, number)
	}
	// fmt.Println(numbers)
	return numbers
}
