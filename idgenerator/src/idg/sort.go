package idg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

// 1. 先使用 seek 读取数据
// 2. 使用快排排序，写入文件
// 3. 分别取出排好序的文件第一行数据，将数据最小的写入新文件(使用堆排、最小堆)

// LIMIT 每个临时文件大小
const LIMIT int = 8757240

// LENGTH  最终排序后每个文件数据条数
const LENGTH int = 875124000

type ID struct {
	MD5    string
	Number uint64
	MD51   uint64
	MD52   uint64
}

func partition(intervals []*ID, p, r int) int {
	// 分区，把大于分区的放分区点右边，小于分区的放左边
	pivot := intervals[r].MD5
	i := p
	for j := p; j <= r-1; j++ {
		if intervals[j].MD5 < pivot {
			intervals[i], intervals[j] = intervals[j], intervals[i]
			i = i + 1
		}
	}
	intervals[i], intervals[r] = intervals[r], intervals[i]
	return i
}

func quickSort(intervals []*ID, p, r int) []*ID {
	if p >= r {
		return intervals
	}
	q := partition(intervals, p, r) // 获取分区点
	quickSort(intervals, p, q-1)
	quickSort(intervals, q+1, r)
	return intervals
}

type SortedID struct {
	FileName string
	Offset   int64
	MD5      string
	Number   uint64
	MD51     uint64
	MD52     uint64
}

type Heap struct {
	nodes       []*SortedID
	startIndex  int
	sortedIndex []*SortedID
}

func insert(elements []*SortedID, index int, value *SortedID) []*SortedID {
	elements = append(elements, nil) // 最后添加一个元素
	copy(elements[index+1:], elements[index:])
	elements[index] = value
	return elements
}

func delete(elements []*SortedID, index int) []*SortedID {
	// 1. 创建一个length - 1 数组
	newElements := make([]*SortedID, len(elements)-1, len(elements)-1)
	copy(newElements[:index], elements[:index]) // 复制index 前的元素到新列表
	if index != len(elements)-1 {
		// 如果index 不是最后一个索引
		copy(newElements[index:], elements[index+1:]) // copy index 后的元素到新列表
	}
	elements = newElements
	return elements
}

// Build 创建一个空堆 时间复杂度 O(n)
func (heap *Heap) Build(values []*SortedID) *Heap {
	// 先按完全二叉树规则依次填入数据，再堆化
	// 从最后一个节点的父节点开始
	// 为了简化计算，把root 节点放到 index 为 1 的位置
	if len(values) == 1 {
		heap.nodes = values
		heap.startIndex = 0
	} else {
		values = insert(values, 0, nil) // 在index = 0 位置添加一个无用元素
		length := len(values) - 1
		for i := length / 2; i >= 1; i-- {
			heap.startIndex = 1
			heap.Heapify(values, length, i)
		}
	}
	return heap
}

// Insert 向堆中插入一个元素 时间复杂度O(n)
func (heap *Heap) Insert(val *SortedID) *Heap {
	// 插入最后一个节点，然后调整
	nodes := heap.nodes
	if nodes == nil || len(nodes) == 0 { // 如果是第一个元素
		nodes = []*SortedID{val}
		heap.nodes = nodes
		heap.startIndex = 0
	} else {
		nodes = append(nodes, val)
		heap = heap.Build(nodes) // 直接重建堆
	}
	return heap
}

// Get 取堆顶元素
func (heap *Heap) Get() (*SortedID, error) {
	nodes := heap.nodes
	if nodes == nil {
		return nil, errors.New("Heap is Empty")
	}
	return nodes[0], nil
}

// Delete 删除堆顶元素 时间复杂度O(logn)
func (heap *Heap) Delete() (*SortedID, error) {
	// 删除二叉树的根或父节点。
	// 删除该节点元素后，队列最后一个元素必须移动到堆得某个位置，
	// 使得堆仍然满足堆序性质
	nodes := heap.nodes
	length := len(nodes)
	if nodes == nil || length == 0 {
		return nil, errors.New("Heap is Empty")
	}
	top := nodes[0]
	nodes[0] = nodes[length-1]
	copy(nodes[1:], nodes[:length-1])
	heap.startIndex = 1
	heap.Heapify(nodes, length-1, 1)
	return top, nil
}

// Heapify 使删除堆顶元素的堆再次成为堆
func (heap *Heap) Heapify(nodes []*SortedID, length, i int) *Heap {
	for {
		maxPos := i
		if i*2 < length && nodes[i].MD5 > nodes[i*2].MD5 {
			// 有左节点 且 节点值大于左节点的值
			maxPos = i * 2
		}
		if i*2+1 < length && nodes[maxPos].MD5 > nodes[i*2+1].MD5 {
			// 有右节点 且 节点值大于左节点的值
			maxPos = i*2 + 1
		}
		if maxPos == i { //  说明无需调整
			break
		}
		nodes[maxPos], nodes[i] = nodes[i], nodes[maxPos]
		i = maxPos
	}
	heap.nodes = nodes[heap.startIndex:]
	heap.startIndex = 0
	return heap
}

func getIds(numbers []uint64) []*ID {
	var ids []*ID
	for i := 0; i < len(numbers); i = i + 3 {
		md51 := numbers[i]
		md52 := numbers[i+1]
		number := numbers[i+2]
		md5 := Uint642Md5(md51, md52)
		id := &ID{md5, number, md51, md52}
		ids = append(ids, id)
		// fmt.Println(number, md51, md52)
	}
	fmt.Println("length: ", len(numbers))
	return ids
}

func WriteData(ids []*ID, fileName string) {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		log.Println("create file error: ", err)
		return
	}
	for _, id := range ids {
		// 将数据转为int64 整数存储
		data := []interface{}{id.MD51, id.MD52, id.Number}
		err = Write2Binary(f, data)
		if err != nil {
			log.Println("write 2 binary err: ", err)
			return
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func splitFile(fileName, tmpPath string, limit int) {
	var offset int64
	var n int
	var bytesData []byte
	wg := sync.WaitGroup{}
	for {
		offset, n, bytesData = ReadFromBinary(fileName, int64(offset), int64(limit))
		fmt.Println("----------", offset, n, len(bytesData))
		if n <= 0 {
			break
		}
		offset = offset + int64(n)
		wg.Add(1)
		go func() {
			defer wg.Done()
			numbers := Bytes2Uint64(bytesData, n)
			ids := getIds(numbers)
			ids = quickSort(ids, 0, len(ids)-1)
			start, end := ids[0].MD5, ids[len(ids)-1].MD5
			// fmt.Println(start, end)
			tf := fmt.Sprintf("%s/%s-%s.bin", tmpPath, start, end)
			WriteData(ids, tf)
		}()
	}
	wg.Wait()
	log.Println("Split file successfully")
}

// NewSortedID 创建一个新的SortedID
func NewSortedID(offset int64, n int, bytesData []byte) *SortedID {
	numbers := Bytes2Uint64(bytesData, n)
	md51 := numbers[0]
	md52 := numbers[1]
	number := numbers[2]
	md5 := Uint642Md5(md51, md52)
	sortedID := &SortedID{Offset: offset, MD5: md5, MD51: md51, MD52: md52, Number: number}
	return sortedID
}

func getSingalID(filePath string, file *os.File, offset int64) *SortedID {
	offset, n, bytesData := ReadData(file, offset, 24) // 每个文件读取一条数据
	offset = offset + int64(n)
	if n <= 0 {
		file.Close()
		return nil
	}
	sortedID := NewSortedID(offset, n, bytesData)
	sortedID.FileName = filePath
	return sortedID
}

func mergeFile(sortedFile string, splitFiles []string) {
	// 首先遍历1000个文件，每个文件里面取第一个数字，组成 (数字, 文件号) 这样的组合加入到堆里（假设是从小到大排序，用小顶堆），遍历完后堆里有1000个 (数字，文件号) 这样的元素
	// 然后不断从堆顶拿元素出来，每拿出一个元素，把它的文件号读取出来，然后去对应的文件里，加一个元素进入堆，直到那个文件被读取完。拿出来的元素当然追加到最终结果的文件里。
	// 按照上面的操作，直到堆被取空了，此时最终结果文件里的全部数字就是有序的了。
	log.Printf("Start merge file: %s。", sortedFile)
	var ids []*SortedID
	files := make(map[string]*os.File)
	for _, filePath := range splitFiles {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		files[filePath] = file
		sortedID := getSingalID(filePath, file, 0)
		if sortedID == nil {
			continue
		}
		ids = append(ids, sortedID)
	}
	heap := Heap{}
	heap.Build(ids)
	sortedfile, err := os.Create(sortedFile)
	if err != nil {
		log.Fatal(err)
	}
	defer sortedfile.Close()
	for {
		// 取最小值
		id, err := heap.Delete()
		if err != nil { // 说明已经没有数据了
			break
		}
		// 写入文件
		data := []interface{}{id.MD51, id.MD52, id.Number}
		err = Write2Binary(sortedfile, data)
		// 再从最小的数据中取出一个新的
		file := files[id.FileName]
		sortedID := getSingalID(id.FileName, file, id.Offset)
		if sortedID == nil {
			// 如果数据为空
			continue
		}
		heap.Insert(sortedID)
	}
	log.Printf("File: %s merged。", sortedFile)
}

// SortSingleAreaIds 给单个地区ID排序
func SortSingleAreaIds(areaCode int, fileName string) {
	tmpPath := fmt.Sprintf("db-%d/tmp", areaCode)
	// create dir
	os.Mkdir(tmpPath, os.ModePerm)
	splitFile(fileName, tmpPath, LIMIT)
	files, err := ioutil.ReadDir(tmpPath)
	if err != nil {
		log.Fatal(err)
	}
	var splitFiles []string
	for _, f := range files {
		fname := f.Name()
		if strings.HasSuffix(fname, ".bin") {
			splitFiles = append(splitFiles, fmt.Sprintf("%s/%s", tmpPath, fname))
		}
	}
	fmt.Println(splitFiles)
	sorteDir := "sorteDB"
	os.Mkdir(sorteDir, os.ModePerm)
	sortedFile := fmt.Sprintf("%s/%d.bin", sorteDir, areaCode)
	mergeFile(sortedFile, splitFiles)
	// delete tmp path
}

// SortAllIds 排序全部数据
func SortAllIds(filepath string) {
	// filepath 为所有已排序数据存放目录
}
