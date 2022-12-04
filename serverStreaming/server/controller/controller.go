package controller

import (
	"context"
	"fmt"
	pb "serverstream/proto"
	"serverstream/server/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userController struct{}

func NewUserControllerServer() pb.UserServiceServer {
	return userController{}
}

var users *mongo.Collection = database.OpenCollection(database.Client, "users")
var chatting *mongo.Collection = database.OpenCollection(database.Client, "chatting")

func (userController) Chatting(in *pb.ProtoChatUserRequest, stream pb.UserService_ChattingServer) error {

	waitc := make(chan struct{})
	data := &pb.ChatProto{}
	go func() {
		for {

			senderEmail := in.ProtoSenderEmail
			receiverEmail := in.ProtoReceiverEmail
			fmt.Println("senderEmail", senderEmail)
			fmt.Println("receiverEmail", receiverEmail)

			opts := options.FindOne().SetSort(bson.M{"$natural": -1})
			filter := bson.M{"$and": []interface{}{bson.M{"$or": []bson.M{{"protosenderemail": senderEmail}, {"protosenderemail": receiverEmail}}}, bson.M{"$or": []bson.M{{"protoreceiveremail": senderEmail}, {"protoreceiveremail": receiverEmail}}}}}
			cursorerr := chatting.FindOne(context.Background(), filter, opts).Decode(data)

			if cursorerr != nil {
				stream.Send(&pb.ProtoChatUserReponse{Res: nil})
			}
			response := &pb.ChatProto{
				ProtoSenderUser:    data.ProtoSenderUser,
				ProtoSenderEmail:   senderEmail,
				ProtoMessage_:      data.ProtoMessage_,
				ProtoReceiverUser:  data.ProtoReceiverUser,
				ProtoReceiverEmail: receiverEmail,
				ProtoTime:          data.ProtoTime,
				ProtoDate:          data.ProtoDate,
			}

			stream.Send(&pb.ProtoChatUserReponse{Res: response})
			time.Sleep(1 * time.Second)
		}
	}()
	<-waitc
	return nil
}
