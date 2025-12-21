package middleware

import (
	"painting/web/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

//用context而不是string或channel的原因是
/*
	1.string为纯文本，而请求链路中有多种类型数据（用户TOKEN，请求参数，响应状态），也无法实现请求
	2.channel控制不了流程，要控制生命周期
	3.context则是一站式数据传递与控制需求
*/
// 接受并处理HTTP请求
func AuthMiddleware() gin.HandlerFunc /*本质就是func(*gin.Context)*/ { // AuthMiddleware(）负责创建func(c *gin.Context)来对HTTP请求进行查验
	return func(c *gin.Context) {
		// 1. 从请求头里拿 Authorization: Bearer xxxxx
		authHeader := c.GetHeader("Authorization") //获取请求头
		if authHeader == "" {                      //为空
			c.AbortWithStatusJSON(401, gin.H{"error": "无证！"}) //立即中止请求链路+返回HTTP状态码与{JSON响应}gin.H类型（map[string]any）
			return
		}

		// 2. 把 Bearer 去掉，只要后面的 token
		parts := strings.SplitN(authHeader, " ", 2) //HTTP有三个要查看的，去除一个还有两个，返回的字符串数组
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "证件格式不对！"}) //立即中止请求链路+返回HTTP状态码与{JSON响应}gin.H类型（map[string]any）
			return
		}

		// 3. 验真伪
		username, err := utils.ParseToken(parts[1]) //第三个为密钥，用utils包中的验证来查看
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "证件过期或造假！"})
			return
		}

		// 4. 重点！把解析出来的用户名，塞进上下文里，传给里面的经理
		c.Set("username", username) //存储键值对数据
		//二者实现了Gin的链路操作
		c.Next() // 放行，将数据给后续的
	}
}

// 1. 外交官：允许跨域
// 浏览器出于安全，默认禁止 JS 跨域名请求，这里是告诉浏览器“我允许”，此处我表示看不懂。。。
/*！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！*/
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 允许任何人访问
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS") // 允许的方法
		// OPTIONS 是浏览器的“预检”请求，问服务器支不支持。
		// 如果是 OPTIONS，直接告诉它支持（204 No Content），不需要再往下处理了。
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

/*！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！*/
