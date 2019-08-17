package main

import (
	"fmt"
	idg "md52id/idg"
)

func main() {
	// id md5 数据初始化，结果为按区号分类已经排好序的文件
	// idg.InitID()
	// 排序全部数据
	// idg.SortAllIds("dbs", "./sorteDB")
	// 索引初始化
	idg.LoadIndex(false)
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
	// 	offset, n, bytesData := idg.ReadFromBinary("sorteDB/330227.bin", int64(offset), 24)
	// 	numbers := idg.Bytes2Uint64(bytesData, n)
	// 	// fmt.Println(offset, numbers)
	// 	fmt.Println(idg.Uint642Md5(numbers[0], numbers[1]), offset)
	// }
	// fmt.Println("offset: ", offset, "data: ", data)
	// idg.SortSingleAreaIds(652222, "db-652222/652222.bin")
	// idg.LoadIndex(false)
	idBase := "33022719891020695"
	idNumber := idg.IDNumberFill(idBase)
	fmt.Println("id: ", idNumber)
	md5 := idg.Md5(idNumber)
	fmt.Println("md5: ", md5)
	idg.FindIDNumberFromMem(md5)
}
