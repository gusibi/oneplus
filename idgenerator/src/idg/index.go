package idg

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

// PageLimit 索引每页数据条数
const PageLimit = 1024

// DBPath 数据目录
const DBPath = "dbs"

/*
Indexes 索引，查找文件
索引结构为：
{
	"file1": ["md51", "md52", "md53"]
}
*/
var Indexes map[string][]string

func getKey(dbPath string) string {
	// 使用dbPath 获取索引key
	// 比如：dbs/7e5f959f6c2c5e8c62696ab4e23dbb78.bin
	// return 7e5f959f6c2c5e8c62696ab4e23dbb78
	paths := strings.Split(dbPath, "/")
	key := strings.Split(paths[len(paths)-1], ".")[0]
	return key
}

func getIndexes(file, key string, wg *sync.WaitGroup) {
	defer wg.Done()
	var offset int64
	var n int
	var bytesData []byte
	var secondIndexes []string
	for {
		offset, n, bytesData = ReadFromBinary(file, offset, 24)
		if n <= 0 {
			// fmt.Println("len: ", len(secondIndexes))
			Indexes[key] = secondIndexes
			break
		}
		numbers := Bytes2Uint64(bytesData, n)
		md5 := Uint642Md5(numbers[0], numbers[1])
		secondIndexes = append(secondIndexes, md5)
		offset = offset + PageLimit*24
	}
}

func loadIndex2Redis() {
	log.Println("Load index from redis")
}

func loadIndex2Mem() {
	dbPaths, dbFiles := GetDirFiles(DBPath)
	wg := sync.WaitGroup{}
	Indexes = make(map[string][]string)
	for index, dbPath := range dbPaths {
		key := getKey(dbFiles[index])
		Indexes[key] = nil
		wg.Add(1)
		go getIndexes(dbPath, key, &wg)
		fmt.Println(dbPath)
	}
	wg.Wait()
	log.Println("Load index to memc")
}

// LoadIndex 加载索引
func LoadIndex(useRedis bool) {
	// 为了加快启动速度，数据可以存储在redis
	if useRedis {
		loadIndex2Redis()
	} else {
		loadIndex2Mem()
	}
}

func partition(intervals []string, p, r int) int {
	// 分区，把大于分区的放分区点右边，小于分区的放左边
	// fmt.Println(intervals)
	pivot := intervals[r]
	i := p
	for j := p; j <= r-1; j++ {
		if intervals[j] < pivot {
			intervals[i], intervals[j] = intervals[j], intervals[i]
			i = i + 1
		}
	}
	intervals[i], intervals[r] = intervals[r], intervals[i]
	return i
}

// quickSort 快排
func quickSort(intervals []string, p, r int) []string {
	if p >= r {
		return intervals
	}
	q := partition(intervals, p, r) // 获取分区点
	quickSort(intervals, p, q-1)
	quickSort(intervals, q+1, r)
	return intervals
}

// BinarySearch 使用二分查找搜索数据
func BinarySearch(values []string, target string) int {
	// 使用二分法查找values 中小于等于target 且差值最小的元素
	left, right := 0, len(values)-1
	for left <= right {
		mid := (right + left) / 2
		// fmt.Println("left: ", left, "right: ", right, "mid: ", mid)
		if values[mid] == target {
			return mid
		} else if values[mid] > target {
			right = mid - 1
		} else if values[mid] < target {
			if (mid < len(values)-1 && values[mid+1] > target) || mid == len(values)-1 {
				// 小于target 且最接近
				return mid
			}
			left = mid + 1
		}
	}
	return -1
}

// BinarySearch 使用二分查找搜索数据
func BinarySearchFromBytes(values []byte, target string) uint64 {
	// 使用二分法查找values 中小于等于target 且差值最小的元素
	length := len(values)/24
	left, right := 0, length-1
	for left <= right {
		mid := (right + left) / 2
		offset := mid * 24
		// fmt.Println("left: ", left, "right: ", right, "mid: ", mid, "offset: ", offset)
		value := values[offset:offset+24]
		numbers := Bytes2Uint64(value, 24)
		md5 := Uint642Md5(numbers[0], numbers[1])
		if md5 == target {
			return numbers[2]
		} else if md5 > target {
			right = mid - 1
		} else if md5 < target {
			left = mid + 1
		}
	}
	return 0
}

func getKeys(indexes map[string][]string) []string {
	keys := make([]string, 0, len(indexes))
	for k := range indexes {
		keys = append(keys, k)
	}
	return keys
}

func findFile(md5 string) string {
	files := getKeys(Indexes)
	// fmt.Println("files:-----------", files)
	sortedFiles := quickSort(files, 0, len(files)-1)
	keyIndex := BinarySearch(sortedFiles, md5)
	// fmt.Println("key index: ", keyIndex)
	fileKey := sortedFiles[keyIndex]
	// fmt.Println(fileKey)
	return fileKey
}

func getIDNumber(md5, fileKey string) uint64 {
	fileName := fmt.Sprintf("%s/%s.bin", DBPath, fileKey)
	// fmt.Println("fileName", fileName)
	startIndex := BinarySearch(Indexes[fileKey], md5)
	// startMd5 := Indexes[fileKey][startIndex]
	// fmt.Println("start: ", startMd5)
	// nextMd5 := Indexes[fileKey][startIndex+1]
	// fmt.Println("next: ", nextMd5)
	offset := PageLimit * 24 * startIndex
	// fmt.Println("offset: ", offset)
	// 读取页数据
	_, n, bytesData := ReadFromBinary(fileName, int64(offset), PageLimit*24)
	if n == 0{
		return 0
	}
    id:=BinarySearchFromBytes(bytesData, md5)	
	return id
}

//FindIDNumberFromMem 从内存中查找数据
func FindIDNumberFromMem(md5 string) uint64 {
	//先从一级索引找到可能存在的文件
	// 再从文件查找对应的身份证号，如果找到返回对应的身份证号码，找不到返回空
	fileKey := findFile(md5)
	// fmt.Println("fileKey: ", fileKey)
	id := getIDNumber(md5, fileKey)
	return id
}

//FindIDNumberFromRedis 从Redis中查找数据
func FindIDNumberFromRedis(md5 string) {

}
