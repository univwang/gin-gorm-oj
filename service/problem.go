package service

import (
	"encoding/json"
	"gin_oj/define"
	"gin_oj/helper"
	"gin_oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Param category_identity query string false "category_identity"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data/":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")
	//page 1 --> 0

	page = (page - 1) * size // 起始位置
	var count int64
	list := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
	err := tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println(err)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "问题标识不能为空",
		})
		return
	}
	date := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").First(&date).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": "-1",
				"msg":  "问题不存在",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": date,
	})
}

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 问题创建
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-create [post]
func ProblemCreate(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMemory, _ := strconv.Atoi(c.PostForm("max_mem"))
	if title == "" || content == "" || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "参数错误",
		})
		return
	}
	identity := helper.GetUUID()
	data := &models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMemory,
	}
	// 处理分类
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryId, _ := strconv.Atoi(id)
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(categoryId),
		})
	}
	data.ProblemCategories = categoryBasics
	//处理测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range testCases {
		// {"input":"1 2\n", "output":"3\n"}
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "-1",
				"msg":  "JSON unmarshal failed",
			})
			return
		}
		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": "-1",
				"msg":  "input is missing",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": "-1",
				"msg":  "output is missing",
			})
			return
		}
		testCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	data.TestCases = testCaseBasics

	// 创建问题
	err := models.DB.Debug().Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Create Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
	})
}
