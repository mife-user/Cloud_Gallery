package db

import (
	"fmt"
	"os"
	"painting/model"
	"strings"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
	RE *redis.Client
}

// 关闭数据库
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// 初始化
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

	dsn := fmt.Sprintf("root:123456@tcp(%s:%s)/paint?charset=utf8mb4&parseTime=True&loc=Local", host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, false
	}
	db.AutoMigrate(&model.User{}, &model.Work{}, &model.Comment{})
	re := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
	if re.Ping(re.Context()).Err() != nil {
		return nil, false
	}
	DB := &Database{DB: db, RE: re}
	return DB, true
}

// 实现添加作品逻辑
func (d *Database) AddWork(username string, work *model.Work) bool {
	var user model.User
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

	if err := d.DB.Create(&user).Error; err != nil {
		return false
	}
	return true

}

// 实现获取work逻辑，应该也为项目特殊
func (d *Database) GetWorks(username string) ([]model.Work, error) {
	// 这里直接返回切片引用（注意：调用方不要直接修改返回切片）
	// 若要防止外部篡改，可返回副本
	// works := WorksMap[username]
	user := model.User{
		Username: username,
		Password: "",
	}
	if err := d.DB.Model(&user).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	var works []model.Work
	if err := d.DB.Preload("Comments").Where("user_id = ?", user.ID).Find(&works).Error; err != nil {
		return nil, err
	}
	// 返回副本更安全（避免外部修改内存结构），但也会多次分配开销
	copyWorks := make([]model.Work, len(works))
	copy(copyWorks, works)
	return copyWorks, nil
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
func (d *Database) AddComment(username string, workname string, comment *model.Comment) bool {
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}

	var work model.Work
	if err := d.DB.Where("user_id = ? AND title = ?", user.ID, workname).First(&work).Error; err != nil {
		return false
	}

	comment.WorkID = work.ID
	if err := d.DB.Create(comment).Error; err != nil {
		return false
	}
	return true
}

// 删除评论
func (d *Database) DelectComment(username string, workname string, comment *model.Comment) bool {
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	var work model.Work
	if err := d.DB.Where("user_id = ? AND title = ?", user.ID, workname).First(&work).Error; err != nil {
		return false
	}
	if err := d.DB.Where("work_id = ? AND from_user = ? AND created_at = ?",
		work.ID, comment.FromUser, comment.CreatedAt).First(comment).Error; err != nil {
		return false
	}
	if err := d.DB.Model(&work).Association("Comments").Delete(comment); err != nil {
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

func (d *Database) AddHand(user *model.User) bool {
	var userTemp model.User
	if err := d.DB.Where("username = ?", user.Username).First(&userTemp).Error; err != nil {
		return false
	}
	if err := d.DB.Model(&userTemp).Update("user_hand", user.UserHand).Error; err != nil {
		return false
	}
	return true
}
