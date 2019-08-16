package idg

/*
索引分为两级结构：
1级索引结构为：
indexes = [
	(start1, file1),
	(start2, file2),
	(start3, file3),
	(start4, file4),
	(start5, file5),
	(start6, file6),
]
2级索引结构为
file = [
	(start1, offset1),
	(start2, offset2),
	(start3, offset3),
	(start4, offset4),
]
*/

func loadIndex2Redis() {

}

func loadIndex2Mem() {

}

// LoadIndex 加载索引
func LoadIndex(useRedis bool) {
	// 为了加快

}
