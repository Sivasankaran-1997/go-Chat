package controller

import (
	"chats/domain"
	"chats/services"
	"chats/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateChat(c *gin.Context) {
	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.CreateChat(newUser)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func ViewsReq(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.RequestView(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func Chat(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.Chat(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func DeleteChat(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.DeleteChat(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func ViewsAcceptReq(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.AcceptReq(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func RequestChat(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.RequestChat(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func AcceptChat(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.AcceptChat(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func AllViews(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.AllView(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func LockingChat(c *gin.Context) {
	var newUser domain.UserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.Locking(newUser)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func UnLockingChat(c *gin.Context) {
	var newUser domain.UserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.UnLocking(newUser)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func UserDetails(c *gin.Context) {
	var newRequest domain.UserRequest

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.UserDetails(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)

}

func UserMessageDetails(c *gin.Context) {
	var newRequest domain.ChatProto

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	email := c.Query("email")

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.UserMessageDetails(newRequest, startDate, endDate, email)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func Chats(c *gin.Context) {
	var newRequest domain.ChatProto

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.Chats(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func AllChats(c *gin.Context) {
	var newRequest domain.ChatProto

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.AllChats(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func FindChats(c *gin.Context) {
	var newRequest domain.ChatProto

	if err := c.ShouldBindJSON(&newRequest); err != nil {
		resterr := utils.BadRequest("Invalid JSON")
		c.JSON(resterr.Status, resterr)
		return
	}

	result, resterr := services.FindChats(newRequest)

	if resterr != nil {
		c.JSON(resterr.Status, resterr)
		return
	}
	c.JSON(http.StatusOK, result)
}
