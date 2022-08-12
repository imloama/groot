package groot

type ApiResult struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	ErrCode string      `json:"errCode"`
	Data    interface{} `json:"data"`
}

func OK(data interface{}) ApiResult {
	return ApiResult{
		Code: 200,
		Data: data,
	}
}

func Fail(msg string, errCode string) ApiResult {
	return ApiResult{
		Code:    400,
		Msg:     msg,
		ErrCode: errCode,
	}
}
