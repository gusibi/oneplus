package main

import (
	idg "md52id/idg"
)

func main() {
	// id md5 数据初始化，结果为按区号分类已经排好序的文件
	idg.InitID()
	// 排序全部数据
	// idg.SortAllIds("./sorted_ids/")
	// 索引初始化
	// idg.LoadIndex(false)
	// start server
	// id := idg.IDNumberFill("10000019220101003")
	// md5 := idg.Md5(id)
	// md51, md52 := idg.Md52Uint64(md5)
	// fmt.Println(id, md5, md51, md52)
	// id = idg.IDNumberFill("10000019230101001")
	// md5 = idg.Md5(id)
	// md51, md52 = idg.Md52Uint64(md5)
	// fmt.Println(id, md5, md51, md52)
	// id = idg.IDNumberFill("10000019800101003")
	// md5 = idg.Md5(id)
	// md51, md52 = idg.Md52Uint64(md5)
	// fmt.Println(id, md5, md51, md52)
	// for offset := 0; offset <= 1200; offset = offset + 24 {
	// 	offset, n, bytesData := idg.ReadFromBinary("sorted-100000.bin", int64(offset), 24)
	// 	numbers := idg.Bytes2Uint64(bytesData, n)
	// 	fmt.Println(offset, numbers)
	// 	fmt.Println(idg.Uint642Md5(numbers[0], numbers[1]))
	// }
	// fmt.Println("offset: ", offset, "data: ", data)
	// idg.SortSingleAreaIds(321201, "321201/321201.bin")
}
