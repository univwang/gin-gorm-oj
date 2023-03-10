package service

import (
	"bytes"
	"errors"
	"gin_oj/define"
	"gin_oj/helper"
	"gin_oj/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))

	page = (page - 1) * size // 起始位置
	var count int64
	list := make([]models.SubmitBasic, 0)

	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err := tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Submit Error:", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
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

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /me/submit [post]
func Submit(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error:" + err.Error(),
		})
		return
	}
	// 代码保存
	path, err := helper.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Save Error:" + err.Error(),
		})
		return
	}
	u, _ := c.Get("user")
	userClaims := u.(*helper.UserClaims)

	sb := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaims.Identity,
		Path:            path,
	}

	//代码的判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).
		Preload("TestCases").First(&pb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	// 答案错误的channel
	WA := make(chan int)
	// 超内存的channel
	OOM := make(chan int)
	// 编译错误的channel
	CE := make(chan int)
	// 通过的个数
	passCount := 0
	var lock sync.Mutex
	//提示信息
	var msg string

	for _, testCase := range pb.TestCases {
		testCase := testCase
		go func() {
			cmd := exec.Command("go", "run", path)
			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out
			pipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			io.WriteString(pipe, testCase.Input)
			//根据测试的输入案例进行运行，拿到输出结果和标准的输出结果进行比对
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					msg = stderr.String()
					CE <- 1
					return
				}
			}

			var em runtime.MemStats
			runtime.ReadMemStats(&em)
			// 答案错误
			if testCase.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			// 运行超内存
			if em.Alloc/1024-bm.Alloc/1024 > uint64(pb.MaxMem) {
				msg = "运行超内存"
				OOM <- 1
				return
			}
			lock.Lock()
			passCount += 1
			lock.Unlock()
		}()
	}
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		if passCount == len(pb.TestCases) {
			sb.Status = 1
			msg = "通过测试"
		} else {
			sb.Status = 3
			msg = "运行超时"
		}
	case <-CE:
		sb.Status = 5
	}
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = models.DB.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + 1")
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + 1")
		}
		// 更新user_basic
		err = tx.Model(new(models.UserBasic)).
			Where("identity =?", userClaims.Identity).Updates(m).Error
		if err != nil {
			return errors.New("Update UserBasic Error:" + err.Error())
		}
		// 更新problem_basic
		err = tx.Model(new(models.ProblemBasic)).
			Where("identity =?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("Update ProblemBasic Error:" + err.Error())
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}
