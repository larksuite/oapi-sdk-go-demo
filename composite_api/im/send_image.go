/*
 发送图片消息，使用到两个OpenAPI：
 1. [上传图片](https://open.feishu.cn/document/server-docs/im-v1/image/create)
 2. [发送消息](https://open.feishu.cn/document/server-docs/im-v1/message/create)
*/

package im

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type SendImageRequest struct {
	Image         io.Reader
	ReceiveIdType string
	ReceiveId     string
	Uuid          string
}

type SendImageResponse struct {
	*larkcore.CodeError
	CreateImageResponse   *larkim.CreateImageRespData
	CreateMessageResponse *larkim.CreateMessageRespData
}

// SendImage 发送图片消息
func SendImage(client *lark.Client, request *SendImageRequest) (*SendImageResponse, error) {
	// 上传图片
	createImageReq := larkim.NewCreateImageReqBuilder().
		Body(larkim.NewCreateImageReqBodyBuilder().
			ImageType("message").
			Image(request.Image).
			Build()).
		Build()
	createImageResp, err := client.Im.Image.Create(context.Background(), createImageReq)
	if err != nil {
		return nil, err
	}
	if !createImageResp.Success() {
		fmt.Printf("client.Im.Image.Create failed, code: %d, msg: %s, log_id: %s\n",
			createImageResp.Code, createImageResp.Msg, createImageResp.RequestId())
		return nil, createImageResp.CodeError
	}

	// 发送消息
	bs, err := json.Marshal(createImageResp.Data)
	if err != nil {
		return nil, err
	}
	createMessageReq := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(request.ReceiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(request.ReceiveId).
			MsgType("image").
			Content(string(bs)).
			Uuid(request.Uuid).
			Build()).
		Build()

	createMessageResp, err := client.Im.Message.Create(context.Background(), createMessageReq)
	if err != nil {
		return nil, err
	}
	if !createMessageResp.Success() {
		fmt.Printf("client.Im.Message.Create failed, code: %d, msg: %s, log_id: %s\n",
			createMessageResp.Code, createMessageResp.Msg, createMessageResp.RequestId())
		return nil, createMessageResp.CodeError
	}

	// 返回结果
	return &SendImageResponse{
		CodeError: &larkcore.CodeError{
			Code: 0,
			Msg:  "success",
		},
		CreateImageResponse:   createImageResp.Data,
		CreateMessageResponse: createMessageResp.Data,
	}, nil
}
