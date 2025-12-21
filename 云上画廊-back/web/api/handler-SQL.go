package api

import (
	"fmt"
	"painting/box"
	"painting/dao/db"
)

// 关闭数据库
func CloseSQL() {
	if err := box.Temp.Close(); err != nil {
		fmt.Print(err)
		return
	}

}

// 初始化数据库
func InitSQL() {
	var ok bool
	box.Temp, ok = db.Init()
	if !ok {
		fmt.Println("数据库初始化失败，请检查 MySQL 是否已启动，以及 dao.Init 的配置（DSN）是否正确。")
		fmt.Println("按回车键退出...")
		fmt.Scanln()
		return
	}
}
