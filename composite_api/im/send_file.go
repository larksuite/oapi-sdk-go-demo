/*
发送文件消息，使用到两个OpenAPI：
1. [上传文件](https://open.feishu.cn/document/server-docs/im-v1/file/create)
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

type SendFileRequest struct {
	FileType      string
	FileName      string
	File          io.Reader
	Duration      int
	ReceiveIdType string
	ReceiveId     string
	Uuid          string
}

type SendFileResponse struct {
	*larkcore.CodeError
	CreateFileResponse    *larkim.CreateFileRespData
	CreateMessageResponse *larkim.CreateMessageRespData
}

func SendFile(client *lark.Client, request *SendFileRequest) (*SendFileResponse, error) {
	// 上传文件
	createFileReq := larkim.NewCreateFileReqBuilder().
		Body(larkim.NewCreateFileReqBodyBuilder().
			FileType(request.FileType).
			FileName(request.FileName).
			Duration(request.Duration).
			File(request.File).
			Build()).
		Build()
	createFileResp, err := client.Im.File.Create(context.Background(), createFileReq)
	if err != nil {
		return nil, err
	}
	if !createFileResp.Success() {
		fmt.Printf("client.Im.File.Create failed, code: %d, msg: %s, log_id: %s\n",
			createFileResp.Code, createFileResp.Msg, createFileResp.RequestId())
		return nil, createFileResp.CodeError
	}

	// 发送消息
	bs, err := json.Marshal(createFileResp.Data)
	if err != nil {
		return nil, err
	}
	createMessageReq := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(request.ReceiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(request.ReceiveId).
			MsgType("file").
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
	return &SendFileResponse{
		CodeError: &larkcore.CodeError{
			Code: 0,
			Msg:  "success",
		},
		CreateFileResponse:    createFileResp.Data,
		CreateMessageResponse: createMessageResp.Data,
	}, nil
}
