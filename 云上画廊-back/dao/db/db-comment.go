package db

import "painting/model"

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
	if err := d.DB.Where("work_id = ? AND from_user = ? AND content = ?",
		work.ID, comment.FromUser, comment.Content).Delete(comment).Error; err != nil {
		return false
	}
	return true
}
