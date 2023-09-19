package common

import "fmt"

// SuccessDataResponse 成功的单一数据应答
type SuccessDataResponse struct {
	Response DataResponse
}

// SuccessListResponse 成功的列表数据应答
type SuccessListResponse struct {
	Response ListResponse
}

// SuccessListStatusResponse 成功的状态数据应答
type SuccessListStatusResponse struct {
	Response ListStatusResponse
}

// FailureResponse 失败的应答
type FailureResponse struct {
	Response ErrorResponse
}

// DataResponse 单一对象
type DataResponse struct {
	Data interface{}
	//RequestId string
}

// ExDataResponse 单一对象
type ExDataResponse struct {
	Data      interface{}
	RequestId string
	Error     Error
}

// ListResponse 集合对象
type ListResponse struct {
	GroupList  []interface{}
	TotalCount int64
	//RequestId  string
}

// ListStatusResponse list status response
type ListStatusResponse struct {
	GroupList  []interface{}
	TotalCount int64
	Status     int
	RequestId  string
}

// ExplorerResponse 成功的单一数据应答
type ExplorerResponse struct {
	Response ExDataResponse
}

// ErrorResponse 异常应答
type ErrorResponse struct {
	Error Error
	//RequestId string
}

// Error 错误
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Message)
}

// StatusIntegerResponse status integer response
type StatusIntegerResponse struct {
	Status int
}

// StatusResponse status response
type StatusResponse struct {
	Status string
}

// NewStatusResponse new status response
func NewStatusResponse() *StatusResponse {
	return &StatusResponse{
		Status: "OK",
	}
}

// TokenResponse token response
type TokenResponse struct {
	Token string
}

// DownloadResponse download response
type DownloadResponse struct {
	Content string
}

// UploadResponse upload response
type UploadResponse struct {
	FileKey string
}

// NewSuccessDataResponse new success data response
func NewSuccessDataResponse(data interface{}) *SuccessDataResponse {
	dataResponse := DataResponse{
		//RequestId: NewRandomRequestId(),
		Data: data,
	}
	return &SuccessDataResponse{
		Response: dataResponse,
	}
}

// NewSuccessListResponse new success list response
func NewSuccessListResponse(datas []interface{}, count int64) *SuccessListResponse {
	listResp := ListResponse{
		GroupList:  datas,
		TotalCount: count,
		//RequestId:  NewRandomRequestId(),
	}
	return &SuccessListResponse{
		Response: listResp,
	}
}

// NewSuccessListStatusResponse new success list status response
func NewSuccessListStatusResponse(datas []interface{}, status int, count int64) *SuccessListStatusResponse {
	listResp := ListStatusResponse{
		GroupList:  datas,
		TotalCount: count,
		Status:     status,
		//RequestId:  NewRandomRequestId(),
	}
	return &SuccessListStatusResponse{
		Response: listResp,
	}
}

// NewFailureResponse new failure response
func NewFailureResponse(err *Error) *FailureResponse {
	errResponse := ErrorResponse{
		Error: *err,
	}
	return &FailureResponse{
		Response: errResponse,
	}
}

// NewError 创建错误
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
