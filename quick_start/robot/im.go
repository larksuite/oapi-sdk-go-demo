package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	oapi_sdk_go_demo "oapi-sdk-go-demo"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var userOpenIds = []string{"ou_a79a0f82add14976e3943f4deb17c3fa", "ou_33c76a4cbeb76bd66608706edb32508e"}

// ListChatHistory 获取会话历史消息
func ListChatHistory(chatId string) error {
	req := larkim.NewListMessageReqBuilder().
		ContainerIdType("chat").
		ContainerId(chatId).
		Build()
	resp, err := oapi_sdk_go_demo.Client.Im.Message.List(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return resp.CodeError
	}

	pwd, _ := os.Getwd()
	file, err := os.OpenFile(pwd+"/chat_history.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	for _, i := range resp.Data.Items {
		senderId := *i.Sender.Id
		content := *i.Body.Content
		createTime := *i.CreateTime
		intCreateTime, err := strconv.ParseInt(createTime, 10, 64)
		if err != nil {
			continue
		}
		createTime = fmt.Sprintf("%v", time.Unix(intCreateTime/1000, 0))
		str := fmt.Sprintf("chatter(%v) at (%v) send: %v", senderId, createTime, content)
		_, _ = write.WriteString(str + "\n")
		_ = write.Flush()
	}

	return nil
}

// CreateAlertChat 创建报警群并拉人入群
func CreateAlertChat() (string, error) {
	req := larkim.NewCreateChatReqBuilder().
		UserIdType("open_id").
		Body(larkim.NewCreateChatReqBodyBuilder().
			Name("P0: 线上事故处理").
			Description("线上紧急事故处理").
			UserIdList(userOpenIds).
			Build()).
		Build()

	resp, err := oapi_sdk_go_demo.Client.Im.Chat.Create(context.Background(), req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", resp.CodeError
	}

	return *resp.Data.ChatId, nil
}

// SendAlertMessage 发送报警消息
func SendAlertMessage(chatId string) error {
	content, err := buildCard("跟进处理")
	if err != nil {
		return err
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(chatId).
			MsgType("interactive").
			Content(content).
			Build()).
		Build()

	resp, err := oapi_sdk_go_demo.Client.Im.Message.Create(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return resp.CodeError
	}

	return nil
}

// 上传图片
func uploadImage() (string, error) {
	image, err := os.Open("./quick_start/robot/alert.png")
	if err != nil {
		return "", err
	}
	req := larkim.NewCreateImageReqBuilder().
		Body(larkim.NewCreateImageReqBodyBuilder().
			ImageType("message").
			Image(image).
			Build()).
		Build()

	resp, err := oapi_sdk_go_demo.Client.Im.Image.Create(context.Background(), req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", resp.CodeError
	}

	return *resp.Data.ImageKey, nil
}

// 构建卡片
func buildCard(buttonName string) (string, error) {
	imageKey, err := uploadImage()
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadFile("./quick_start/robot/card.json")
	if err != nil {
		return "", err
	}

	card := string(bs)
	card = strings.Replace(card, "${image_key}", imageKey, -1)
	card = strings.Replace(card, "${button_name}", buttonName, -1)
	return card, nil
}

// 获取会话信息
func getChatInfo(chatId string) (*larkim.GetChatRespData, error) {
	req := larkim.NewGetChatReqBuilder().
		ChatId(chatId).
		Build()

	resp, err := oapi_sdk_go_demo.Client.Im.Chat.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if !resp.Success() {
		return nil, resp.CodeError
	}

	return resp.Data, nil
}

// 更新会话名称
func updateChatName(chatId string, chatName string) error {
	req := larkim.NewUpdateChatReqBuilder().
		ChatId(chatId).
		Body(larkim.NewUpdateChatReqBodyBuilder().
			Name(chatName).
			Build()).
		Build()

	resp, err := oapi_sdk_go_demo.Client.Im.Chat.Update(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success() {
		return resp.CodeError
	}

	return nil
}

// DoP2ImMessageReceiveV1 处理消息回调
func DoP2ImMessageReceiveV1(ctx context.Context, data *larkim.P2MessageReceiveV1) error {
	msg := data.Event.Message
	if strings.Contains(*msg.Content, "/solve") {
		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType("chat_id").
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(*msg.ChatId).
				MsgType("text").
				Content("{\"text\":\"问题已解决，辛苦了!\"}").
				Build()).
			Build()

		resp, err := oapi_sdk_go_demo.Client.Im.Message.Create(context.Background(), req)
		if err != nil {
			return err
		}
		if !resp.Success() {
			return resp.CodeError
		}

		// 获取会话信息
		chatInfo, err := getChatInfo(*msg.ChatId)
		if err != nil {
			return err
		}
		name := *chatInfo.Name
		if strings.HasPrefix(name, "[跟进中]") {
			name = "[已解决]" + name[len("[跟进中]"):]
		} else if !strings.HasPrefix(name, "[已解决]") {
			name = "[已解决]" + name
		}
		// 修改会话名称
		err = updateChatName(*msg.ChatId, name)
		if err != nil {
			return err
		}
	}

	return nil
}

// DoInteractiveCard 处理卡片回调
func DoInteractiveCard(ctx context.Context, data *larkcard.CardAction) (interface{}, error) {
	if data.Action.Value["key"] == "follow" {
		chatInfo, err := getChatInfo(data.OpenChatId)
		if err != nil {
			return nil, err
		}
		name := *chatInfo.Name
		if !strings.HasPrefix(name, "[跟进中]") && !strings.HasPrefix(name, "[已解决]") {
			name = "[跟进中] " + name
		}
		// 修改会话名称
		err = updateChatName(data.OpenChatId, name)
		if err != nil {
			return nil, err
		}

		return buildCard("跟进中")
	}

	return nil, nil
}
