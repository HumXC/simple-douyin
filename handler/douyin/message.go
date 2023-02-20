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

}
