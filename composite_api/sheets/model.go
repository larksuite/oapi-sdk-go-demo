package sheets

import larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

type ValueRange struct {
	MajorDimension string
	Range          string
	Revision       int
	Values         []interface{}
}

type SpreadsheetRespData struct {
	ValueRange *ValueRange
}

type SpreadsheetResp struct {
	*larkcore.ApiResp
	larkcore.CodeError
	Data *SpreadsheetRespData
}
