package app

import (
	"chats/controller"
	"chats/middleware"
)

func Routers() {
	r.POST("/users/chatcreate", middleware.CORSMiddleware(), controller.CreateChat)
	r.POST("/users/chatrequest", middleware.CORSMiddleware(), controller.RequestChat)
	r.POST("/users/reqviews", middleware.CORSMiddleware(), controller.ViewsReq)
	r.POST("/users/acceptchat", middleware.CORSMiddleware(), controller.AcceptChat)
	r.POST("/users/reqacceptviews", middleware.CORSMiddleware(), controller.ViewsAcceptReq)
	r.POST("/users/chatting", middleware.CORSMiddleware(), controller.Chat)
	r.POST("/users/deletechat", middleware.CORSMiddleware(), controller.DeleteChat)
	r.POST("/users/allviews", middleware.CORSMiddleware(), controller.AllViews)
	r.POST("/users/locking", middleware.CORSMiddleware(), controller.LockingChat)
	r.POST("/users/unlocking", middleware.CORSMiddleware(), controller.UnLockingChat)
	r.POST("/users/userDetails", middleware.CORSMiddleware(), controller.UserDetails)
	r.POST("/users/viewmessage", middleware.CORSMiddleware(), controller.UserMessageDetails)
	r.POST("/users/chats", middleware.CORSMiddleware(), controller.Chats)
	r.POST("/users/allchats", middleware.CORSMiddleware(), controller.AllChats)
	r.POST("/users/findchats", middleware.CORSMiddleware(), controller.FindChats)

}
