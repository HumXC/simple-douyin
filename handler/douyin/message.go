package douyin

import (
	"github.com/HumXC/simple-douyin/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
	toUserId, _ := strconv.Atoi(c.Query("to_user_id")) //对方用户ID
	var messages []model.Message

	var messages1 []model.Message
	var messages2 []model.Message

	err := messageMan.QueryMessageRecord(fromUserId, int64(toUserId), &messages1)
	err = messageMan.QueryMessageRecord(int64(toUserId), fromUserId, &messages2)

	messages = append(messages1, messages2...)

	if err != nil {
		resp.Status(StatusFailedChatList)
		log.Println("拉取聊天记录失败", err.Error())
		return
	}
	messageList := make([]Message, len(messages))
	idx := 0
	for _, message := range messages {
		messageData := Message{
			Id:         int64(message.Model.ID),
			ToUserId:   message.ToUserId,
			FromUserId: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		messageList[idx] = messageData
		idx = idx + 1
	}
	resp.MessageList = messageList
}
