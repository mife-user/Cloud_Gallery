package api

import (
	"fmt"
	"painting/model"
	"painting/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册缓存检查
func register1(c *gin.Context, u *model.User) bool {
	if err := Temp.RE.HGet(c, u.Username, "username").Err(); err == nil {
		c.JSON(400, gin.H{"error": "用户已存在"})
		return false
	}
	return true
}

// 注册缓存设置
func register2(c *gin.Context, u *model.User) {
	if err := Temp.RE.HSet(c, u.Username, "username", u.Username).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
	Temp.RE.Expire(c, u.Username, 2*time.Hour)
}

// 登录缓存检查与处理
func login1(c *gin.Context, u *model.User) int {
	userTemp, err := Temp.RE.HMGet(c, u.Username, "username", "password").Result()
	if err != nil {
		return 0
	}
	if userTemp[0] == nil || userTemp[1] == nil {
		return 0
	}
	userName, ok1 := userTemp[0].(string)
	passWord, ok2 := userTemp[1].(string)
	if !ok1 || !ok2 {
		return 0
	}
	if ok := Temp.CheckUser(userName, passWord); ok {
		token, _ := utils.GenerateToken(userName)
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		return 2
	} else {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return 1
	}
}

// 登录缓存设置与处理
func login2(c *gin.Context, u *model.User) {
	if err := Temp.RE.HMSet(c, u.Username,
		"username", u.Username,
		"password", u.Password,
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
	Temp.RE.Expire(c, u.Username, 2*time.Hour)
}

// 挂画缓存处理
func uploadPaint2(c *gin.Context, u *model.Work) {
	userWork := fmt.Sprintf("%s:%s", u.Author, u.Title)
	if err := Temp.RE.HMSet(c,
		userWork,
		"title", u.Title,
		"author", u.Author,
		"content", u.Content,
		"image", u.Image,
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
	Temp.RE.Expire(c, userWork, 2*time.Hour)
}

// 删画缓存处理
func delectPaint2(c *gin.Context, u *model.Work) {
	userWork := fmt.Sprintf("%s:%s", u.Author, u.Title)
	if err := Temp.RE.HDel(c,
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
func view1(c *gin.Context, who string) bool {
	pattern := fmt.Sprintf("%s:*", who)
	var works []model.Work
	iter := Temp.RE.Scan(c, 0, pattern, 0).Iterator()
	for iter.Next(c) {
		workKey := iter.Val()
		workData, err := Temp.RE.HMGet(c, workKey, "title", "author", "content", "image").Result()
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
		iterComment := Temp.RE.Scan(c, 0, patternComment, 0).Iterator()
		for iterComment.Next(c) {
			commentKey := iterComment.Val()
			commentData, err := Temp.RE.HMGet(c, commentKey, "from_user", "content", "created_at").Result()
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
func view2(c *gin.Context, works *[]model.Work) {
	for _, work := range *works {
		userWork := fmt.Sprintf("%s:%s", work.Author, work.Title)
		if err := Temp.RE.HMSet(c,
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

			if err := Temp.RE.HMSet(c,
				commentKey,
				"from_user", comment.FromUser,
				"content", comment.Content,
				"created_at", comment.CreatedAt.Format(time.RFC3339),
			).Err(); err != nil {
				c.JSON(400, gin.H{"error": "redis服务器错误"})
				return
			}
			Temp.RE.Expire(c, commentKey, 2*time.Hour)
		}

		Temp.RE.Expire(c, userWork, 2*time.Hour)
	}
}

// 评论缓存处理
func postComment1(c *gin.Context, newComment *model.Comment, req *model.CommentRequest) bool {

	workComment := fmt.Sprintf("%s:%s:%s:%s:comments",
		req.TargetAuthor, req.WorkTitle, newComment.FromUser, newComment.CreatedAt.Format(time.RFC3339))

	if err := Temp.RE.HMSet(c,
		workComment,
		"from_user", newComment.FromUser,
		"content", newComment.Content,
		"created_at", newComment.CreatedAt.Format(time.RFC3339),
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return false
	}
	Temp.RE.Expire(c, workComment, 2*time.Hour)
	return true
}

// 作者评论删除缓存处理
func delectCommentMaster1(c *gin.Context, currentMaster string, req *model.DeleteCommentReq) {
	workComment := fmt.Sprintf("%s:%s:%s:%s:comments", currentMaster, req.Title, req.FromUser, req.CreatedAt.Format(time.RFC3339))
	if err := Temp.RE.HDel(c,
		workComment,
		"from_user",
		"content",
		"created_at",
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

// 用户评论删除缓存处理
func delectCommentPoster1(c *gin.Context, req *model.DeleteCommentReq) {
	workComment := fmt.Sprintf("%s:%s:%s:%s:comments", req.Owner, req.Title, req.FromUser, req.CreatedAt.Format(time.RFC3339))
	if err := Temp.RE.HDel(c,
		workComment,
		"from_user",
		"content",
		"created_at",
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

// 添加头像缓存处理
func addUserHand2(c *gin.Context, user *model.User) {
	if err := Temp.RE.HSet(c, user.Username, "userhand", user.UserHand).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}
