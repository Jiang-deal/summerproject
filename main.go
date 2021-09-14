package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func oper(add string) {
	// var clo []string
	// i := 0
	array := make([]int, 0, 100)
	array_rep := make([]int, 0, 100)

	fi, err := os.Open(add)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	//新建一个缓冲区,把内容放在缓冲区
	br := bufio.NewReader(fi)
	i := 0

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		var cli string

		cli = string(a)
		data, _ := strconv.Atoi(cli) //将string转为int
		array = append(array, data)
		array_rep = append(array_rep, 0)
		i++
		//fmt.Println(array)

		//这部分是进行byte数组转换为int
		// fmt.Println(a)
		// b_buf := bytes.NewBuffer(a)

		// var x int32

		// binary.Read(b_buf, binary.LittleEndian, &x)

		// fmt.Println(x)

	}

	copy(array_rep, array) //copy的时候如果没有初始化的话，不能进行复制
	fmt.Println("原始数据:", array)
	sup := 0
	sort.Ints(array_rep)
	len := len(array_rep)
	fmt.Println(array_rep)
	if len%2 == 1 {
		sup = array_rep[len/2]

		fmt.Println("中位数为:", sup)
	} else {
		u := len / 2
		sup = (array_rep[u] + array_rep[u-1]) / 2
		fmt.Println("中位数为:", sup)
	}

	for index, value := range array {
		if value == 0 {
			array[index] = sup
		}
	}
	fmt.Println("中位数补充数据", array)

}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func main() {
	router := gin.Default()
	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB

	router.Static("xxx", "./static")
	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.LoadHTMLFiles("form2.tmpl")

	router.GET("/posts/index", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "get/index.tmpl", nil)
	})
	router.POST("/upload", func(c *gin.Context) {
		// 单个文件
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Println(file.Filename)
		dst := fmt.Sprintf("./file/%s", file.Filename)
		fmt.Println("nishi", dst)
		// 上传文件到指定的目录
		c.SaveUploadedFile(file, dst)
		oper(dst) //填补缺失数据
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})

	})
	router.Run()
}
