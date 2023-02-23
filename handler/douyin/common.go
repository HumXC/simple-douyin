package douyin

import (
	"github.com/HumXC/simple-douyin/handler/douyin/videos"
	"github.com/HumXC/simple-douyin/model"
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
	InvalidParams
	UnKnownActionType
	StatusFailedPostComment
	StatusCommentNotFound
	StatusFailedDelComment
	StatusFailedCommentList
	StatusVideoHasNoComment
	StatusFailedChatList
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
	StatusMsgs[InvalidParams] = "参数错误"
	StatusMsgs[UnKnownActionType] = "未知操作类型"
	StatusMsgs[StatusFailedPostComment] = "发布评论失败"
	StatusMsgs[StatusCommentNotFound] = "未找到该评论"
	StatusMsgs[StatusFailedDelComment] = "删除评论失败"
	StatusMsgs[StatusFailedCommentList] = "拉取评论列表失败"
	StatusMsgs[StatusVideoHasNoComment] = "该视频暂无评论"
	StatusMsgs[StatusFailedChatList] = "拉取聊天记录失败"
}

// 所有 gin.HandlerFunc 都应该绑定到 Handler 上
type Handler struct {
	DB            *DBMan
	RDB           *RDBMan
	StorageClient StorageClient
	VideoButcher  *videos.Butcher
	Avatars       []string
	Backgrounds   []string
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
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count,omitempty"`    // 关注数
	FollowerCount  int64  `json:"follower_count,omitempty"`  // 粉丝数
	TotalFavorited int64  `json:"total_favorited,omitempty"` // 获赞数
	WorkCount      int64  `json:"work_count,omitempty"`      // 作品数量
	FavoriteCount  int64  `json:"favorite_count,omitempty"`  // 点赞数
	IsFollow       bool   `json:"is_follow,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
	Background     string `json:"background_image,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
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

// 转换数据库 model.User 到 douyin.User
func (h *Handler) ConvertUser(u model.User, isFollow bool) User {
	return User{
		Id:            u.ID,
		FollowCount:   h.DB.User.CountFollow(u.ID),
		FollowerCount: h.DB.User.CountFollower(u.ID),
		IsFollow:      isFollow,
		Name:          u.Name,
		Avatar:        h.StorageClient.GetURL("avatars", u.Avatar),
		Background:    h.StorageClient.GetURL("backgrounds", u.Background),
		WorkCount:     h.DB.User.CountPublished(u.ID),
		FavoriteCount: h.RDB.User.CountFavorite(u.ID),
	}
}

// 转换数据库 model.User 到 douyin.User
func (h *Handler) ConvertUsers(us *[]model.User, isFollow bool) *[]User {
	result := make([]User, len(*us), len(*us))
	for i := 0; i < len(*us); i++ {
		result[i] = h.ConvertUser((*us)[i], isFollow)
	}
	return &result
}

// 转换视频，user 是请求对象的 userID，用于获取 IsFavorite 等字段
// author 是视频发布者，未知时传入 nil
func (h *Handler) ConvertVideo(src model.Video, user int64, author *User) Video {
	v := Video{}
	_u := model.User{}
	if author == nil {
		err := h.DB.User.QueryById(src.UserID, &_u)
		if err != nil {
			panic(err)
		}
		a := h.ConvertUser(_u, h.DB.User.IsFollow(user, _u.ID))
		author = &a
	}
	v.Author = *author
	v.CommentCount = src.CommentCount
	v.FavoriteCount = h.RDB.Video.CountFavorite(src.ID)
	if user != 0 {
		v.IsFavorite = h.RDB.User.IsFavorite(user, src.ID)
	}
	v.Id = src.ID
	v.CoverUrl = h.StorageClient.GetURL("covers", src.Cover)
	v.PlayUrl = h.StorageClient.GetURL("videos", src.Video)
	return v
}

// 转换视频，user 是请求对象的 userID，用于获取 IsFavorite 等字段
// author 是视频发布者，未知时传入 nil
func (h *Handler) ConvertVideos(vs *[]model.Video, user int64, author *User) *[]Video {
	result := make([]Video, len(*vs), len(*vs))
	for i := 0; i < len(*vs); i++ {
		result[i] = h.ConvertVideo((*vs)[i], user, author)
	}
	return &result
}
