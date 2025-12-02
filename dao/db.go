package dao //可以看出dao包用于数据库的处理

import (
	"encoding/json"
	"os"
	"painting/model"
	"sync"
)

/*
paintMu 是一把锁。
因为 Go 的 map 在多个人同时写入时会报错（并发冲突），
所以每次写文件前要先锁住，写完再开锁
*/
var paintMu sync.Mutex
var UserMap = map[string]string{}
var WorksMap = map[string][]model.Work{}
var dbFile = "database.json"

type DataStorage struct {
	Users map[string]string       `json:"users"`
	Works map[string][]model.Work `json:"works"`
}

func Init() {
	//初始化
	file, err := os.ReadFile(dbFile) //打开数据库dbFile放入file中
	if err == nil {                  //无错
		var data DataStorage        //声明临时数据库结构
		json.Unmarshal(file, &data) //将file（JSON）转换到data（临时数据库）
		UserMap = data.Users        //将临时数据库的数据存入传输数据中
		WorksMap = data.Works       //同上
	}
	// 防止空指针
	if UserMap == nil { //即数据库中无内容则初始化
		UserMap = make(map[string]string)
	}
	if WorksMap == nil {
		WorksMap = make(map[string][]model.Work)
	}
}

// 向数据库中写入数据
func Save() {
	paintMu.Lock()         //将要对数据库操作
	defer paintMu.Unlock() //最后解锁

	data := DataStorage{ //将传输数据转到临时库结构中
		Users: UserMap,
		Works: WorksMap,
	}
	// MarshalIndent 比 Marshal 多了缩进（"  "），让保存的 JSON 文件人眼看起来更好看
	file, _ := json.MarshalIndent(data, "", "  ") //将临时库的数据转换成JSON进入file中
	os.WriteFile(dbFile, file, 0666)              //将file（Json内容）存入数据库中，0666是读写权限
}

// 实现添加作品逻辑
func AddWork(username string, work model.Work) {
	//参数为要保存的数据，参考model中，从此处可以看出model用于定义数据结构
	work.Author = username                                //此处为本项目特殊要求，因为work数据结构中需要作者名
	WorksMap[username] = append(WorksMap[username], work) //将work放入传输数据中
	Save()                                                //将传输数据保存
}

// 查人（给登录用）
func CheckUser(username, password string) bool {
	//通过从context读入username
	// ok 代表 map 里有没有这个 key（人存不存在）
	// pwd 是 map 里存的真密码
	if pwd, ok := UserMap[username]; ok && pwd == password {
		return true
	}
	return false
}

// 加人（给注册用）
func AddUser(username, password string) bool {
	if _, ok := UserMap[username]; ok {
		return false // 已存在，即不为空
	}
	UserMap[username] = password
	Save() //保存
	return true
}

// 实现获取work逻辑，应该也为项目特殊
func GetWorks(username string) []model.Work {
	return WorksMap[username] // 直接返回这个人的作品切片
}

// 删除画作
func DelectPaint(username string, workname string) bool {
	for num := range WorksMap[username] {
		if WorksMap[username][num].Title == workname {
			WorksMap[username] = append(WorksMap[username][:num], WorksMap[username][num+1:]...)
			Save()
			return true
		}
	}
	return false
}

// 添加评论
func AddComment(username string, workname string, comment model.Comment) bool {
	paintMu.Lock()
	defer paintMu.Unlock()
	for num := range WorksMap[username] {
		if WorksMap[username][num].Title == workname {
			WorksMap[username][num].Comments = append(WorksMap[username][num].Comments, comment)
			WorksMap[username] = WorksMap[username]
			Save()
			return true
		}
	}
	return false
}
