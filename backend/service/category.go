package service

import (
	"gin_oj/define"
	"gin_oj/helper"
	"gin_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param authorization header string true "authorization"
// @Param page query int false "请输入当前页，默认第一页"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data/":""}"
// @Router /admin/category-list [get]
func GetCategoryList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	keyword := c.Query("keyword")
	//page 1 --> 0

	page = (page - 1) * size // 起始位置
	var count int64

	categoryList := make([]*models.CategoryBasic, 0)
	err := models.DB.Model(new(models.CategoryBasic)).
		Where("name like ?", "%"+keyword+"%").
		Count(&count).Limit(size).Offset(page).Find(&categoryList).Error

	if err != nil {
		log.Println("GetCategoryList err:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"list":  categoryList,
		},
	})
}

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData int false "parent_id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-create [post]
func CategoryCreate(c *gin.Context) {
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	category := &models.CategoryBasic{
		Name:     name,
		ParentId: parentId,
		Identity: helper.GetUUID(),
	}
	err := models.DB.Create(category).Error
	if err != nil {
		log.Println("CategoryCreate err:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分类创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类创建成功",
	})
}

// CategoryUpdate
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData int false "parent_id"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-update [put]
func CategoryUpdate(c *gin.Context) {
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	if name == "" || identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}

	category := &models.CategoryBasic{
		Name:     name,
		Identity: identity,
		ParentId: parentId,
	}
	err := models.DB.Model(new(models.CategoryBasic)).
		Where("identity = ?", identity).
		Updates(category).Error
	if err != nil {
		log.Println("CategoryUpdate err:", err)

		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分类修改失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分类修改成功",
	})

}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-delete [delete]
func CategoryDelete(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}
	var cnt int64
	err := models.DB.Model(new(models.ProblemCategory)).
		Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)", identity).
		Count(&cnt).Error
	if err != nil {
		log.Println("Get ProblemCategory err:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下存在问题，不可删除",
		})
		return
	}
	err = models.DB.Where("identity =?", identity).Delete(new(models.CategoryBasic)).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
