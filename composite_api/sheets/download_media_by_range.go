/*
 下载指定范围单元格的所有素材列表，使用到两个OpenAPI：
 1. [读取单个范围](https://open.feishu.cn/document/server-docs/docs/sheets-v3/data-operation/reading-a-single-range)
 2. [下载素材](https://open.feishu.cn/document/server-docs/docs/drive-v1/media/download)
*/

package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
)

type DownloadMediaByRangeRequest struct {
	SpreadsheetToken string
	Range            string
}

type DownloadMediaByRangeResponse struct {
	*larkcore.CodeError
	ReadResponse          *SpreadsheetRespData
	DownloadMediaResponse []*larkdrive.DownloadMediaResp
}

func DownloadMediaByRange(client *lark.Client, request *DownloadMediaByRangeRequest) (*DownloadMediaByRangeResponse, error) {
	// 读取单个范围
	readResp, err := client.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod:                http.MethodGet,
		ApiPath:                   fmt.Sprintf("/open-apis/sheets/v2/spreadsheets/%s/values/%s", request.SpreadsheetToken, request.Range),
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

	// 下载文件
	files := make([]*larkdrive.DownloadMediaResp, 0)
	values := readSpreadsheetResp.Data.ValueRange.Values
	tokens := parseFileToken(values, make(map[string]bool))
	for _, token := range tokens {
		downloadMediaReq := larkdrive.NewDownloadMediaReqBuilder().
			FileToken(token).
			Build()

		downloadMediaResp, err := client.Drive.Media.Download(context.Background(), downloadMediaReq)
		if err != nil {
			return nil, err
		}
		if !downloadMediaResp.Success() {
			fmt.Printf("client.Drive.Media.Download failed, code: %d, msg: %s, log_id: %s\n",
				downloadMediaResp.Code, downloadMediaResp.Msg, downloadMediaResp.RequestId())
			return nil, downloadMediaResp.CodeError
		}

		files = append(files, downloadMediaResp)
	}

	// 返回结果
	return &DownloadMediaByRangeResponse{
		CodeError: &larkcore.CodeError{
			Code: 0,
			Msg:  "success",
		},
		ReadResponse:          readSpreadsheetResp.Data,
		DownloadMediaResponse: files,
	}, nil
}

func parseFileToken(values []interface{}, tokens map[string]bool) []string {
	if len(values) == 0 {
		res := make([]string, 0, len(tokens))
		for k := range tokens {
			res = append(res, k)
		}
		return res
	}
	for _, i := range values {
		kind := reflect.TypeOf(i).Kind()
		if kind == reflect.Slice {
			parseFileToken(i.([]interface{}), tokens)
		} else if kind == reflect.Map {
			m := i.(map[string]interface{})
			if val, ok := m["fileToken"]; ok {
				tokens[val.(string)] = true
			}
		}
	}

	res := make([]string, 0, len(tokens))
	for k := range tokens {
		res = append(res, k)
	}
	return res
}
