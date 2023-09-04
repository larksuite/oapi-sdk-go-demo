/*
 创建多维表格同时添加数据表，使用到两个OpenAPI：
 1. [创建多维表格](https://open.feishu.cn/document/server-docs/docs/bitable-v1/app/create)
 2. [新增一个数据表](https://open.feishu.cn/document/server-docs/docs/bitable-v1/app-table/create)
*/

package base

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

type CreateAppAndTablesRequest struct {
	Name        string
	FolderToken string
	Tables      []*larkbitable.ReqTable
}

type CreateAppAndTablesResponse struct {
	*larkcore.CodeError
	CreateAppResponse       *larkbitable.CreateAppRespData
	CreateAppTablesResponse []*larkbitable.CreateAppTableRespData
}

// CreateAppAndTables 创建多维表格同时添加数据表
func CreateAppAndTables(client *lark.Client, request *CreateAppAndTablesRequest) (*CreateAppAndTablesResponse, error) {
	// 创建多维表格
	createAppReq := larkbitable.NewCreateAppReqBuilder().
		ReqApp(larkbitable.NewReqAppBuilder().
			Name(request.Name).
			FolderToken(request.FolderToken).
			Build()).
		Build()

	createAppResp, err := client.Bitable.App.Create(context.Background(), createAppReq)
	if err != nil {
		return nil, err
	}
	if !createAppResp.Success() {
		fmt.Printf("client.Bitable.App.Create failed, code: %d, msg: %s, log_id: %s\n",
			createAppResp.Code, createAppResp.Msg, createAppResp.RequestId())
		return nil, createAppResp.CodeError
	}

	// 添加数据表
	tables := make([]*larkbitable.CreateAppTableRespData, 0)
	for _, table := range request.Tables {
		req := larkbitable.NewCreateAppTableReqBuilder().
			AppToken(*createAppResp.Data.App.AppToken).
			Body(larkbitable.NewCreateAppTableReqBodyBuilder().
				Table(table).
				Build()).
			Build()

		createAppTableResp, err := client.Bitable.AppTable.Create(context.Background(), req)
		if err != nil {
			return nil, err
		}
		if !createAppTableResp.Success() {
			fmt.Printf("client.Bitable.AppTable.Create failed, code: %d, msg: %s, log_id: %s\n",
				createAppTableResp.Code, createAppTableResp.Msg, createAppTableResp.RequestId())
			return nil, createAppTableResp.CodeError
		}

		tables = append(tables, createAppTableResp.Data)
	}

	// 返回结果
	return &CreateAppAndTablesResponse{
		CodeError: &larkcore.CodeError{
			Code: 0,
			Msg:  "success",
		},
		CreateAppResponse:       createAppResp.Data,
		CreateAppTablesResponse: tables,
	}, nil
}
