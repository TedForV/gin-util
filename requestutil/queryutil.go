package requestutil

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

const daySeconds int64 = 86400

// GetDateRangeFromURLQuery fetch startDate and endDate from url querystring
func GetDateRangeFromURLQuery(c *gin.Context, maxIntervalDay int64) (int64, int64, error) {
	startDateStr := c.DefaultQuery("sDate", "")
	endDateStr := c.DefaultQuery("eDate", "")
	startDate := now.BeginningOfDay().Unix()
	endDate := now.EndOfDay().Unix()
	var err error
	if len(startDateStr) != 0 {
		startDate, err = strconv.ParseInt(startDateStr, 10, 64)
		if err != nil {
			startDate = now.BeginningOfDay().Unix()
		}
	}
	if len(endDateStr) != 0 {
		endDate, err = strconv.ParseInt(endDateStr, 10, 64)
		if err != nil {
			endDate = now.EndOfDay().Unix()
		}
	}
	if endDate-startDate > daySeconds*maxIntervalDay {
		return 0, 0, errors.New("interval time is overflowed")
	}
	return startDate, endDate, nil
}

// GetPageInfo fetch pageNo & pageRow in url querystring
func GetPageInfo(c *gin.Context, maxRows int) (int, int, error) {
	pNo := GetQueryInt(c, "pageNo", 1)
	pRow := GetQueryInt(c, "pageRow", 50)

	if pRow > maxRows {
		return 1, 50, errors.New("request rows is overflowed")
	}
	return pNo, pRow, nil
}

// TRUEString for true string
const TRUEString = "TRUE"

// FALSEString for false string
const FALSEString = "FALSE"

// ConvertBoolenFromString get bool value from string value
func ConvertBoolenFromString(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	formatValue := strings.ToUpper(value)
	switch formatValue {
	case TRUEString:
		return true
	case FALSEString:
		return false
	default:
		return defaultValue
	}
}

// GetQueryInt get int value from url querystring
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.DefaultQuery(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetQueryBool get bool value from url querystring
func GetQueryBool(c *gin.Context, key string, defaultValue bool) bool {
	valueStr := c.DefaultQuery(key, "")
	return ConvertBoolenFromString(valueStr, defaultValue)

}

// GetQueryDate get datetime value from url querystring
func GetQueryDate(c *gin.Context, key string, defaultValue time.Time) time.Time {
	valueStr := c.DefaultQuery(key, "")
	unix, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return time.Unix(unix, 0)
}

// GetQueryInt64 get int64 value from url querystring
func GetQueryInt64(c *gin.Context, key string, defaultValue int64) int64 {
	valueStr := c.DefaultQuery(key, "")
	unix, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return unix
}
