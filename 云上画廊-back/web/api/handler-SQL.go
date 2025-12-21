package api

import (
	"fmt"
	"painting/box"
	"painting/dao/db"
	"time"
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
	num := 0
F:
	var ok bool
	time.Sleep(3 * time.Second)
	box.Temp, ok = db.Init()
	if !ok {
		if num < 10 {
			fmt.Println("加载中...")
			num++
			goto F
		} else {
			fmt.Println("失败!")
			return
		}
	}
	fmt.Println("成功!")
}
