package redis

import (
	"fmt"
	"painting/box"
	"painting/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*扫描作品*/
func getWorks(c *gin.Context, works *[]model.Work, who string) bool {
	pattern := fmt.Sprintf("%s:*", who)
	iter := box.Temp.RE.Scan(c, 0, pattern, 0).Iterator()
	for iter.Next(c) {
		workKey := iter.Val()
		workData, err := box.Temp.RE.HMGet(c, workKey, "title", "author", "content", "image", "created_at", "updated_at").Result()
		/*处理错误*/
		if err != nil {
			c.JSON(400, gin.H{"error": "redis服务器错误"})
			return false
		}
		/*空值检查*/
		if workData[0] == nil || workData[1] == nil || workData[2] == nil || workData[3] == nil || workData[4] == nil || workData[5] == nil {
			continue
		}
		/*转换类型*/
		title, okTitle := workData[0].(string)
		author, okAuthor := workData[1].(string)
		content, okContent := workData[2].(string)
		image, okImage := workData[3].(string)
		created_at, okCreated := workData[4].(string)
		updated_at, okUpdated := workData[5].(string)
		/*错误检查*/
		if !okTitle || !okAuthor || !okContent || !okImage || !okCreated || !okUpdated {
			continue
		}
		createdTime, _ := time.Parse(time.RFC3339, created_at)
		updatedTime, _ := time.Parse(time.RFC3339, updated_at)

		/*--------------------------------------------------------------------------*/
		/*扫描评论*/
		var comments []model.Comment
		patternComment := fmt.Sprintf("%s:%s:*:*:comments", who, title)
		iterComment := box.Temp.RE.Scan(c, 0, patternComment, 0).Iterator()
		for iterComment.Next(c) {
			commentKey := iterComment.Val()
			commentData, err := box.Temp.RE.HMGet(c, commentKey, "id", "from_user", "content", "created_at", "updated_at").Result()
			/*处理错误*/
			if err != nil {
				c.JSON(400, gin.H{"error": "redis服务器错误"})
				return false
			}
			/*错误处理*/
			if len(commentData) < 5 || commentData[0] == nil || commentData[1] == nil || commentData[2] == nil || commentData[3] == nil || commentData[4] == nil {
				continue
			}
			/*转换类型*/
			idC, okIdC := commentData[0].(string)
			from_userC, okFromC := commentData[1].(string)
			contentC, okContentC := commentData[2].(string)
			created_atC, okCreatedC := commentData[3].(string)
			updated_atC, okUpdatedC := commentData[4].(string)
			if !okIdC || !okFromC || !okContentC || !okCreatedC || !okUpdatedC {
				continue
			}
			createdTimeComment, _ := time.Parse(time.RFC3339, created_atC)
			updatedTimeComment, _ := time.Parse(time.RFC3339, updated_atC)
			idCommentUint, _ := strconv.ParseUint(idC, 10, 64)
			/*加入数组*/
			cm := model.Comment{
				FromUser: from_userC,
				Content:  contentC,
			}
			cm.ID = uint(idCommentUint)
			cm.CreatedAt = createdTimeComment
			cm.UpdatedAt = updatedTimeComment
			comments = append(comments, cm)
		}

		work := model.Work{
			Title:    title,
			Author:   author,
			Content:  content,
			Image:    image,
			Comments: comments,
		}
		work.CreatedAt = createdTime
		work.UpdatedAt = updatedTime

		*works = append(*works, work)
	}
	return true
}
