package db

import (
	"os"
	"painting/model"
	"strings"
)

// 实现获取work逻辑，应该也为项目特殊
func (d *Database) GetWorks(username string) ([]model.Work, error) {
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
	return works, nil
}

// 删除画作
func (d *Database) DelectPaint(username string, workname string) error {
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	var work model.Work
	if err := d.DB.Where("user_id = ? AND title = ?", user.ID, workname).First(&work).Error; err != nil {
		return err
	}
	if err := d.DB.Where("work_id = ?", work.ID).Delete(&model.Comment{}).Error; err != nil {
		return err
	}
	if err := d.DB.Delete(&work).Error; err != nil {
		return err
	}
	if work.Image != "" {
		filePath := strings.TrimPrefix(work.Image, "/")
		if strings.HasPrefix(filePath, "uploads/") {
			err := os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 实现添加作品逻辑
func (d *Database) AddWork(username string, work *model.Work) error {
	var user model.User
	if err := d.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}
	if err := d.DB.Model(&user).Association("Work").Append(work); err != nil {
		return err
	}
	return nil
}
