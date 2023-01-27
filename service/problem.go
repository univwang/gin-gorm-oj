package service

import (
	"encoding/json"
	"errors"
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
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Param max_runtime formData int false "max_runtime"
// @Param max_mem formData int false "max_mem"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
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

// ProblemUpdate
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int true "max_runtime"
// @Param max_mem formData int true "max_mem"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-update [put]
func ProblemUpdate(c *gin.Context) {
	identity := c.PostForm("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	if identity == "" || title == "" || content == "" || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": "-1",
			"msg":  "参数错误",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息保存 problem_basic
		problemBasic := models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMem:     maxMem,
		}
		err := tx.Where("identity =?", identity).Updates(&problemBasic).Error
		if err != nil {
			return err
		}

		// 查询问题详情
		err = tx.Where("identity =?", identity).Find(&problemBasic).Error
		if err != nil {
			return err
		}

		// 关联问题分类的更新
		// 1、删除原来的
		err = tx.Where("problem_id = ?", problemBasic.ID).
			Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		// 2. 新增新的
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			intId, _ := strconv.Atoi(id)
			pcs = append(pcs, &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(intId),
			})
		}
		err = tx.Create(&pcs).Error

		if err != nil {
			return err
		}
		// 关联测试案例的更新
		// 1. 删除已存在的关联关系
		err = tx.Where("problem_identity =?", identity).
			Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		// 2. 增加新的关联关系
		tcs := make([]*models.TestCase, 0)
		for _, testCases := range testCases {
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(testCases), &caseMap)
			if err != nil {
				return err
			}
			if _, ok := caseMap["input"]; !ok {
				return errors.New("测试案例[input]格式错误")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("测试案例[output]格式错误")
			}
			tcs = append(tcs, &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Update Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "问题更新成功",
	})

}
