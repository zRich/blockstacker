package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zRich/cm-api-server/src/logger"
)

var (
	Log = logger.GetLogger("API")
)

// ConvergeDataResponse 汇聚单一对象应答结果
func ConvergeDataResponse(ctx *gin.Context, data interface{}, err *Error) {
	// 首先判断err是否为空
	if err == nil {
		successResponse := NewSuccessDataResponse(data)
		ctx.JSON(http.StatusOK, successResponse)
	} else {
		ConvergeHandleFailureResponse(ctx, err)
	}
}

// ConvergeListResponse 汇聚集合对象应答结果
func ConvergeListResponse(ctx *gin.Context, datas []interface{}, count int64, err *Error) {
	// 首先判断err是否为空
	if err == nil {
		successResponse := NewSuccessListResponse(datas, count)
		ctx.JSON(http.StatusOK, successResponse)
	} else {
		ConvergeHandleFailureResponse(ctx, err)
	}
}

// ConvergeFailureResponse 汇聚失败应答
func ConvergeFailureResponse(ctx *gin.Context, errCode ErrCode) {
	err := Error{
		Code:    ErrCodeName[errCode],
		Message: ErrCodeMsg[errCode][1],
	}
	Log.Errorf("Http request[%s]'s error = [%s]", ctx.Request.URL.String(), err.Error())
	failureResponse := NewFailureResponse(&err)
	ctx.JSON(http.StatusOK, failureResponse)
}

// CreateError create error
func CreateError(errCode ErrCode) *Error {
	return &Error{
		Code:    ErrCodeName[errCode],
		Message: ErrCodeMsg[errCode][1],
	}
}

// ConvergeHandleFailureResponse 汇聚处理异常的应答
func ConvergeHandleFailureResponse(ctx *gin.Context, err error) {
	newError := &Error{
		Code:    ErrCodeName[ErrorHandleFailure],
		Message: err.Error(),
	}
	Log.Errorf("Http request[%s]'s error = [%s]", ctx.Request.URL.String(), err.Error())
	failureResponse := NewFailureResponse(newError)
	ctx.JSON(http.StatusOK, failureResponse)
}

// ConvergeHandleErrorResponse converge handle error
func ConvergeHandleErrorResponse(ctx *gin.Context, err *Error) {
	Log.Errorf("Http request[%s]'s error = [%s]", ctx.Request.URL.String(), err.Error())
	failureResponse := NewFailureResponse(err)
	ctx.JSON(http.StatusOK, failureResponse)
}

// BindBody 绑定参数
func BindBody(ctx *gin.Context, body RequestBody) error {
	if err := ctx.ShouldBindJSON(body); err != nil {
		Log.Error("resolve param error:", err)
		return err
	}
	return nil
}
