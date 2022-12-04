package domain

import (
	"chats/utils"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (user *User) Create() (*mongo.InsertOneResult, *utils.Resterr) {
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	emailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if emailCount > 0 {
		return nil, utils.BadRequest("Email Already Register")
	}

	insertResult, err1 := usersC.InsertOne(ctx, user)

	if err1 != nil {
		restErr := utils.InternalErr("can't insert user to the database.")
		return nil, restErr
	}

	return insertResult, nil

}

func (request *UserRequest) ReqViews() ([]UserRequest, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	emailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})
	if emailCount == 0 {
		return nil, utils.BadRequest("Email Not Found")
	}

	emaillockerr := LockedVaildate(request.SenderEmail, "")

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	defer cancel()

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"senderemail": request.SenderEmail}, bson.M{"senderstatus": bson.M{"$eq": true}}, bson.M{"receiverstatus": bson.M{"$eq": false}}}}}}

	cursor, err1 := RequestD.Aggregate(context.Background(), pipeline)
	if err1 != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []UserRequest
	if err := cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Datas Not Found")
		return nil, restErr
	}

	if len(results) == 0 {
		return nil, utils.NotFound("Chat Request Not Found")
	}

	return results, nil
}

func (request *UserRequest) ReqChat() (*mongo.InsertOneResult, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	sendermailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})
	defer cancel()

	if sendermailCount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receivermailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.ReceiverEmail})

	if receivermailCount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	emaillockerr := LockedVaildate(request.SenderEmail, request.ReceiverEmail)

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"$or": []bson.M{{"senderemail": request.SenderEmail}, {"senderemail": request.ReceiverEmail}}}, bson.M{"$or": []bson.M{{"receiveremail": request.ReceiverEmail}, {"receiveremail": request.SenderEmail}}}}}}}
	cursor, err := RequestD.Aggregate(context.TODO(), pipeline)
	if err != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}

	if len(results) > 0 {
		restErr := utils.BadRequest("Already Request Sent")
		return nil, restErr
	}

	result, err1 := RequestD.InsertOne(ctx, request)

	if err1 != nil {
		restErr := utils.InternalErr("can't insert user to the database.")
		return nil, restErr
	}

	return result, nil

}

func (request *UserRequest) AcceptChat() (*mongo.UpdateResult, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()

	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.ReceiverEmail})

	if receiveremailcount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	emaillockerr := LockedVaildate(request.SenderEmail, request.ReceiverEmail)

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"senderstatus": bson.M{"$eq": true}}, bson.M{"receiveremail": request.SenderEmail}, bson.M{"senderemail": request.ReceiverEmail}, bson.M{"receiverstatus": bson.M{"$eq": false}}}}}}

	cursor, err := RequestD.Aggregate(context.TODO(), pipeline)
	if err != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}

	if len(results) == 0 {
		restErr := utils.BadRequest("Request Already Accepted")
		return nil, restErr
	}
	accepttimevalue := time.Now().Format("2006-01-02 15:04:05")

	filter := bson.M{"$and": []interface{}{bson.M{"$or": []bson.M{{"senderemail": request.SenderEmail}, {"senderemail": request.ReceiverEmail}}}, bson.M{"$or": []bson.M{{"receiveremail": request.SenderEmail}, {"receiveremail": request.ReceiverEmail}}}}}
	updateValue := bson.M{"$set": bson.M{"receiverstatus": true, "accepttime": accepttimevalue}}
	opts := options.Update().SetUpsert(true)
	result, errs := RequestD.UpdateOne(ctx, filter, updateValue, opts)
	if result.ModifiedCount == 0 {
		return nil, utils.BadRequest("not modified")
	}

	if errs != nil {
		return nil, utils.InternalErr("Data not Updated")
	}

	return result, nil
}

func (request *UserRequest) AcceptReq() ([]UserRequest, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	emailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})
	defer cancel()

	if emailCount == 0 {
		return nil, utils.NotFound("Email Not Found")
	}

	emaillockerr := LockedVaildate(request.SenderEmail, "")

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"receiveremail": request.SenderEmail}, bson.M{"senderstatus": bson.M{"$eq": true}}, bson.M{"receiverstatus": bson.M{"$eq": false}}}}}}

	cursor, err1 := RequestD.Aggregate(context.Background(), pipeline)
	if err1 != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []UserRequest
	if err := cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Datas Not Found")
		return nil, restErr
	}

	if len(results) == 0 {
		return nil, utils.NotFound("Accept Chat Not Found")
	}

	return results, nil
}

func (request *UserRequest) Chat() ([]UserRequest, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	emailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})
	defer cancel()

	if emailCount == 0 {
		return nil, utils.NotFound("Email Not Found")
	}

	emaillockerr := LockedVaildate(request.SenderEmail, "")

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"$or": []bson.M{{"senderemail": request.SenderEmail}, {"receiveremail": request.SenderEmail}}}, bson.M{"senderstatus": bson.M{"$eq": true}}, bson.M{"receiverstatus": bson.M{"$eq": true}}}}}}

	cursor, err1 := RequestD.Aggregate(context.Background(), pipeline)
	if err1 != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []UserRequest
	if err := cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Datas Not Found")
		return nil, restErr
	}

	if len(results) == 0 {
		return nil, utils.NotFound("Chatting Not Found")
	}

	return results, nil
}

func (request *UserRequest) DeleteChat() (*mongo.DeleteResult, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.ReceiverEmail})

	if receiveremailcount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	emaillockerr := LockedVaildate(request.SenderEmail, request.ReceiverEmail)

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"$or": []bson.M{{"senderemail": request.SenderEmail}, {"senderemail": request.ReceiverEmail}}}, bson.M{"$or": []bson.M{{"receiveremail": request.ReceiverEmail}, {"receiveremail": request.SenderEmail}}}, bson.M{"senderstatus": true}, bson.M{"$or": []bson.M{{"receiverstatus": false}, {"receiverstatus": true}}}}}}}

	cursor, err1 := RequestD.Aggregate(context.Background(), pipeline)
	if err1 != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []UserRequest
	if err := cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Datas Not Found")
		return nil, restErr
	}

	if len(results) == 0 {
		return nil, utils.NotFound("Not Found")
	}

	filter := bson.M{"$and": []interface{}{bson.M{"$or": []bson.M{{"senderemail": request.SenderEmail}, {"senderemail": request.ReceiverEmail}}}, bson.M{"$or": []bson.M{{"receiveremail": request.SenderEmail}, {"receiveremail": request.ReceiverEmail}}}}}

	result, errs := RequestD.DeleteOne(ctx, filter)
	if result.DeletedCount == 0 {
		return nil, utils.BadRequest("No Record Found")
	}

	if errs != nil {
		return nil, utils.NotFound("Email is Not Found")
	}

	return result, nil

}

func (request *UserRequest) AllViews() ([]User, *utils.Resterr) {
	RequestD := DB.Collection("RequestChat")
	usersC := DB.Collection("users")
	var value []string
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	emailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": request.SenderEmail})

	if emailCount == 0 {
		return nil, utils.NotFound("Email Not Found")
	}

	defer cancel()

	emaillockerr := LockedVaildate(request.SenderEmail, "")

	if emaillockerr != nil {
		return nil, emaillockerr
	}

	pipeline := []bson.M{bson.M{"$match": bson.M{"$and": []bson.M{bson.M{"$or": []bson.M{{"senderemail": request.SenderEmail}, {"receiveremail": request.SenderEmail}}}, bson.M{"senderstatus": bson.M{"$eq": true}}, bson.M{"$or": []bson.M{{"receiverstatus": bson.M{"$eq": true}}, {"receiverstatus": bson.M{"$eq": false}}}}}}}}

	cursor, err1 := RequestD.Aggregate(context.Background(), pipeline)
	if err1 != nil {
		restErr := utils.NotFound("Data Not Found")
		return nil, restErr
	}
	var results []UserRequest
	if err := cursor.All(context.TODO(), &results); err != nil {
		restErr := utils.NotFound("Datas Not Found")
		return nil, restErr
	}

	if len(results) == 0 {
		filter := bson.M{"email": bson.M{"$ne": request.SenderEmail}}
		cursor, _ := usersC.Find(ctx, filter)
		var results []User
		if err := cursor.All(context.TODO(), &results); err != nil {
			return nil, utils.NotFound("Users is Not Found")
		}
		return results, nil
	}
	for i := 0; i < len(results); i++ {
		if strings.TrimSpace(results[i].ReceiverEmail) == strings.TrimSpace(request.SenderEmail) {
			value = append(value, results[i].SenderEmail)
		} else {
			value = append(value, results[i].ReceiverEmail)
		}

	}

	filter := bson.M{"email": bson.M{"$nin": value, "$ne": request.SenderEmail}}
	cursor1, _ := usersC.Find(ctx, filter)
	var result []User
	if err := cursor1.All(context.TODO(), &result); err != nil {
		return nil, utils.NotFound("Users is Not Found")
	}

	return result, nil
}

func (user *UserRequest) Locking() (*mongo.UpdateResult, *utils.Resterr) {
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.SenderEmail})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.ReceiverEmail})

	if receiveremailcount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	adminvalidate := AdminVaildate(user.SenderEmail)

	if adminvalidate != nil {
		return nil, adminvalidate
	}

	filter := bson.M{"$and": []interface{}{bson.M{"email": user.ReceiverEmail, "admin": bson.M{"$ne": true}}}}
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	updateValue := bson.M{"$set": bson.M{"lock": true, "updatetime": updateTime}}

	opts := options.Update().SetUpsert(true)

	result, err := usersC.UpdateOne(ctx, filter, updateValue, opts)

	defer cancel()

	if result.ModifiedCount == 0 {
		return nil, utils.BadRequest("not modified")
	}

	if err != nil {
		return nil, utils.InternalErr("Data not Updated")
	}

	return result, nil
}

func (user *UserRequest) UnLocking() (*mongo.UpdateResult, *utils.Resterr) {
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.SenderEmail})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.ReceiverEmail})

	if receiveremailcount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	adminvalidate := AdminVaildate(user.SenderEmail)

	if adminvalidate != nil {
		return nil, adminvalidate
	}

	filter := bson.M{"email": user.ReceiverEmail}
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	updateValue := bson.M{"$set": bson.M{"lock": false, "updatetime": updateTime}}

	opts := options.Update().SetUpsert(true)

	result, err := usersC.UpdateOne(ctx, filter, updateValue, opts)

	defer cancel()

	if result.ModifiedCount == 0 {
		return nil, utils.BadRequest("not modified")
	}

	if err != nil {
		return nil, utils.InternalErr("Data not Updated")
	}

	return result, nil
}

func (user *UserRequest) UserDetails() ([]User, *utils.Resterr) {
	usersC := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.SenderEmail})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	adminvalidate := AdminVaildate(user.SenderEmail)

	if adminvalidate != nil {
		return nil, adminvalidate
	}

	filter := bson.M{"admin": bson.M{"$ne": true}}
	cursor1, _ := usersC.Find(ctx, filter)
	var result []User
	if err := cursor1.All(context.TODO(), &result); err != nil {
		return nil, utils.NotFound("Users is Not Found")
	}

	return result, nil
}

func (chat *ChatProto) UserMessageDetails(startDate string, endDate string, email string) ([]ChatProto, *utils.Resterr) {
	usersC := DB.Collection("users")
	chattingD := DB.Collection("chatting")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protosenderemail})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protoreceiveremail})

	if receiveremailcount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	adminvalidate := AdminVaildate(email)

	if adminvalidate != nil {
		return nil, adminvalidate
	}

	filter := bson.M{"$and": []interface{}{bson.M{"$or": []bson.M{{"protosenderemail": chat.Protosenderemail}, {"protosenderemail": chat.Protoreceiveremail}}}, bson.M{"$or": []bson.M{{"protoreceiveremail": chat.Protosenderemail}, {"protoreceiveremail": chat.Protoreceiveremail}}}, bson.M{"protodate": bson.M{"$gte": startDate}}, bson.M{"protodate": bson.M{"$lte": endDate}}}}
	opts := options.Find().SetSort(bson.D{{"protodate", 1}})
	cursor1, _ := chattingD.Find(ctx, filter, opts)
	var result []ChatProto
	if err := cursor1.All(context.TODO(), &result); err != nil {
		return nil, utils.NotFound("Users is Not Found")
	}

	return result, nil
}

func (chat *ChatProto) Chats() (*mongo.InsertOneResult, *utils.Resterr) {
	usersC := DB.Collection("users")
	chatD := DB.Collection("chatting")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	senderemailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protosenderemail, "lock": false})
	defer cancel()
	if senderemailCount == 0 {
		return nil, utils.BadRequest("Sender Email Not Found")
	}

	receiveremailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protoreceiveremail, "lock": false})
	if receiveremailCount == 0 {
		return nil, utils.BadRequest("Receiver Email Not Found")
	}

	insertResult, err1 := chatD.InsertOne(ctx, chat)

	if err1 != nil {
		restErr := utils.InternalErr("can't insert user to the database.")
		return nil, restErr
	}

	return insertResult, nil

}

func (chat *ChatProto) AllChats() ([]ChatProto, *utils.Resterr) {
	usersC := DB.Collection("users")
	chattingD := DB.Collection("chatting")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protosenderemail, "lock": false})

	if senderemailcount == 0 {
		return nil, utils.NotFound("Sender Email Not Found")
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protoreceiveremail, "lock": false})

	if receiveremailcount == 0 {
		return nil, utils.NotFound("Receiver Email Not Found")
	}

	filter := bson.M{"$and": []interface{}{bson.M{"$or": []bson.M{{"protosenderemail": chat.Protosenderemail}, {"protosenderemail": chat.Protoreceiveremail}}}, bson.M{"$or": []bson.M{{"protoreceiveremail": chat.Protosenderemail}, {"protoreceiveremail": chat.Protoreceiveremail}}}}}
	opts := options.Find().SetSort(bson.D{{"protodate", -1}})
	cursor1, _ := chattingD.Find(ctx, filter, opts)
	var result []ChatProto
	if err := cursor1.All(context.TODO(), &result); err != nil {
		return nil, utils.NotFound("Users is Not Found")
	}

	return result, nil
}

func (chat *ChatProto) FindChats() *utils.Resterr {
	usersC := DB.Collection("users")
	chattingD := DB.Collection("chatting")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	defer cancel()
	senderemailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protosenderemail, "lock": false})

	if senderemailcount == 0 {
		return nil
	}

	receiveremailcount, _ := usersC.CountDocuments(ctx, bson.M{"email": chat.Protoreceiveremail, "lock": false})

	if receiveremailcount == 0 {
		return nil
	}

	filter := bson.M{"$and": []interface{}{bson.M{"$or": []bson.M{{"protosenderemail": chat.Protosenderemail}, {"protosenderemail": chat.Protoreceiveremail}}}, bson.M{"$or": []bson.M{{"protoreceiveremail": chat.Protosenderemail}, {"protoreceiveremail": chat.Protoreceiveremail}}}}}
	//opts := options.FindOne().SetSort(bson.D{{"protodate", -1}})
	cursor1err := chattingD.FindOne(ctx, filter).Decode(&chat)

	if cursor1err != nil {
		return utils.NotFound("Users is Not Found")
	}

	return nil
}
