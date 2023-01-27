package helper

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

// GetMd5
// 生成MD5
func GetMd5(s string) string {
	//返回16进制的密文
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("gin-gorm-oj-key")

// GeneratorToken
// 生成token
func GeneratorToken(identity string, name string, IsAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        IsAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		log.Println(err.Error())
	}
	return tokenString, nil
}

// AnalyseToken
// 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	UserClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, UserClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error + %v", err)
	}
	return UserClaim, err
}

func SendCode(toUserEmail string, code string) error {
	e := email.NewEmail()
	e.From = "GO-OJ练习平台 <1151695676@qq.com>"
	e.To = []string{toUserEmail}
	e.Subject = "gin-gorm-oj验证码"
	e.HTML = []byte("<h1>您的验证码：<b>" + code + "</b><h1>")
	//xujzslhdysowhgia
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1151695676@qq.com", "xujzslhdysowhgia", "smtp.qq.com"))
	//err := e.SendWithTLS("smtp.qq.com:587",
	//	smtp.PlainAuth("", "2162768982@qq.com", "wqy610366", "smtp.qq.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	return err
}

// GetUUID
// 生成唯一码
func GetUUID() string {
	return uuid.NewV4().String()
}

// GetRand
// 生成验证码
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}
