package redis

import (
	"fmt"
	"painting/box"
	"painting/model"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 挂画缓存处理
func UpPaint_write(c *gin.Context, u *model.Work) {
	userWork := fmt.Sprintf("%s:%s", u.Author, u.Title)
	if err := box.Temp.RE.HMSet(c,
		userWork,
		"title", u.Title,
		"author", u.Author,
		"content", u.Content,
		"image", u.Image,
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
	box.Temp.RE.Expire(c, userWork, 2*time.Hour)
}

// 删画缓存处理
func DePaint_write(c *gin.Context, u *model.Work) {
	userWork := fmt.Sprintf("%s:%s", u.Author, u.Title)
	if err := box.Temp.RE.HDel(c,
		userWork,
		"title",
		"author",
		"content",
		"image",
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

// 看画缓存检查
func View_read(c *gin.Context, who string) bool {
	/*扫描作品*/
	pattern := fmt.Sprintf("%s:*", who)
	var works []model.Work
	iter := box.Temp.RE.Scan(c, 0, pattern, 0).Iterator()
	for iter.Next(c) {
		workKey := iter.Val()
		workData, err := box.Temp.RE.HMGet(c, workKey, "title", "author", "content", "image", "created_at", "updated_at").Result()
		if err != nil {
			c.JSON(400, gin.H{"error": "redis服务器错误"})
			return false
		}
		/*检查空值*/
		if len(workData) < 6 || workData[0] == nil || workData[1] == nil || workData[2] == nil || workData[3] == nil {
			continue
		}
		/*转换类型*/
		title, ok1 := workData[0].(string)
		author, ok2 := workData[1].(string)
		content, ok3 := workData[2].(string)
		image, ok4 := workData[3].(string)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			continue
		}
		/*时间检查*/
		var workCreatedAt, workUpdatedAt time.Time
		if len(workData) >= 6 {
			if workData[4] != nil {
				if createdStr, ok := workData[4].(string); ok && createdStr != "" {
					if t, err := time.Parse(time.RFC3339, createdStr); err == nil {
						workCreatedAt = t
					}
				}
			}
			if workData[5] != nil {
				if updatedStr, ok := workData[5].(string); ok && updatedStr != "" {
					if t, err := time.Parse(time.RFC3339, updatedStr); err == nil {
						workUpdatedAt = t
					}
				}
			}
		}

		/*--------------------------------------------------------------------------*/
		/*扫描评论*/
		patternComment := fmt.Sprintf("%s:%s:*:*:comments", who, title)
		var comments []model.Comment
		iterComment := box.Temp.RE.Scan(c, 0, patternComment, 0).Iterator()
		for iterComment.Next(c) {
			commentKey := iterComment.Val()
			commentData, err := box.Temp.RE.HMGet(c, commentKey, "id", "from_user", "content", "created_at", "updated_at").Result()
			if err != nil {
				c.JSON(400, gin.H{"error": "redis服务器错误"})
				return false
			}
			if len(commentData) < 5 || commentData[0] == nil || commentData[1] == nil || commentData[2] == nil || commentData[3] == nil {
				parts := strings.Split(commentKey, ":")
				if len(parts) >= 5 {
					createdFromKey := parts[3]
					createdTime, err := time.Parse(time.RFC3339, createdFromKey)
					if err != nil {
						continue
					}
					var id uint
					if commentData[0] != nil {
						if s, ok := commentData[0].(string); ok {
							if v, err := strconv.ParseUint(s, 10, 64); err == nil {
								id = uint(v)
							}
						}
					}
					fromUser := ""
					if commentData[1] != nil {
						if s, ok := commentData[1].(string); ok {
							fromUser = s
						}
					}
					contentStr := ""
					if commentData[2] != nil {
						if s, ok := commentData[2].(string); ok {
							contentStr = s
						}
					}
					cm := model.Comment{
						FromUser: fromUser,
						Content:  contentStr,
					}
					cm.ID = id
					cm.CreatedAt = createdTime
					comments = append(comments, cm)
					continue
				}
				continue
			}

			var id uint
			if s, ok := commentData[0].(string); ok {
				if v, err := strconv.ParseUint(s, 10, 64); err == nil {
					id = uint(v)
				}
			}
			fromUser, okFU := commentData[1].(string)
			contentStr, okCT := commentData[2].(string)
			createdStr, okCR := commentData[3].(string)
			if !okFU || !okCT || !okCR {
				continue
			}
			createdTime, err := time.Parse(time.RFC3339, createdStr)
			if err != nil {
				parts := strings.Split(commentKey, ":")
				if len(parts) >= 5 {
					if t, err2 := time.Parse(time.RFC3339, parts[3]); err2 == nil {
						createdTime = t
					}
				}
			}

			cm := model.Comment{
				FromUser: fromUser,
				Content:  contentStr,
			}
			cm.ID = id
			cm.CreatedAt = createdTime
			comments = append(comments, cm)
		}

		work := model.Work{
			Title:    title,
			Author:   author,
			Content:  content,
			Image:    image,
			Comments: comments,
		}
		work.CreatedAt = workCreatedAt
		work.UpdatedAt = workUpdatedAt

		works = append(works, work)
	}

	c.JSON(200, gin.H{"owner": who, "works": works})
	return true
}

// 看画缓存写入
func View_write(c *gin.Context, works *[]model.Work) {
	for _, work := range *works {
		userWork := fmt.Sprintf("%s:%s", work.Author, work.Title)
		if err := box.Temp.RE.HMSet(c,
			userWork,
			"title", work.Title,
			"author", work.Author,
			"content", work.Content,
			"image", work.Image,
			"created_at", work.CreatedAt.Format(time.RFC3339),
			"updated_at", work.UpdatedAt.Format(time.RFC3339),
		).Err(); err != nil {
			c.JSON(400, gin.H{"error": "redis服务器错误"})
			return
		}

		for _, comment := range work.Comments {
			commentKey := fmt.Sprintf("%s:%s:%s:%s:comments",
				work.Author, work.Title, comment.FromUser, comment.CreatedAt.Format(time.RFC3339))

			if err := box.Temp.RE.HMSet(c,
				commentKey,
				"id", strconv.FormatUint(uint64(comment.ID), 10),
				"from_user", comment.FromUser,
				"content", comment.Content,
				"created_at", comment.CreatedAt.Format(time.RFC3339),
			).Err(); err != nil {
				c.JSON(400, gin.H{"error": "redis服务器错误"})
				return
			}
			box.Temp.RE.Expire(c, commentKey, 2*time.Hour)
		}
		box.Temp.RE.Expire(c, userWork, 2*time.Hour)
	}
}
