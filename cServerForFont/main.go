package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/exec"
	path2 "path"
)

/*
	端口默认5000
	路径默认test
	配置文件默认test.json
*/

var Res []byte

func ReturnData(c *gin.Context) {
	res := string(Res)
	c.JSON(http.StatusOK, res)
}

/**
跨域资源共享(CORS) 是一种机制，它使用额外的 HTTP 头来告诉浏览器  让运行在一个 origin (domain) 上的Web应用
被准许访问来自不同源服务器上的指定的资源。当一个资源从与该资源本身所在的服务器不同的域、协议或端口请求一个资源时，
资源会发起一个跨域 HTTP 请求。比如，站点 http://domain-a.com 的某 HTML 页面通过 <img> 的 src 请求 http://domain-b.com/image.jpg。
网络上的许多页面都会加载来自不同域的CSS样式表，图像和脚本等资源。出于安全原因，浏览器限制从脚本内发起的跨源HTTP请求。
例如，XMLHttpRequest和Fetch API遵循同源策略。 这意味着使用这些API的Web应用程序只能从加载应用程序的同一个域请求HTTP资源，
除非响应报文包含了正确CORS响应头。
*/
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			origin = c.Request.Header.Get("Origin") //请求头部
			//method     = c.Request.Method
			//headerKeys []string                         // 声明请求头keys
			//headerStr  string
		)
		//for k, _ := range c.Request.Header {
		//	headerKeys = append(headerKeys, k)
		//}
		//headerStr = strings.Join(headerKeys, ", ")
		//if headerStr != "" {
		//	headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		//} else {
		//	headerStr = "access-control-allow-origin, access-control-allow-headers"
		//}
		if origin != "" {
			//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			// 这是允许访问所有域
			c.Header("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//header的类型
			c.Header("Access-Control-Allow-Headers", `Authorization, Content-Length, X-CSRF-Token, Token,
							session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,
							DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since,
							Cache-Control, Content-Type, Pragma`)
			//允许跨域设置 可以返回其他子段
			// 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Expose-Headers", `Content-Length, Access-Control-Allow-Origin,
							Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,
							Last-Modified,Pragma,FooBar`)
			// 缓存请求信息 单位为秒
			c.Header("Access-Control-Max-Age", "172800")
			//跨域请求是否需要带cookie信息 默认设置为true
			c.Header("Access-Control-Allow-Credentials", "true")
			// 设置返回格式是json
			c.Set("content-type", "application/json")
		}
		////放行所有OPTIONS方法
		//if method == "OPTIONS" {
		//	c.JSON(http.StatusOK, "Options Request!")
		//}
		// 处理请求
		c.Next()
	}
}

func main() {

	//var tp string
	var method string
	var port string
	var src string
	var json_path string
	//flag.StringVar(&tp,"t","e","类别，[可选:e(简单),o（仅一个）]")
	flag.StringVar(&method, "m", "get", "请求方式，[可选:get,post]")
	flag.StringVar(&port, "p", "5000", "端口")
	flag.StringVar(&src, "r", "/test", "路径")
	flag.StringVar(&json_path, "c", "test.json", "json文件")
	flag.Parse()

	base, _ := exec.LookPath(os.Args[0])

	path := path2.Dir(base) + "/" + json_path
	fmt.Printf("method%v", method)
	fmt.Printf("port%v", port)
	fmt.Printf("src%v", src)
	fmt.Printf("json_path%v", json_path)
	fmt.Printf("path%v", path)
	jsonfile, _ := os.OpenFile(path, os.O_RDONLY, 777)
	var data = make([]byte, 1024)
	for {
		n, err := jsonfile.Read(data)
		if err == io.EOF {
			break
		}
		Res = append(Res, data[:n]...)
	}
	r := gin.New()
	r.Use(Cors())
	if method == "get" {
		r.GET(src, ReturnData)
	}
	if method == "post" {
		r.POST(src, ReturnData)
	}
	_ = r.Run(":" + port)
}
