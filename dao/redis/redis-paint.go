package redis

import (
	"fmt"
	"painting/box"
	"painting/model"
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
	pattern := fmt.Sprintf("%s:*", who)
	var works []model.Work
	iter := box.Temp.RE.Scan(c, 0, pattern, 0).Iterator()
	for iter.Next(c) {
		workKey := iter.Val()
		workData, err := box.Temp.RE.HMGet(c, workKey, "title", "author", "content", "image").Result()
		if err != nil {
			c.JSON(400, gin.H{"error": "redis服务器错误"})
			return false
		}
		if len(workData) < 4 || workData[0] == nil || workData[1] == nil || workData[2] == nil || workData[3] == nil {
			continue
		}

		title, ok1 := workData[0].(string)
		author, ok2 := workData[1].(string)
		content, ok3 := workData[2].(string)
		image, ok4 := workData[3].(string)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			continue
		}
		/*--------------------------------------------------------------------------*/
		patternComment := fmt.Sprintf("%s:%s:*:*:comments", who, title)
		var comments []model.Comment
		iterComment := box.Temp.RE.Scan(c, 0, patternComment, 0).Iterator()
		for iterComment.Next(c) {
			commentKey := iterComment.Val()
			commentData, err := box.Temp.RE.HMGet(c, commentKey, "from_user", "content", "created_at").Result()
			if err != nil {
				c.JSON(400, gin.H{"error": "redis服务器错误"})
				return false
			}
			if len(commentData) < 3 || commentData[0] == nil || commentData[1] == nil || commentData[2] == nil {
				continue
			}

			fromUser, okFU := commentData[0].(string)
			contentStr, okCT := commentData[1].(string)
			createdAtStr, okCA := commentData[2].(string)
			if !okFU || !okCT || !okCA {
				continue
			}
			createdAt, err := time.Parse(time.RFC3339, createdAtStr)
			if err != nil {
				continue
			}

			comments = append(comments, model.Comment{
				FromUser:  fromUser,
				Content:   contentStr,
				CreatedAt: createdAt,
			})
		}

		work := model.Work{
			Title:    title,
			Author:   author,
			Content:  content,
			Image:    image,
			Comments: comments,
		}
		works = append(works, work)
	}

	c.JSON(200, gin.H{"owner": who, "works": works})
	return true
}

// 看画缓存处理
func View_write(c *gin.Context, works *[]model.Work) {
	for _, work := range *works {
		userWork := fmt.Sprintf("%s:%s", work.Author, work.Title)
		if err := box.Temp.RE.HMSet(c,
			userWork,
			"title", work.Title,
			"author", work.Author,
			"content", work.Content,
			"image", work.Image,
		).Err(); err != nil {
			c.JSON(400, gin.H{"error": "redis服务器错误"})
			return
		}

		for _, comment := range work.Comments {
			commentKey := fmt.Sprintf("%s:%s:%s:%s:comments",
				work.Author, work.Title, comment.FromUser, comment.CreatedAt.Format(time.RFC3339))

			if err := box.Temp.RE.HMSet(c,
				commentKey,
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
