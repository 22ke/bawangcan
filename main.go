package main

import (
	"bawangcan/allrequest"
	"bawangcan/config"
	"os"
)

var Con config.Config

func main() {
	f, _ := os.OpenFile("log.txt", os.O_APPEND, 0644) //配置日志文件
	defer f.Close()

	con := config.Makeconfig() //初始化配置文件
	act, actid := allrequest.Getmenuinfos(con)
	////报名
	allrequest.Baoming(act, actid, con, f)
}
