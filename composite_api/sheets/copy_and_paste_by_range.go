/*
 复制粘贴某个范围的单元格数据，使用到两个OpenAPI：
 1. [读取单个范围](https://open.feishu.cn/document/server-docs/docs/sheets-v3/data-operation/reading-a-single-range)
 2. [向单个范围写入数据](https://open.feishu.cn/document/server-docs/docs/sheets-v3/data-operation/write-data-to-a-single-range)
*/

package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type CopyAndPasteByRangeRequest struct {
	SpreadsheetToken string
	SrcRange         string
	DstRange         string
}

type CopyAndPasteRangeResponse struct {
	*larkcore.CodeError
	ReadResponse  *SpreadsheetRespData
	WriteResponse *SpreadsheetRespData
}

// CopyAndPasteRange 复制粘贴某个范围的单元格数据
func CopyAndPasteRange(client *lark.Client, request *CopyAndPasteByRangeRequest) (*CopyAndPasteRangeResponse, error) {
	// 读取单个范围
	readResp, err := client.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod:                http.MethodGet,
		ApiPath:                   fmt.Sprintf("/open-apis/sheets/v2/spreadsheets/%s/values/%s", request.SpreadsheetToken, request.SrcRange),
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
	})
	if err != nil {
		return nil, err
	}

	readSpreadsheetResp := &SpreadsheetResp{}
	err = json.Unmarshal(readResp.RawBody, readSpreadsheetResp)
	if err != nil {
		return nil, err
	}
	readSpreadsheetResp.ApiResp = readResp
	if readSpreadsheetResp.Code != 0 {
		fmt.Printf("read spreadsheet failed, code: %d, msg: %s, log_id: %s\n",
			readSpreadsheetResp.Code, readSpreadsheetResp.Msg, readSpreadsheetResp.RequestId())
		return nil, readSpreadsheetResp.CodeError
	}

	// 向单个范围写入数据
	valueRange := map[string]interface{}{}
	valueRange["range"] = request.DstRange
	valueRange["values"] = readSpreadsheetResp.Data.ValueRange.Values
	body := map[string]interface{}{}
	body["valueRange"] = valueRange

	writeResp, err := client.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod:                http.MethodPut,
		ApiPath:                   fmt.Sprintf("/open-apis/sheets/v2/spreadsheets/%s/values", request.SpreadsheetToken),
		Body:                      body,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
	})
	if err != nil {
		return nil, err
	}

	writeSpreadsheetResp := &SpreadsheetResp{}
	err = json.Unmarshal(writeResp.RawBody, writeSpreadsheetResp)
	if err != nil {
		return nil, err
	}
	writeSpreadsheetResp.ApiResp = writeResp
	if writeSpreadsheetResp.Code != 0 {
		fmt.Printf("write spreadsheet failed, code: %d, msg: %s, log_id: %s\n",
			writeSpreadsheetResp.Code, writeSpreadsheetResp.Msg, writeSpreadsheetResp.RequestId())
		return nil, writeSpreadsheetResp.CodeError
	}

	// 返回结果
	return &CopyAndPasteRangeResponse{
		CodeError: &larkcore.CodeError{
			Code: 0,
			Msg:  "success",
		},
		ReadResponse:  readSpreadsheetResp.Data,
		WriteResponse: writeSpreadsheetResp.Data,
	}, nil
}
