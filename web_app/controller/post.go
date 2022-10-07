package controller

import (
	"strconv"
	"web_app/logic"
	"web_app/models"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//CreatePostHandler  创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	//1.获取参数及参数的校验

	//c.ShouldBindJSON() //vaildator---->bingd tag
	p := new(models.Post)

	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invaild param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从c取到当前发请求的用户的id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

//CetPostDetailHandler 获取帖子详情
func CetPostDetailHandler(c *gin.Context) {
	//1.获取参数（url中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail	with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id取出帖子数据（查数据库）
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

//CetPostListlHandler //获取帖子列表的处理函数
func CetPostListlHandler(c *gin.Context) {
	page, size := getPageInfo(c)
	//获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

//CetPostListlHandler2 //升级版本的铁子列表接口
//根据前端传来的参数动态的获取帖子列表
//按创建时间排序。或者按照分数排序

//1.获取参数
//2.去redis查询id 列表
//3.根据id去数据库查询帖子详细信息

func CetPostListlHandler2(c *gin.Context) {
	//Get 请求参数 /api/v1/post2?page=1&size=10&order=time query sting
	//初始化结构体指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //magic string
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("CetPostListlHandler2 withinvalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))

		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, data)
}

//根据社区去查询帖子列表

func GetCommunityPostList() {

}
