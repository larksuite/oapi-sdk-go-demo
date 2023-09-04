package base

import (
	"fmt"
	"testing"

	oapi_sdk_go_demo "oapi-sdk-go-demo"
	"oapi-sdk-go-demo/composite_api/base"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

func TestCreateAppAndTables(t *testing.T) {
	tableName := "这是数据表"
	fieldName := "这是字段名"
	type_ := 1
	req := &base.CreateAppAndTablesRequest{
		Name:        "这是多维表格",
		FolderToken: "Y9LhfoWNZlKxWcdsf2fcPP0SnXc",
		Tables: []*larkbitable.ReqTable{
			{
				Name: &tableName,
				Fields: []*larkbitable.AppTableCreateHeader{
					{
						FieldName: &fieldName,
						Type:      &type_,
					},
				},
			},
		},
	}

	resp, err := base.CreateAppAndTables(oapi_sdk_go_demo.Client, req)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
}
