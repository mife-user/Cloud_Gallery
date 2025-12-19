package utils

//发放通行证
import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySecret = []byte("lao-tou-de-mi-yao") // 您的私房密钥

// 发证
func GenerateToken(username string) (string, error) { //参数为用户名，为他发放
	claims := jwt.MapClaims{ //原始的钥匙
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 2小时过期（钥匙必须声明过期时间）
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //利用算法包装原始钥匙为未签名的钥匙
	return token.SignedString(mySecret)                        //将包装结果打造成密钥（用我的私房密钥签名）
}

// 提取验证（解析 Token）
func ParseToken(tokenString string) (string, error) { //传入密钥
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { //将密钥转换为*jwt.Token对象，即未签名的钥匙
		return mySecret, nil // 告诉解析器，请用我的私钥去验章
	})
	//转换后实现判断
	if claims /*ok为true则claims赋值为有效的jwt.MapClaims类型*/, ok /*ok为断言结果*/ := token.Claims.(jwt.MapClaims); /*对token.Claims做类型断言，判断是否为jwt.MapClaims类型*/
	ok /*先判断ok再判断token是否过期*/ && token.Valid {
		return claims["username"].(string), nil //成功后返回密钥原本数据中的“username”其中.(string)为断言string类型，因为claims为MapClaims类型，他的值为inerface{}
	}
	return "", err
}
