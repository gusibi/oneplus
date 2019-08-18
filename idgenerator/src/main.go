package main

import (
	// "fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	idg "md52id/idg"
)

func idNumber(c *gin.Context){
	idBase := c.Query("base") // get id base
	if idBase == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id_base"})
	}
	idNumber := idg.IDNumberFill(idBase)
	c.JSON(200, gin.H{
		"id_number": idNumber,
	})
}

func id2md5(c *gin.Context){
	id := c.Query("id") // get id base
	if !idg.ValidateIDNumber(id){
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
	}
	md5 := idg.Md5(id)
	c.JSON(200, gin.H{
		"md5": md5,
	})
}

func search(c *gin.Context){
	md5 := c.Query("md5") // get query param
	idNumber:=idg.FindIDNumberFromMem(md5)
	c.JSON(200, gin.H{
		"id_number": idNumber,
	})
}

func main() {
	// id md5 数据初始化，结果为按区号分类已经排好序的文件
	// idg.InitID()
	// 排序全部数据
	// idg.SortAllIds("dbs", "./sorteDB")
	// 索引初始化
	idg.LoadIndex(false)
	// start server
	r := gin.Default()
	r.GET("/search", search)
	r.GET("/id2md5", id2md5)
	r.GET("/idNumber", idNumber)
	r.Run() // listen and serve on 0.0.0.0:8080
}
