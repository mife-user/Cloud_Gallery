package db

import (
	"painting/model"
)

// 查人（给登录用）
func (d *Database) CheckUser(username, password string) bool {
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
	user := model.User{
		Username: username,
		Password: password,
	}

	if err := d.DB.Create(&user).Error; err != nil {
		return false
	}
	return true

}

// 添加用户头像
func (d *Database) AddHand(user *model.User) bool {
	var userTemp model.User
	if err := d.DB.Where("username = ?", user.Username).First(&userTemp).Error; err != nil {
		return false
	}
	if err := d.DB.Model(&userTemp).Update("userhand", user.UserHand).Error; err != nil {
		return false
	}
	return true
}
