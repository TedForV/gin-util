package base

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/tedforv/gin-util/log"
	"github.com/tedforv/gin-util/restfulrouter"
)

var (
	corsAllowAllOrigins = []string{"*"}
)

// CreateDefaultGin Create *gin.Engine with default log,using logrus
func CreateDefaultGin(
	isProduct bool,
	isCors bool,
	logFolderPath string,
	corsAllowOrigins []string,
	corsAllowHeaders []string,
	groupedControllers map[string][]restfulrouter.IBaseController) (*gin.Engine, error) {

	if len(logFolderPath) == 0 {
		return nil, errors.New("log folder path is nil")
	}
	logger, err := log.NewLogrusLogger(logFolderPath)
	if err != nil {
		return nil, err
	}
	return CreateGin(isProduct, isCors, logger, corsAllowAllOrigins, corsAllowHeaders, groupedControllers)
}

// CreateGin create *gin.Engine with custom logger
func CreateGin(
	isProduct bool,
	isCors bool,
	logger restfulrouter.ILogger,
	corsAllowOrigins []string,
	corsAllowHeaders []string,
	groupedControllers map[string][]restfulrouter.IBaseController) (*gin.Engine, error) {

	setMode(isProduct)

	r := gin.Default()

	if logger != nil {
		restfulrouter.SetLogger(logger)
	}

	if isCors {
		setCors(r, corsAllowOrigins, corsAllowHeaders)
	}

	r.RedirectFixedPath = true

	for k, v := range groupedControllers {
		restfulrouter.RegisterGroupAPIRoute(k, r, v)
	}
	return r, nil
}

func setMode(isProduct bool) {
	if isProduct {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func hookLogger(logFolderPath string) (*log.LogrusLogger, error) {
	return log.NewLogrusLogger(logFolderPath)
}

func setCors(r *gin.Engine, corsAllowOrigins []string, allowHeaders []string) {
	if corsAllowOrigins == nil || len(corsAllowOrigins) == 0 {
		corsAllowOrigins = corsAllowAllOrigins
	}
	c := cors.DefaultConfig()
	c.AllowOrigins = corsAllowAllOrigins
	c.AllowCredentials = true
	if allowHeaders != nil && len(allowHeaders) > 0 {
		c.AddAllowHeaders(allowHeaders...)
	}
	r.Use(cors.New(c))
}
