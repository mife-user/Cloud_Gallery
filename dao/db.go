package dao //可以看出dao包用于数据库的处理

import (
	"fmt"
	"os"
	"painting/model"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
paintMu 是一把锁。
因为 Go 的 map 在多个人同时写入时会报错（并发冲突），
所以每次写文件前要先锁住，写完再开锁
*/
// var paintMu sync.Mutex
// var UserMap = map[string]string{}
// var WorksMap = map[string][]model.Work{}
// var dbFile = "database.json"

type DataStorage struct {
	Users map[string]string       `json:"users"`
	Works map[string][]model.Work `json:"works"`
}
type Database struct {
	DB *gorm.DB
}

//	func Init() {
//		// 初始化
//		file, err := os.ReadFile(dbFile) //打开数据库dbFile放入file中
//		if err == nil {                  //无错
//			var data DataStorage        //声明临时数据库结构
//			json.Unmarshal(file, &data) //将file（JSON）转换到data（临时数据库）
//			UserMap = data.Users        //将临时数据库的数据存入传输数据中
//			WorksMap = data.Works       //同上
//		}
//		// 防止空指针
//		if UserMap == nil { //即数据库中无内容则初始化
//			UserMap = make(map[string]string)
//		}
//		if WorksMap == nil {
//			WorksMap = make(map[string][]model.Work)
//		}
//	}
func Init() (*Database, bool) {
	// 从环境变量读取，默认 localhost（本地开发）或 mysql（Docker 环境）
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}

	dsn := fmt.Sprintf("root:314159@tcp(%s:%s)/paint?charset=utf8mb4&parseTime=True&loc=Local", host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, false
	}
	db.AutoMigrate(&model.User{}, &model.Work{}, &model.Comment{})
	DB := &Database{DB: db}
	return DB, true
}

/*
内部写文件函数（假定调用方已经持有 paintMu）
把写文件的具体逻辑拆出来，避免死锁（调用方锁 -> 调用内部写）。
外部如果需要也可以直接调用 Save() 自动上锁。
*/
// func saveLocked() {
// 	data := DataStorage{ //将传输数据转到临时库结构中
// 		Users: UserMap,
// 		Works: WorksMap,
// 	}
// 	// MarshalIndent 比 Marshal 多了缩进（"  "），让保存的 JSON 文件人眼看起来更好看
// 	file, _ := json.MarshalIndent(data, "", "  ") //将临时库的数据转换成JSON进入file中
// 	// 写文件时不在意错误（简洁），但在生产环境应处理错误并考虑原子写法
// 	_ = os.WriteFile(dbFile, file, 0644) //将file（Json内容）存入数据库中，0644 是常用权限
// }

// // 对外暴露的 Save，会自动上锁（安全）
// func Save() {
// 	paintMu.Lock()
// 	defer paintMu.Unlock()
// 	saveLocked()
// }

// 实现添加作品逻辑
func (d *Database) AddWork(username string, work *model.Work) bool {
	// paintMu.Lock() // 上锁保护内存数据修改 + 持久化
	// defer paintMu.Unlock()

	// work.Author = username                                //此处为本项目特殊要求，因为work数据结构中需要作者名
	// WorksMap[username] = append(WorksMap[username], work) //将work放入传输数据中

	// 写文件（内部函数假定当前已上锁，避免再次尝试加锁导致死锁）
	// saveLocked()
	user := model.User{
		Username: username,
		Password: "",
	}
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	if err := d.DB.Model(&user).Association("Work").Append(work); err != nil {
		return false
	}
	return true
}

// 查人（给登录用）
func (d *Database) CheckUser(username, password string) bool {
	// 读操作也应当是并发安全的；这里为了性能没有加锁（map 读在大多数场景下是安全的）
	// 如果你期望严格安全，也可以在读时使用锁
	// if pwd, ok := UserMap[username]; ok && pwd == password {
	// 	return true
	// }

	// return false
	user := model.User{
		Username: username,
		Password: password,
	}
	if err := d.DB.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return false
	}
	return true
}

// 加人（给注册用）
func (d *Database) AddUser(username, password string) bool {
	// paintMu.Lock()
	// defer paintMu.Unlock()

	// if _, ok := UserMap[username]; ok {
	// 	return false // 已存在，即不为空
	// }
	// UserMap[username] = password
	// // 保存（内部函数假定已上锁）
	// // saveLocked()
	user := model.User{
		Username: username,
		Password: password,
	}
	if err := d.DB.Model(&user).Where("username = ?", username).First(&user); err != nil {
		// 用户不存在，创建新用户
		if err := d.DB.Create(&user).Error; err != nil {
			return false
		}
		return true
	}
	return false
}

// 实现获取work逻辑，应该也为项目特殊
func (d *Database) GetWorks(username string) []model.Work {
	// 这里直接返回切片引用（注意：调用方不要直接修改返回切片）
	// 若要防止外部篡改，可返回副本
	// works := WorksMap[username]
	user := model.User{
		Username: username,
		Password: "",
	}
	if err := d.DB.Model(&user).Where("username = ?", username).First(&user).Error; err != nil {
		return nil
	}
	var works []model.Work
	if err := d.DB.Preload("Comments").Where("user_id = ?", user.ID).Find(&works).Error; err != nil {
		return nil
	}
	// 返回副本更安全（避免外部修改内存结构），但也会多次分配开销
	copyWorks := make([]model.Work, len(works))
	copy(copyWorks, works)
	return copyWorks
}

// 删除画作
func (d *Database) DelectPaint(username string, workname string) bool {
	// paintMu.Lock()
	// defer paintMu.Unlock()

	// 检查用户是否存在
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	var work model.Work
	if err := d.DB.Where("user_id = ? AND title = ?", user.ID, workname).First(&work).Error; err != nil {
		return false // 作品不存在
	}
	// works, exists := WorksMap[username]
	// if !exists {
	// 	return false
	// }
	if err := d.DB.Delete(&work).Error; err != nil {
		return false
	}
	if work.Image != "" {
		filePath := strings.TrimPrefix(work.Image, "/")
		if strings.HasPrefix(filePath, "uploads/") {
			_ = os.Remove(filePath) // 忽略文件删除错误
		}
	}
	// // 创建一个新的切片来存储保留的作品
	// // newWorks := make([]model.Work, 0, len(works))
	// // found := false

	// // 遍历所有作品
	// for _, work := range works {
	// 	if work.Title == workname && !found {
	// 		// 找到要删除的作品，跳过它
	// 		found = true
	// 		// 删除本地图片文件 - 简化的版本
	// 		if work.Image != "" {
	// 			// 先移除开头的斜杠，然后检查是否为 uploads/ 路径
	// 			filePath := strings.TrimPrefix(work.Image, "/")
	// 			if strings.HasPrefix(filePath, "uploads/") {
	// 				os.Remove(filePath)
	// 			}
	// 		}
	// 	} else {
	// 		// 保留其他作品
	// 		newWorks = append(newWorks, work)
	// 	}
	// }

	// // 如果找到了要删除的作品
	// if found {
	// 	WorksMap[username] = newWorks
	// 	saveLocked()
	// 	return true
	// }

	return true
}

// 添加评论
func (d *Database) AddComment(username string, workname string, comment model.Comment) bool {
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}

	var work model.Work
	if err := d.DB.Where("user_id = ? AND title = ?", user.ID, workname).First(&work).Error; err != nil {
		return false
	}

	comment.WorkID = work.ID
	if err := d.DB.Create(&comment).Error; err != nil {
		return false
	}
	return true
}

// 删除评论
func (d *Database) DelectComment(username string, workname string, comment model.Comment) bool {
	// paintMu.Lock()
	// defer paintMu.Unlock()
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	var work model.Work
	if err := d.DB.Where("user_id = ? AND title = ?", user.ID, workname).First(&work).Error; err != nil {
		return false
	}
	if err := d.DB.Where("work_id = ? AND content = ?", work.ID, comment.Content).First(&comment).Error; err != nil {
		return false
	}
	if err := d.DB.Model(&work).Association("Comments").Delete(&comment); err != nil {
		return false
	}
	return true
	// for num := range WorksMap[username] {
	// 	if WorksMap[username][num].Title == workname {
	// 		// 遍历评论列表，找到匹配的评论并删除
	// 		for commentNum := len(WorksMap[username][num].Comments) - 1; commentNum >= 0; commentNum-- {
	// 			c := WorksMap[username][num].Comments[commentNum]
	// 			// 根据发布者和时间匹配评论（使用 CreatedAt 精确匹配）
	// 			if c.FromUser == comment.FromUser && c.CreatedAt.Equal(comment.CreatedAt) {
	// 				// 使用 append 删除：拼接前半部分和后半部分
	// 				WorksMap[username][num].Comments = append(
	// 					WorksMap[username][num].Comments[:commentNum],
	// 					WorksMap[username][num].Comments[commentNum+1:]...,
	// 				)
	// 				saveLocked()
	// 				return true
	// 			}
	// 		}
	// 	}
	// }
}

//删除画作
// ...existing code... 继续等待
