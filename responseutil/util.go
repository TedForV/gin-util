package responseutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseObj define  a response object
type ResponseObj struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg"`
	Data interface{}  `json:"data"`
}

// ResponseCode define the new type for response code
type ResponseCode int8

const (
	_ ResponseCode = iota
	// RSuccess is success status
	RSuccess
	// RError is error status
	RError
	// ROther is other status
	ROther
)

// GinResponseObj convert ResponseObj into gin.H
func GinResponseObj(o *ResponseObj) gin.H {
	return gin.H{
		"code": o.Code,
		"msg":  o.Msg,
		"data": o.Data,
	}
}

// Error response error message
func Error(c *gin.Context, additionalInfo string) {
	var msg string
	if len(additionalInfo) == 0 {
		msg = "system error occured"
	} else {
		msg = additionalInfo
	}

	result := ResponseObj{
		RError,
		msg,
		nil,
	}
	c.JSON(http.StatusOK, GinResponseObj(&result))
}
