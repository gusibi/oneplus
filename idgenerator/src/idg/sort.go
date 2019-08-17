package idg

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 1. 先使用 seek 读取数据
// 2. 使用快排排序，写入文件
// 3. 分别取出排好序的文件第一行数据，将数据最小的写入新文件(使用堆排、最小堆)

// SplitLimit 每个临时文件大小
const SplitLimit int = 8757240

// MaxLength  最终排序后每个文件数据最大条数 每个文件768M
// const MaxLength int = 729270
const MaxLength int = 33554432

// ProgressFile 进度存储文件
const ProgressFile string = "db-progress"

type ID struct {
	MD5    string
	Number uint64
	MD51   uint64
	MD52   uint64
}

func partitionByID(intervals []*ID, p, r int) int {
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

// QuickSort 快排
func QuickSort(intervals []*ID, p, r int) []*ID {
	if p >= r {
		return intervals
	}
	q := partitionByID(intervals, p, r) // 获取分区点
	QuickSort(intervals, p, q-1)
	QuickSort(intervals, q+1, r)
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
		if i*2 <= length && nodes[i].MD5 > nodes[i*2].MD5 {
			// 有左节点 且 节点值大于左节点的值
			maxPos = i * 2
		}
		if i*2+1 <= length && nodes[maxPos].MD5 > nodes[i*2+1].MD5 {
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
	// log.Println("ids length: ", len(numbers))
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
		// fmt.Println("----------", offset, n, len(bytesData))
		if n <= 0 {
			break
		}
		offset = offset + int64(n)
		wg.Add(1)
		go func() {
			defer wg.Done()
			numbers := Bytes2Uint64(bytesData, n)
			ids := getIds(numbers)
			ids = QuickSort(ids, 0, len(ids)-1)
			if len(ids) > 0 {
				start, end := ids[0].MD5, ids[len(ids)-1].MD5
				// fmt.Println(start, end)
				tf := fmt.Sprintf("%s/%s-%s.bin", tmpPath, start, end)
				WriteData(ids, tf)
			}
		}()
	}
	wg.Wait()
	log.Println("Split file successfully")
}

// IDCache id 缓存，读出的数据先暂存到cache 中
/*
{
	"file1": {
		0: [md51, md52, number],
		8: [md51, md52, number],
		16: [md51, md52, number]
	}
}
*/
var IDCache = make(map[string]map[int64][]uint64)

func getIdFromCache(filePath string, file *os.File, offset int64) ([]uint64, int) {
	data, ok := IDCache[filePath]
	var n int
	var bytesData []byte
	if !ok {
		// 如果没有数据，从文件读取，再写入cache
		offset, n, bytesData = ReadData(file, offset, 2400) // 每个文件读取100条数据，存入IDCache
	} else {
		numbers, ok := data[offset]
		if !ok {
			// 如果没有则需要从文件读取，再写入cache
			offset, n, bytesData = ReadData(file, offset, 2400) // 每个文件读取100条数据，存入IDCache
		} else {
			return numbers, 24
		}
	}
	if n <= 0 {
		file.Close()
		return nil, 0
	}
	numbers := Bytes2Uint64(bytesData, n)
	data = make(map[int64][]uint64)
	for i := 0; i < len(numbers); i = i + 3 {
		data[offset] = numbers[i : i+3]
		offset = offset + 24
	}
	IDCache[filePath] = data
	return numbers[:3], 24
}

func getSingalID(filePath string, file *os.File, offset int64) *SortedID {
	// 这里可以优化，每次从文件读出多条数据先放入队列，
	// 每次从队列取一条数据返回，队列为空后再从文件读取多条
	// 减少文件读取次数，提高执行速度
	// 先判断缓存有没有数据
	// 存入IDCache
	numbers, n := getIdFromCache(filePath, file, offset)
	offset = offset + int64(n)
	if len(numbers) == 0 {
		return nil
	}
	md5 := Uint642Md5(numbers[0], numbers[1])
	// fmt.Println(md5, offset, filePath)
	sortedID := &SortedID{Offset: offset, MD5: md5, MD51: numbers[0], MD52: numbers[1], Number: numbers[2]}
	sortedID.FileName = filePath
	return sortedID
}

func mergeFile(sortedFile string, splitFiles []string) {
	// 首先遍历多个文件，每个文件里面取第一个数字，组成 (数字, 文件号) 这样的组合加入到堆里（假设是从小到大排序，用小顶堆），遍历完后堆里有1000个 (数字，文件号) 这样的元素
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
	log.Printf("File: %s merged", sortedFile)
}

// GetDirFiles 获取目录下的文件路径
func GetDirFiles(path string) (subFiles, fs []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fname := f.Name()
		if strings.HasSuffix(fname, ".bin") {
			fs = append(fs, fname)
			subFiles = append(subFiles, fmt.Sprintf("%s/%s", path, fname))
		}
	}
	return subFiles, fs
}

// SortSingleAreaIds 给单个地区ID排序
func SortSingleAreaIds(areaCode int, fileName string) {
	tmpPath := fmt.Sprintf("db-%d/tmp", areaCode)
	// create dir
	os.Mkdir(tmpPath, os.ModePerm)
	splitFile(fileName, tmpPath, SplitLimit)
	splitFiles, _ := GetDirFiles(tmpPath)
	sorteDir := "sorteDB"
	os.Mkdir(sorteDir, os.ModePerm)
	sortedFile := fmt.Sprintf("%s/%d.bin", sorteDir, areaCode)
	mergeFile(sortedFile, splitFiles)
	err := os.RemoveAll(tmpPath) // 执行完删除
	if err != nil {
		log.Fatal("Remove tmp path error: ", err)
	}
}

// Progress 进度
type Progress struct {
	LogFile string // 记录进度的文件
}

// NewProgress create new progress
func NewProgress(logFile string) *Progress {
	return &Progress{LogFile: logFile}
}

// Load 加载进度，如果没有数据 返回空
func (p *Progress) Load() (progress map[string]int64) {
	file, err := os.Open(p.LogFile)
	defer file.Close()
	progress = make(map[string]int64)
	if err == nil {
		// read data
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			rows := strings.Split(row, " ")
			fileName := rows[0]
			offset, err := strconv.ParseInt(rows[1], 10, 64)
			if err != nil {
				log.Fatal("Read progress offset error: ", err)
			}
			progress[fileName] = offset
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return progress
}

// Refresh 刷新进度，一个文件写满时，刷新；如果一个文件已经读取完成 offset=-1
func (p *Progress) Refresh(progress map[string]int64) {
	f, err := os.Create(p.LogFile)
	defer f.Close()
	if err != nil {
		log.Fatal("Create progress file error: ", err)
	}
	for fileName, offset := range progress {
		row := fmt.Sprintf("%s %d", fileName, offset)
		_, err = fmt.Fprintln(f, row)
		if err != nil {
			log.Fatal("Write progress file error: ", err)
		}
	}
}

func createNewSortedFile(fileName string) *os.File {
	sortedfile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return sortedfile
}

func writeSorteIds(dbPath string, sortedAreaFiles []string) {
	// 加载进度，如果没有则说明是从头开始
	// 进度为每个areaFile 读取到的进度
	// 首先遍历多个文件，每个文件里面取第一个数字，组成 (数字, 文件号) 这样的组合加入到堆里（假设是从小到大排序，用小顶堆），遍历完后堆里有1000个 (数字，文件号) 这样的元素
	// 然后不断从堆顶拿元素出来，每拿出一个元素，把它的文件号读取出来，然后去对应的文件里，加一个元素进入堆，直到那个文件被读取完。拿出来的元素当然追加到最终结果的文件里。
	// 按照上面的操作，直到堆被取空了，此时最终结果文件里的全部数字就是有序的了。
	log.Printf("Start merge file: %s。", sortedAreaFiles)
	var ids []*SortedID
	files := make(map[string]*os.File)
	progress := NewProgress(ProgressFile)
	progressData := progress.Load()
	for _, filePath := range sortedAreaFiles {
		offset, ok := progressData[filePath]
		if ok && offset == -1 {
			continue
		} else if !ok {
			offset = 0
		}
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		files[filePath] = file
		sortedID := getSingalID(filePath, file, offset)
		if sortedID == nil {
			continue
		}
		ids = append(ids, sortedID)
	}
	heap := Heap{}
	heap.Build(ids)
	id, _ := heap.Get()
	// 第一个文件
	sortedFile := fmt.Sprintf("%s/%s.bin", dbPath, id.MD5)
	sortedfile := createNewSortedFile(sortedFile)
	count := 0
	for {
		// 取最小值
		id, err := heap.Delete()
		if err != nil { // 说明已经没有数据了
			break
		}
		progressData[id.FileName] = id.Offset
		count++
		if count%MaxLength == 1 {
			sortedFile = fmt.Sprintf("%s/%s.bin", dbPath, id.MD5)
			sortedfile = createNewSortedFile(sortedFile)
			log.Printf("start write data to file: %s !\n", sortedFile)
		}
		// 写入文件
		data := []interface{}{id.MD51, id.MD52, id.Number}
		// fmt.Println(id)
		// fmt.Println("write   id from: ", id.FileName, id.MD5)
		err = Write2Binary(sortedfile, data)
		if count%MaxLength == 0 {
			// 说明文件已经到达写入上限
			sortedfile.Close()
			log.Println("close file: ", sortedFile)
			// 写入进度
			progress.Refresh(progressData)
		}
		// 再从最小的数据中取出一个新的
		file := files[id.FileName]
		sortedID := getSingalID(id.FileName, file, id.Offset)
		// fmt.Println("get new id from: ", sortedID.FileName, sortedID.MD5)
		// time.Sleep(1 * time.Second)
		// fmt.Println(sortedID)
		if sortedID == nil {
			// 如果数据为空 说明该文件已经读完 更新progress
			progressData[id.FileName] = -1
			continue
		}
		heap.Insert(sortedID)
	}
	if count%MaxLength != 0 {
		// 说明最后文件没有关闭
		sortedfile.Close()
		log.Println("close file: ", sortedFile)
	}
}

// SortAllIds 排序全部数据
func SortAllIds(targetPath, sourcePath string) {
	// filepath 为所有已排序数据存放目录
	// 读取 sorteDB 的文件，使用归并排序
	sourceFiles, _ := GetDirFiles(sourcePath)
	os.Mkdir(targetPath, os.ModePerm)
	writeSorteIds(targetPath, sourceFiles)
}
