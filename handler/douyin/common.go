package douyin

import (
	"github.com/HumXC/simple-douyin/model"
	"github.com/HumXC/simple-douyin/videos"
)

// 状态码定义，所有状态码应当在写在 StatusOK 与 StatusOtherError 之间
// 并且新定义的状态码不需要赋值
// 例如：
// StatusOK = iota
// StatusNotFound
// StatusOtherError = -1
// 此时 StatusNotFound 的值为 1，不要关心值是多少
const (
	StatusOK         = iota
	StatusAuthFailed = iota + 400
	StatusAuthKeyTimeout
	StatusUserNotFound
	StatusUploadNotAVideo
	StatusFailedToFetchVideo
	StatusNeedLogin
	StatusOtherError = -1
)

var StatusMsgs = make(map[int32]string, 16)

func init() {
	StatusMsgs[StatusOK] = "OK"
	StatusMsgs[StatusOtherError] = "其他错误"
	StatusMsgs[StatusAuthFailed] = "身份验证失败"
	StatusMsgs[StatusAuthKeyTimeout] = "登录信息已过期"
	StatusMsgs[StatusUserNotFound] = "用户不存在"
	StatusMsgs[StatusUploadNotAVideo] = "上传的文件不是视频"
	StatusMsgs[StatusFailedToFetchVideo] = "获取视频列表失败"
	StatusMsgs[StatusNeedLogin] = "需要登录"
}

// 所有 gin.HandlerFunc 都应该绑定到 Handler 上
type Handler struct {
	DB            *model.DouyinDB
	StorageClient StorageClient
	VideoButcher  *videos.Butcher
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 返回一个状态码 200 的 Response
func BaseResponse() Response {
	return Response{
		StatusCode: StatusOK,
		StatusMsg:  StatusMsgs[StatusOK],
	}
}

// 设定 Response 的状态码，如果状态消息在 StatusMsgs 中找不到，则状态消息会被设定为 "未定义状态"
func (r *Response) Status(code int32) {
	r.StatusCode = code
	r.StatusMsg = StatusMsgs[code]
	if r.StatusMsg == "" {
		r.StatusMsg = "未定义状态"
	}
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
