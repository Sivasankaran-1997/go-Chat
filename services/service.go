package services

import (
	"chats/domain"
	"chats/utils"
	"time"

	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateChat(user domain.User) (*mongo.InsertOneResult, *utils.Resterr) {
	if err := user.Vaildate(); err != nil {
		return nil, err
	}
	guid := xid.New()
	user.ID = guid.String()
	user.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
	userResult, restErr := user.Create()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func RequestChat(request domain.UserRequest) (*mongo.InsertOneResult, *utils.Resterr) {
	if err := request.RequestVaildate(); err != nil {
		return nil, err
	}
	guid := xid.New()
	request.ID = guid.String()
	request.SenderStatus = true
	request.RequestTime = time.Now().Format("2006-01-02 15:04:05")
	userRequest, restErr := request.ReqChat()
	if restErr != nil {
		return nil, restErr
	}
	return userRequest, nil
}

func AcceptChat(request domain.UserRequest) (*mongo.UpdateResult, *utils.Resterr) {
	if err := request.AcceptVaildate(); err != nil {
		return nil, err
	}
	userRequest, restErr := request.AcceptChat()
	if restErr != nil {
		return nil, restErr
	}
	return userRequest, nil
}

func RequestView(request domain.UserRequest) ([]domain.UserRequest, *utils.Resterr) {
	if err := request.ReqViewVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.ReqViews()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func AcceptReq(request domain.UserRequest) ([]domain.UserRequest, *utils.Resterr) {
	if err := request.ReqViewVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.AcceptReq()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func Chat(request domain.UserRequest) ([]domain.UserRequest, *utils.Resterr) {
	if err := request.ReqViewVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.Chat()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func DeleteChat(request domain.UserRequest) (*mongo.DeleteResult, *utils.Resterr) {
	if err := request.AcceptVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.DeleteChat()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func AllView(request domain.UserRequest) ([]domain.User, *utils.Resterr) {
	if err := request.ReqViewVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.AllViews()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func Locking(request domain.UserRequest) (*mongo.UpdateResult, *utils.Resterr) {
	if err := request.LockingVaildate(); err != nil {
		return nil, err
	}
	userRequest, restErr := request.Locking()
	if restErr != nil {
		return nil, restErr
	}
	return userRequest, nil
}

func UnLocking(request domain.UserRequest) (*mongo.UpdateResult, *utils.Resterr) {
	if err := request.LockingVaildate(); err != nil {
		return nil, err
	}
	userRequest, restErr := request.UnLocking()
	if restErr != nil {
		return nil, restErr
	}
	return userRequest, nil
}

func UserDetails(request domain.UserRequest) ([]domain.User, *utils.Resterr) {
	if err := request.ReqViewVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.UserDetails()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func UserMessageDetails(request domain.ChatProto, startDate string, endDate string, email string) ([]domain.ChatProto, *utils.Resterr) {
	if err := request.MessageVaildate(startDate, endDate, email); err != nil {
		return nil, err
	}
	userResult, restErr := request.UserMessageDetails(startDate, endDate, email)
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func Chats(chat domain.ChatProto) (*mongo.InsertOneResult, *utils.Resterr) {
	if err := chat.ChatsVaildate(); err != nil {
		return nil, err
	}
	chat.Prototime = time.Now().Format("2006-01-02 15:04:05")
	chat.Protodate = time.Now().Format("2006-01-02")
	receivername := domain.GetReceivername(chat.Protoreceiveremail)
	chat.Protoreceiveruser = receivername
	userResult, restErr := chat.Chats()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func AllChats(request domain.ChatProto) ([]domain.ChatProto, *utils.Resterr) {
	if err := request.AllChatsVaildate(); err != nil {
		return nil, err
	}
	userResult, restErr := request.AllChats()
	if restErr != nil {
		return nil, restErr
	}
	return userResult, nil
}

func FindChats(request domain.ChatProto) (*domain.ChatProto, *utils.Resterr) {
	if err := request.AllChatsVaildate(); err != nil {
		return nil, err
	}
	restErr := request.FindChats()
	if restErr != nil {
		return nil, restErr
	}
	return &request, nil
}
