package cmd

import (
	"log"
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
}

//错误状态信息
type Errinfo struct {
	Code int    `json:"code" xml:"code"`
	Msg  string `json:"msg" xml:"msg"`
}

// Web server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start web api server",
	Long:    `Start web api server`,
	Example: "nali server --port 8080",
	Run: func(cmd *cobra.Command, args []string) {
		/*

			ip := res[0].Text
			addr := strings.Split(strings.Replace(res[0].Info,"\t"," ",-1)," ")
			country := addr[0]
			area := addr[1]
			fmt.Println(ip,country,area)
		*/
		port, _ := cmd.Flags().GetString("port")
		port = ":" + port

    log.Println(">>> Nali Server running on *"+port)
    //关闭调试模式
    gin.SetMode(gin.ReleaseMode)
		//启动Gin服务
		ec := gin.Default()
		//参数方式查询
		ec.GET("/", func(c *gin.Context) {
			ip := c.Query("ip")
			if ip == "help" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "Param: /223.5.5.5 or ip=223.5.5.5 (Default: visitor ip)",
				}
				c.JSON(200, reinfo)
				return
			} else if ip == "" {
				ip = c.ClientIP()
			}
			//fmt.Println(ip)
			args = nil
			args = append(args, ip)
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
				reinfo := &Ipinfo{
					Code: 1,
					Ip:   res[0].Text,
					Addr: strings.Replace(res[0].Info, "\t", " ", -1),
				}
				c.JSON(200, reinfo)
				return
			}

		})

		//路径查询方式
		//参数方式查询
		ec.GET("/:ip", func(c *gin.Context) {
			ip := c.Param("ip")
			if ip == "help" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "Param: /223.5.5.5 or ip=223.5.5.5 (Default: visitor ip)",
				}
				c.JSON(200, reinfo)
				return
			} else if ip == "" {
				ip = c.ClientIP()
			}
			//fmt.Println(ip)
			args = nil
			args = append(args, ip)
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
				reinfo := &Ipinfo{
					Code: 1,
					Ip:   res[0].Text,
					Addr: strings.Replace(res[0].Info, "\t", " ", -1),
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
