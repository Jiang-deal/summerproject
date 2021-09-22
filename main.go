package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
)

type meg struct {
	Electrict float64
	Absence   int
	Len1      int //一共多少天
	Len2      int //缺失的天数的总数
}

// type Date1 struct {
// 	Date int
// }

func oper(add string, array map[int]float64, array_rep []int, array_rep1 []int) ([]int, []int) {
	//add是需要读取的文件地址

	//array := make(map[int]float64, 31) //定义map存储(日期，耗电量)
	//array_rep := make([]int, 0, 31) //保存确实数据的日期

	fi, err := os.Open(add)
	if err != nil {
		fmt.Printf("Error: %s\n", err)

	}
	defer fi.Close()

	//新建一个缓冲区,把内容放在缓冲区
	br := bufio.NewReader(fi)
	i := 1
	//给map中填入数据
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {

			break
		}
		var cli string

		fmt.Println("hfahfd", a)

		cli = string(a)
		data, _ := strconv.ParseFloat(cli, 64) //将string转为int
		fmt.Println("hfahfd", a)
		if data == 0 {
			array_rep = append(array_rep, i) //append相当于创建了一个新的变量
		} else {
			array_rep1 = append(array_rep1, i)
		}

		array[i] = float64(data)
		// array = append(array, data)
		// array_rep = append(array_rep, 0)
		i++
		//fmt.Println(array)

		//这部分是进行byte数组转换为int
		// fmt.Println(a)
		// b_buf := bytes.NewBuffer(a)

		// var x int32

		// binary.Read(b_buf, binary.LittleEndian, &x)

		// fmt.Println(x)

	}

	//copy(array_rep, array) //copy的时候如果没有初始化的话，不能进行复制
	fmt.Println("原始数据:", array)
	return array_rep, array_rep1

	// sup := 0
	// sort.Ints(array_rep)
	// len := len(array_rep)
	// fmt.Println(array_rep)
	// if len%2 == 1 {
	// 	sup = array_rep[len/2]

	// 	fmt.Println("中位数为:", sup)
	// } else {
	// 	u := len / 2
	// 	sup = (array_rep[u] + array_rep[u-1]) / 2
	// 	fmt.Println("中位数为:", sup)
	// }

	// for index, value := range array {
	// 	if value == 0 {
	// 		array[index] = sup
	// 	}
	// }
	// fmt.Println("中位数补充数据", array)

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

//数据集
func createdata(array_rep1 []int, array map[int]float64) ([]float64, []float64) {
	x := make([]float64, 0)
	y := make([]float64, 0)
	for i := 0; i < len(array_rep1); i++ {
		x = append(x, float64(array_rep1[i]))
		y = append(y, float64(array[array_rep1[i]])) //存在日期的用电量

	}
	fmt.Println("存在天数", x)
	fmt.Println("存在天数的用电量", y)

	return x, y
}

func main() {
	go registerServer("81.68.192.244:8500", "02", "DataSupply", "81.68.192.244", 9092, []string{"v1000"}, "Data", 8889) //开启服务

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
	router.LoadHTMLFiles("./static/index.tmpl")

	router.GET("/DataSupply/posts/index", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "get/index.tmpl", nil)
	})
	router.POST("/DataSupply/upload", func(c *gin.Context) {
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

		array := make(map[int]float64, 31) //数据
		array_rep := make([]int, 0, 31)    //保存缺失数据的日期
		array_rep1 := make([]int, 0, 31)   //保存存在数据的日期

		array_rep, array_rep1 = oper(dst, array, array_rep, array_rep1) //填补缺失数据
		fmt.Println(array)
		fmt.Println(array_rep)
		fmt.Println(array_rep1)
		//现在需要使用算法去搭建二次函数模型

		x, y := createdata(array_rep1, array)
		fmt.Println(x, y)

		// 开始数据拟合
		actPoints := plotter.XYs{}
		for a := 0; a < len(x); a++ {
			actPoints = append(actPoints, plotter.XY{
				X: x[a],
				Y: y[a],
			})
		}
		result, err := optimize.Minimize(optimize.Problem{
			Func: func(x []float64) float64 {

				if len(x) != 3 {
					panic("illegal x")
				}
				a := x[0]

				b := x[1]
				c := x[2]
				var sum float64
				for _, point := range actPoints {
					y := a*point.X*point.X + b*point.X + c
					sum += math.Abs(y - point.Y) //返回绝对值
				}
				return sum
			},
		},
			[]float64{1, 1, 1}, &optimize.Settings{}, &optimize.NelderMead{})
		if err != nil {
			panic(err)
		}
		fa, fb, fc := result.X[0], result.X[1], result.X[2]
		fmt.Println(fa, fb, fc) //a*point.X*point.X + b*point.X + c
		//二次曲线拟合完成，进行数据预测,将预测结果放入map中
		for _, v := range array_rep {
			array[v] = (fa*float64(v)*float64(v) + fb*float64(v) + fc)
		}
		fmt.Println(array)
		fmt.Println(len(array))
		length := len(array_rep)
		fmt.Println("mii", len(array_rep))
		//对缺失的数据进行填充
		for i := len(array_rep); i < len(array); i++ {
			array_rep = append(array_rep, 0)
		}

		fmt.Println("mii", array)
		fmt.Println("mii", len(array_rep))
		//定义返回的结构体
		u := make([]meg, 0, 31)
		for i := 0; i < len(array); i++ {
			var temp meg
			temp.Len1 = len(array)
			temp.Len2 = length
			temp.Absence = array_rep[i]
			temp.Electrict = array[i+1]
			u = append(u, temp)

		}

		c.JSON(http.StatusOK, u)

	})
	router.Run(":9092")
}
