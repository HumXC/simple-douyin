package service

import (
	"fmt"
	"time"

	"github.com/HumXC/simple-douyin/config"
	"github.com/HumXC/simple-douyin/handler/douyin"
	"github.com/HumXC/simple-douyin/middlewares"
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/videos"
	"github.com/gin-gonic/gin"
)

type DouYin struct {
	engine *gin.Engine
}

func NewDouyin(g *gin.Engine, conf config.Douyin, db *model.DouyinDB, storageClient douyin.StorageClient) *DouYin { // 初始化 douyin
	handler := douyin.Handler{
		DB:            db,
		StorageClient: storageClient,
		VideoButcher: videos.NewButcher(db.VideoJob, conf.VideoButCherMaxJob, func(job videos.Job, video, cover string, err error) (delete bool) {
			if err != nil {
				fmt.Println("视频任务失败: " + err.Error())
				return false
			}
			vHash, err := storageClient.Upload(video, "videos")
			if err != nil {
				fmt.Println("视频任务失败: " + err.Error())
				return false
			}
			cHash, err := storageClient.Upload(cover, "covers")
			if err != nil {
				fmt.Println("视频任务失败: " + err.Error())
				return false
			}
			// 将视频信息写入数据库
			_ = db.Video.Put(model.Video{
				Video:  vHash,
				Cover:  cHash,
				Title:  job.Title,
				UserID: job.UserID,
				Time:   time.Now(),
			})
			return true
		}),
	}
	douyin := g.Group("douyin")
	douyin.GET("feed", handler.Feed(conf.FeedNum))
	douyin.POST("user/register/", middlewares.PwdHashMiddleWare(), handler.UserRegister)
	douyin.POST("user/login/", handler.UserLogin)
	douyin.GET("user/", middlewares.JWTMiddleWare(), handler.User)
	douyin.POST("comment/action/", middlewares.JWTMiddleWare(), handler.CommentAction)
	douyin.GET("comment/list/", handler.CommentList)

	publish := douyin.Group("publish")
	publish.Use(middlewares.JWTMiddleWare())
	publish.POST("action/", handler.PublishAction)
	publish.GET("list/", handler.PublishList)

	relation := douyin.Group("relation")
	relation.POST("action/", middlewares.JWTMiddleWare(), handler.RelationAction)
	relation.GET("/follow/list/", middlewares.AuthUserCheck(), handler.FollowList)
	relation.GET("/follower/list/", middlewares.AuthUserCheck(), handler.FollowerList)
	relation.GET("/friend/list/", middlewares.AuthUserCheck(), handler.FriendList)
	
	return &DouYin{
		engine: g,
	}
}
