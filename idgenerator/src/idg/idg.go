package idg

// 身份证号生成器
/*
身份证号生成规则
公民身份号码是特征组合码，由前十七位数字本体码和最后一位数字校验码组成。排列顺序从左至右依次为六位数字地址码，八位数字出生日期码，三位数字顺序码和一位数字校验码。
地址码： 表示编码对象常住户口所在县(市、旗、区)的行政区划代码。对于新生儿，该地址码为户口登记地行政区划代码。需要没说明的是，随着行政区划的调整，同一个地方进行户口登记的可能存在地址码不一致的情况。行政区划代码按GB/T2260的规定执行。
出生日期码：表示编码对象出生的年、月、日，年、月、日代码之间不用分隔符，格式为YYYYMMDD，如19880328。按GB/T 7408的规定执行。原15位身份证号码中出生日期码还有对百岁老人特定的标识，其中999、998、997、996分配给百岁老人。
顺序码： 表示在同一地址码所标识的区域范围内，对同年、同月、同日出生的人编定的顺序号，顺序码的奇数分配给男性，偶数分配给女性。
校验码： 根据本体码，通过采用ISO 7064:1983,MOD 11-2校验码系统计算出校验码。算法可参考下文。前面有提到数字校验码，我们知道校验码也有X的，实质上为罗马字符X，相当于10.

校验码算法

将本体码各位数字乘以对应加权因子并求和，除以11得到余数，根据余数通过校验码对照表查得校验码。

加权因子：
+-----------------------------------------------------------+
|位置序号|1 |2 |3 |4 |5 |6 |7 |8 |9 |10|11|12|13|14|15|16|17|
+-----------------------------------------------------------+
|加权因子|7 |9 |10|5 |8 |4 |2 |1 |6 |3 |7 |9 |10|5 |8 |4 |2 |
+-----------------------------------------------------------+

校验码:

+----------------------------------------------------+
| 余数  | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 |
+----------------------------------------------------+
| 校验码| 1 | 0 | X | 9 | 8 | 7 | 6 | 5 | 4 | 3 | 2  |
+----------------------------------------------------+

算法举例：

本体码为11010519491231002

第一步：各位数与对应加权因子乘积求和1*7+1*9+0*10+1*5+***=167
第二步：对求和进行除11得余数167%11=2
第三步：根据余数2对照校验码得X

因此完整身份证号为：11010519491231002X

https://zhuanlan.zhihu.com/p/21286417
*/

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// STARTYEAR 开始年份
const STARTYEAR = 1920

// const STARTYEAR = 2008

// ENDYEAR 结束年份
const ENDYEAR = 2019

// MONTHS 月份是否为31天
var MONTHS = [12]int{1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1}

// CHECKCODE 校验码
var CHECKCODE = map[int]string{
	0:  "1",
	1:  "0",
	2:  "10",
	3:  "9",
	4:  "8",
	5:  "7",
	6:  "6",
	7:  "5",
	8:  "4",
	9:  "3",
	10: "2",
}

// WEIGHT  加权因子
var WEIGHT = map[int]int{
	1:  7,
	2:  9,
	3:  10,
	4:  5,
	5:  8,
	6:  4,
	7:  2,
	8:  1,
	9:  6,
	10: 3,
	11: 7,
	12: 9,
	13: 10,
	14: 5,
	15: 8,
	16: 4,
	17: 2,
}

// LeapYear 判断是否是闰年
func LeapYear(year int) bool {
	if (year%4 == 0 && year%100 != 0) || (year%400 == 0 && year%3200 != 0) {
		return true
	} else {
		return false
	}
}

func ValidateIDNumber(id string) bool{
	if len(id) != 18{
		return false
	}
	idBase := id[:17]
	if id != IDNumberFill(idBase){
		return false
	}
	return true
}

// IDNumberFill 添加校验位
func IDNumberFill(id string) string {
	sum := 0
	for index, r := range id {
		number := int(r - '0')
		sum = sum + number*WEIGHT[index+1]
	}
	checkCode := CHECKCODE[sum%11]
	if checkCode == "10" {
		checkCode = "X"
	}
	return fmt.Sprintf("%s%s", id, checkCode)
}

// IDGeneratorByDays 按天生成身份证号
func IDGeneratorByDays(areaCode, year, month, day int, idChan chan string) {
	for i := 1; i < 1000; i++ {
		idBase := fmt.Sprintf("%d%d%02d%02d%03d", areaCode, year, month, day, i)
		// fmt.Println(id_base)
		idNumber := idBase
		idStr := IDNumberFill(idBase)
		idChan <- fmt.Sprintf("%s %s", Md5(idStr), idNumber)
	}
}

// IDGeneratorByMonths 按月生成身份证号
func IDGeneratorByMonths(areaCode, year, month int, idChan chan string) {
	days := 30
	if month == 2 {
		if LeapYear(year) == true {
			days = 29
		} else {
			days = 28
		}
	} else if MONTHS[month-1] == 1 {
		days = 31
	}
	for day := 1; day <= days; day++ {
		IDGeneratorByDays(areaCode, year, month, day, idChan)
	}
}

// IDGeneratorByYear 按年生成身份证号
func IDGeneratorByYear(areaCode, year int, idChan chan string, wg *sync.WaitGroup) {
	wg.Add(1)

	for month := 1; month <= 12; month++ {
		// fmt.Println(MONTHS, "month: ", month, month-1)
		IDGeneratorByMonths(areaCode, year, month, idChan)
	}
	wg.Done()
}

// IDGenerator 生成身份证号
func IDGenerator(areaCode int) string {
	idChan := make(chan string)
	done := make(chan bool)
	wg := sync.WaitGroup{}
	for year := STARTYEAR; year <= ENDYEAR; year++ {
		go IDGeneratorByYear(areaCode, year, idChan, &wg)
	}
	// create dir
	os.Mkdir(fmt.Sprintf("db-%d", areaCode), os.ModePerm)
	fileName := fmt.Sprintf("db-%d/%d.bin", areaCode, areaCode)
	// go Write2File(fileName, idChan, done)
	// go WriteBytes2File(fileName, idChan, done)
	go WriteBinary2File(fileName, idChan, done)
	time.Sleep(1 * time.Second)
	wg.Wait()
	close(idChan)
	d := <-done
	if d == true {
		log.Printf("AreaCode: %d written successfully", areaCode)
	} else {
		fmt.Println("File writing failed")
	}
	return fileName
}

func Md5(id string) string {
	// md5 结果使用字符串表示
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(id))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func InitID() {
	i := 0
	for k, v := range GB2260 {
		if i > 5 {
			break
		}
		log.Printf("%s ids start init!\n", v)
		fileName := IDGenerator(int(k))
		log.Printf("%s ids init finished!!\n", v)
		//排序
		SortSingleAreaIds(int(k), fileName)
		log.Printf("%s ids sort finished!!\n", v)
	}
}
