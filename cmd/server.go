package cmd

import (
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/zu1k/nali/pkg/entity"
)

//查询到的IP信息
type Ipinfo struct {
	Code int    `json:"code" xml:"code"`
	Ip   string `json:"ip" xml:"ip"`
	Addr string `json:"addr" xml:"addr"`
	Cdn  string `json:"cdn" xml:"cdn"`
}

//错误状态信息
type Errinfo struct {
	Code int    `json:"code" xml:"code"`
	Msg  string `json:"msg" xml:"msg"`
}

//解析cname
func getCname(host string) string {
	host = strings.TrimSpace(host)
	cname, err := net.LookupCNAME(host)
	if err != nil {
		return ""
	}
	//判断cname是否存在
	if host == strings.TrimRight(cname, ".") {
		//cname不存在，返回空值
		return ""
	}
	//fmt.Println("CNAME地址：", strings.TrimRight(cname, "."))
	return strings.TrimRight(cname, ".")
}

//域名解析ip
func getIp(host string) string {
	host = strings.TrimSpace(host)
	ipall, err := net.LookupIP(host)
	if err != nil {
		return host
	}
	//fmt.Println("IP地址：", ipall)
	return ipall[0].String()
}

// Web server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start web api server",
	Long:    `Start web api server`,
	Example: "nali server --port 8080",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		port = ":" + port

		log.Println(">>> Nali Server running on *" + port + " ...")
		//关闭调试模式
		gin.SetMode(gin.ReleaseMode)
		//启动Gin服务
		ec := gin.Default()
		//参数方式查询
		ec.GET("/", func(c *gin.Context) {
			//设置响应头信息
			c.Header("Server", "nginx/1.22.0")
			ip := c.Query("ip")
			cname := ""
			if ip == "help" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "Param: /223.5.5.5 or ip=223.5.5.5 (Default: visitor ip)",
				}
				c.JSON(200, reinfo)
				return
			} else if ip == "" {
				//获取访客ip
				ip = c.ClientIP()
			} else if net.ParseIP(ip) == nil {
				//不是ip地址，当作域名解析
				cname = getCname(ip) //解析cname
				ip = getIp(ip)       //解析ip
			}
			//fmt.Println(ip)
			args = nil
			//判断cdn信息是否为空
			if cname == "" {
				args = append(args, ip)
			} else {
				args = append(args, ip, cname)
			}
			res := entity.ParseLine(strings.Join(args, " "))
			//判断是否有结果返回
			if res[0].Info == "" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "No record, please try again",
				}
				c.JSON(200, reinfo)
				return
			} else {
				//判断是否有CDN信息返回
				if cname != "" && res[2].Info != "" {
					cname = strings.Replace(res[2].Info, "\t", " ", -1)
				}
				reinfo := &Ipinfo{
					Code: 1,
					Ip:   res[0].Text,
					Addr: strings.Replace(res[0].Info, "\t", " ", -1),
					Cdn:  cname,
				}
				c.JSON(200, reinfo)
				return
			}

		})

		//路径查询方式
		ec.GET("/:ip", func(c *gin.Context) {
			//设置响应头信息
			c.Header("Server", "nginx/1.22.0")
			ip := c.Param("ip")
			cname := ""
			if ip == "favicon.ico" {
				//不响应浏览器请求favicon.ico图标
				return
			} else if ip == "help" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "Param: /223.5.5.5 or ip=223.5.5.5 (Default: visitor ip)",
				}
				c.JSON(200, reinfo)
				return
			} else if ip == "" {
				//获取访客ip
				ip = c.ClientIP()
			} else if net.ParseIP(ip) == nil {
				//不是ip地址，当作域名解析
				cname = getCname(ip) //解析cname
				ip = getIp(ip)       //解析ip
			}
			//fmt.Println(ip)
			args = nil
			//判断cdn信息是否为空
			if cname == "" {
				args = append(args, ip)
			} else {
				args = append(args, ip, cname)
			}
			res := entity.ParseLine(strings.Join(args, " "))
			//判断是否有结果返回
			if res[0].Info == "" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "No record, please try again",
				}
				c.JSON(200, reinfo)
				return
			} else {
				//判断是否有CDN信息返回
				if cname != "" && res[2].Info != "" {
					cname = strings.Replace(res[2].Info, "\t", " ", -1)
				}
				reinfo := &Ipinfo{
					Code: 1,
					Ip:   res[0].Text,
					Addr: strings.Replace(res[0].Info, "\t", " ", -1),
					Cdn:  cname,
				}
				c.JSON(200, reinfo)
				return
			}

		})
		//ec.Logger.Fatal(ec.Start(port))
		ec.Run(port)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("port", "8080", "Set web service listen port")
}
