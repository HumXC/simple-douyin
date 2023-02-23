package douyin

import (
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) MessageAction(c *gin.Context) {
	type Resp struct {
		Response
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()

	messageMan := h.DB.Message
	actionType, _ := strconv.Atoi(c.Query("action_type"))
	toUserId, _ := strconv.Atoi(c.Query("to_user_id"))
	userId := c.GetInt64("user_id")

	//发送消息
	if actionType == 1 {
		content := c.Query("content")
		message := model.Message{
			FromUserId: userId,
			ToUserId:   int64(toUserId),
			Content:    content,
		}
		err := messageMan.AddMessage(&message)
		if err != nil {
			resp.Status(StatusOtherError)
			log.Println("发送消息失败:", err.Error())
			return
		}
	} else {
		resp.Status(UnKnownActionType)
		log.Println("未定义的操作类型")
	}
}

func (h *Handler) MessageChatListAction(c *gin.Context) {

	type Resp struct {
		Response
		MessageList []Message ` json:"message_list,omitempty"`
	}
	resp := Resp{
		Response: BaseResponse(),
	}
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()
	messageMan := h.DB.Message
	fromUserId := c.GetInt64("user_id")
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64) //对方用户ID
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}
	preMsgTime, err := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64) //上次最新消息的时间(第一次请求该值为0)
	if err != nil {
		resp.Status(StatusOtherError)
		return
	}
	var messages []model.Message

	time := time.Unix(preMsgTime, 0).Format("2006-01-02 15:04:05")
	time += ".9999999+08:00"

	err = messageMan.QueryChat(fromUserId, toUserId, time, &messages) //获取时间大于time的聊天记录
	if err != nil {
		resp.Status(StatusFailedChatList)
		return
	}

	if len(messages) > 0 {
		messageList := make([]Message, len(messages))
		convert(messages, messageList)
		resp.MessageList = messageList
		return
	}
	return
}

// 将model.message切片转为Message切片
func convert(messages []model.Message, messageList []Message) {
	idx := 0
	for _, message := range messages {
		messageData := Message{
			Id:         int64(message.Model.ID),
			ToUserId:   message.ToUserId,
			FromUserId: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.CreatedAt.Unix(),
		}
		messageList[idx] = messageData
		idx = idx + 1
	}

}
