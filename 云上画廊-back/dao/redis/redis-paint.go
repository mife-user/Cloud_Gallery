package redis

import (
	"fmt"
	"painting/box"
	"painting/model"
	"strconv"
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
		"created_at", u.CreatedAt.Format(time.RFC3339),
		"updated_at", u.CreatedAt.Format(time.RFC3339),
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
		"created_at",
		"updated_at",
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

// 看画缓存检查
func View_read(c *gin.Context, who string) bool {
	works, ok := getWorks(c, who)
	if !ok {
		return false
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
				"updated_at", comment.UpdatedAt.Format(time.RFC3339),
			).Err(); err != nil {
				c.JSON(400, gin.H{"error": "redis服务器错误"})
				return
			}
			box.Temp.RE.Expire(c, commentKey, 2*time.Hour)
		}
		box.Temp.RE.Expire(c, userWork, 2*time.Hour)
	}
}
