package domain

import (
	"chats/utils"
	"context"
	"time"

	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID          string `json:id"`
	Name        string `json:name bson:"name,omitempty"`
	Email       string `json:email bson:"email,omitempty"`
	Lock        bool   `json:lock bson:"lock,omitempty"`
	CreatedTime string `json:createdtime bson:"createdtime,omitempty"`
	UpdateTime  string `json:updatetime bson:"updatetime,omitempty"`
	Admin       bool   `json:admin bson:"admin,omitempty"`
}

type UserRequest struct {
	ID             string `json:id"`
	SenderName     string `json:sendername bson:"sendername,omitempty"`
	SenderEmail    string `json:senderemail bson:"senderemail,omitempty"`
	SenderStatus   bool   `json:senderstatus bson:"senderstatus,omitempty"`
	RequestTime    string `json:requesttime bson:"requesttime,omitempty"`
	ReceiverName   string `json:receivername bson:"receivername,omitempty"`
	ReceiverEmail  string `json:receiveremail bson:"receiveremail,omitempty"`
	ReceiverStatus bool   `json:receiverstatus bson:"receiverstatus,omitempty"`
	AcceptTime     string `json:accepttime bson:"accepttime,omitempty"`
}

type ChatProto struct {
	Protosenderuser    string `json:protosenderuser bson:"protosenderuser,omitempty"`
	Protosenderemail   string `json:protosenderemail bson:"protosenderemail,omitempty"`
	Protoreceiveruser  string `json:protoreceiveruser bson:"protoreceiveruser,omitempty"`
	Protoreceiveremail string `json:protoreceiveremail bson:"protoreceiveremail,omitempty"`
	Protomessage_      string `json:protomessage_ bson:"protomessage_,omitempty"`
	Prototime          string `json:prototime bson:"prototime,omitempty"`
	Protodate          string `json:protodate bson:"protodate,omitempty"`
}

type UserJWTsigneDetails struct {
	Email string
	jwt.RegisteredClaims
}

func (user *User) Vaildate() *utils.Resterr {
	if strings.TrimSpace(user.Name) == "" {
		return utils.BadRequest("Name Required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return utils.BadRequest("Email Required")
	}

	return nil
}

func (request *UserRequest) RequestVaildate() *utils.Resterr {

	if strings.TrimSpace(request.SenderName) == "" {
		return utils.BadRequest("SenderName Required")
	}

	if strings.TrimSpace(request.SenderEmail) == "" {
		return utils.BadRequest("SenderEmail Required")
	}

	if request.SenderEmail == request.ReceiverEmail {
		return utils.BadRequest("Sender Email and Receiver Email is Not Same")
	}

	if strings.TrimSpace(request.ReceiverName) == "" {
		return utils.BadRequest("ReceiverName Required")
	}

	if strings.TrimSpace(request.ReceiverEmail) == "" {
		return utils.BadRequest("ReceiverEmail Required")
	}

	return nil
}

func (request *UserRequest) ReqViewVaildate() *utils.Resterr {

	if strings.TrimSpace(request.SenderEmail) == "" {
		return utils.BadRequest("SenderEmail Required")
	}
	return nil
}

func (request *UserRequest) AcceptVaildate() *utils.Resterr {
	if strings.TrimSpace(request.SenderEmail) == "" {
		return utils.BadRequest("SenderEmail Required")
	}
	if strings.TrimSpace(request.ReceiverEmail) == "" {
		return utils.BadRequest("ReceiverEmail Required")
	}

	if request.SenderEmail == request.ReceiverEmail {
		return utils.BadRequest("Sender Email and Receiver Email is Not Same")
	}

	return nil
}

func (user *UserRequest) LockingVaildate() *utils.Resterr {
	if strings.TrimSpace(user.SenderEmail) == "" {
		return utils.BadRequest("Sender Email Required")
	}

	if strings.TrimSpace(user.ReceiverEmail) == "" {
		return utils.BadRequest("Receiver Email Required")
	}

	if user.SenderEmail == user.ReceiverEmail {
		return utils.BadRequest("Sender Email and Receiver Email is Not Same")
	}

	return nil
}

func (user *ChatProto) MessageVaildate(startDate string, endDate string, email string) *utils.Resterr {
	if strings.TrimSpace(user.Protosenderemail) == "" {
		return utils.BadRequest("Sender Email Required")
	}

	if strings.TrimSpace(user.Protoreceiveremail) == "" {
		return utils.BadRequest("Receiver Email Required")
	}

	if strings.TrimSpace(email) == "" {
		return utils.BadRequest("Admin Email Required")
	}

	if user.Protosenderemail == user.Protoreceiveremail {
		return utils.BadRequest("Sender Email and Receiver Email is Not Same")
	}

	timeValue := strings.Compare(startDate, endDate)
	if timeValue == 1 {
		return utils.BadRequest("Date Value is Invalid")
	}

	return nil
}

func LockedVaildate(senderEmail string, receiverEmail string) *utils.Resterr {
	usersC := DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if strings.TrimSpace(receiverEmail) == "" {
		senderemailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": senderEmail, "lock": false})

		if senderemailCount == 0 {
			return utils.BadRequest("Sender Email is Locked")
		}
	} else {
		senderemailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": senderEmail, "lock": false})

		if senderemailCount == 0 {
			return utils.BadRequest("Sender Email is Locked")
		}

		receiveremailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": receiverEmail, "lock": false})

		if receiveremailCount == 0 {
			return utils.BadRequest("Receiver Email is Locked")
		}
	}
	return nil
}

func AdminVaildate(senderEmail string) *utils.Resterr {
	usersC := DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	adminemailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": senderEmail, "admin": true})

	if adminemailCount == 0 {
		return utils.BadRequest("Admin Only Accesible")
	}
	return nil
}

func (chat *ChatProto) ChatsVaildate() *utils.Resterr {
	if strings.TrimSpace(chat.Protosenderuser) == "" {
		return utils.BadRequest("Sender User Required")
	}

	if strings.TrimSpace(chat.Protosenderemail) == "" {
		return utils.BadRequest("Sender Email Required")
	}

	if strings.TrimSpace(chat.Protoreceiveremail) == "" {
		return utils.BadRequest("Receiver Email Required")
	}

	if strings.TrimSpace(chat.Protomessage_) == "" {
		return utils.BadRequest("Message Required")
	}

	if chat.Protosenderemail == chat.Protoreceiveremail {
		return utils.BadRequest("Sender Email and Receiver Email is Not Same")
	}

	return nil
}

func (chat *ChatProto) AllChatsVaildate() *utils.Resterr {

	if strings.TrimSpace(chat.Protosenderemail) == "" {
		return utils.BadRequest("Sender Email Required")
	}

	if strings.TrimSpace(chat.Protoreceiveremail) == "" {
		return utils.BadRequest("Receiver Email Required")
	}

	if chat.Protosenderemail == chat.Protoreceiveremail {
		return utils.BadRequest("Sender Email and Receiver Email is Not Same")
	}

	return nil
}

func GetReceivername(receiver string) string {

	var user User
	usersC := DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	filter := bson.M{"email": receiver}

	filtererr := usersC.FindOne(ctx, filter).Decode(&user)

	if filtererr != nil {
		return ""
	}
	return user.Name
}
