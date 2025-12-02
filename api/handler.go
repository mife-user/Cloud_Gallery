package api

import (
	"painting/dao"
	"painting/model"
	"painting/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义前端传来的参数结构（DTO）
type CommentRequest struct {
	TargetAuthor string `json:"target_author"` // 评论谁的画
	WorkTitle    string `json:"work_title"`    // 哪幅画
	Content      string `json:"content"`       // 评论内容
}

// 注册
func Register(c *gin.Context) {
	// ShouldBindJSON 自动把前端传来的 JSON 对照着 model.User 的结构填进去
	// 如果格式不对（比如传了 int 而不是 string），err 就不为空
	var u model.User                             //临时创建一个结构（model中的User）
	if err := c.ShouldBindJSON(&u); err != nil { //将context的`json:"username"`与`json:"password"`放入u中
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	if dao.AddUser(u.Username, u.Password) {
		c.JSON(200, gin.H{"message": "注册成功！"})
	} else {
		c.JSON(400, gin.H{"error": "用户名已存在"})
	}
}

// 登录
func Login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	if dao.CheckUser(u.Username, u.Password) {
		token, _ := utils.GenerateToken(u.Username)
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token, // 把证件发给用户
		})
	} else {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
	}
}

// 挂画（需要登录）
func UploadPaint(c *gin.Context) {
	// 从保安那儿知道是谁
	// 从上下文里拿用户名
	// 这个 "username" 是 AuthMiddleware 里用 c.Set 塞进去的
	username, _ := c.Get("username")
	var w model.Work
	if err := c.ShouldBindJSON(&w); /*向临时w写入JSON中的model.Work成功返回nil*/ err != nil {
		c.JSON(400, gin.H{"error": "画的信息没填对"})
		return
	}

	dao.AddWork(username.(string), w)
	c.JSON(200, gin.H{"message": "上传成功", "work": w})
}

// 删画
func Delect(c *gin.Context) {
	username, err1 := c.Get("username") //获取当前登录用户的username
	if !err1 {
		c.JSON(400, gin.H{"error": "用户名不存在"})
		return
	}
	name, ok1 := username.(string)
	if !ok1 {
		c.JSON(400, gin.H{"error": "用户名格式错误"})
		return
	}
	var w model.Work
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(400, gin.H{"error": "作品名不存在"})
		return
	}
	last := dao.DelectPaint(name, w.Title)
	if !last {
		c.JSON(400, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(200, gin.H{"error": "删除成功"})
}

// 看展（公开）
func View(c *gin.Context) {
	who := c.Param("who") //获取路径中的而非上下文中的
	works := dao.GetWorks(who)
	c.JSON(200, gin.H{"owner": who, "works": works})
}

// 评论（需要登录）
func PostComment(c *gin.Context) {
	//当前登录者
	username, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "登录名错误"})
		return
	}
	commentator := username.(string)
	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数没填对"})
		return
	}
	newComment := model.Comment{
		FromUser:  commentator,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	succeces := dao.AddComment(req.TargetAuthor, req.WorkTitle, newComment)
	if succeces {
		c.JSON(200, gin.H{"message": "评论成功", "data": newComment})
	} else {
		c.JSON(404, gin.H{"error": "找不到这幅画"})
	}
}
