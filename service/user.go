package service

import (
	"fmt"
	"gin_oj/define"
	"gin_oj/helper"
	"gin_oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":"200","msg":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}
	date := new(models.UserBasic)
	err := models.DB.Omit("password").Where("identity = ?", identity).Find(&date).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Detail By Identity:" + identity + " Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  date,
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "username"
// @Param password formData string false "password"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /user-login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	//123456 -> e10adc3949ba59abbe56e057f20f883e
	//md5 存储到数据库中的是密文
	fmt.Print(password)
	password = helper.GetMd5(password)
	fmt.Println(" -> " + password)
	data := new(models.UserBasic)
	err := models.DB.Where("name = ? AND password = ?", username, password).
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": "-1",
				"msg":  "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "登录出错:" + err.Error(),
		})
		return
	}
	token, err := helper.GeneratorToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "GeneratorToken Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string false "email"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /user-send_code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	code := helper.GetRand()
	models.RDB.Set(c, email, code, time.Second*300)
	err := helper.SendCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "Send Code Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "Success",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param name formData string true "name"
// @Param mail formData string true "mail"
// @Param password formData string true "password"
// @Param phone formData string false "phone"
// @Param code formData string true "code"
// @Success 200 {string} json "{"code":"200","msg":"","data":""}"
// @Router /user-register [post]
func Register(c *gin.Context) {
	mail := c.PostForm("mail")
	name := c.PostForm("name")
	password := c.PostForm("password")
	userCode := c.PostForm("code")
	phone := c.PostForm("phone")

	if mail == "" || name == "" || password == "" || userCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "信息不能为空",
		})
		return
	}
	//验证码是否正确
	sysCode, err := models.RDB.Get(c, mail).Result()
	if err != nil {
		log.Printf("Get Code Error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "验证码错误",
		})
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "验证码错误",
		})
		return
	}
	//判断邮箱是否已经存在
	var cnt int64
	err = models.DB.Where("mail = ?", mail).Model(new(models.UserBasic)).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Error:" + err.Error(),
		})
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已经注册",
		})
		return
	}
	//数据的插入

	data := &models.UserBasic{
		Identity: helper.GetUUID(),
		Name:     name,
		Password: helper.GetMd5(password),
		Phone:    phone,
		Mail:     mail,
	}

	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create User Error" + err.Error(),
		})
		return
	}

	//生成token
	token, err := helper.GeneratorToken(data.Identity, data.Name, data.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Token Error" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	//page 1 --> 0

	page = (page - 1) * size // 起始位置
	var count int64

	list := make([]*models.UserBasic, 0)
	err := models.DB.Model(new(models.UserBasic)).Count(&count).Order("finish_problem_num DESC, submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "Get Rank List Error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
