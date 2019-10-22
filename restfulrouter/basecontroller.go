package restfulrouter

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// KeySeperator seperator for keys
	KeySeperator = ":"
)

// IBaseController that satisfies restfulrouter.IBaseController can be
//auto mapping request to the certain method
type IBaseController interface {
	//Get is for the HttpGet route for the certain request
	Get(c *gin.Context)
	//Post is for the HttpPost route for the certain request
	Post(c *gin.Context)
	//Put is for the HttpPut route for the certain request
	Put(c *gin.Context)
	//Delete is for the HttpDelete route for the certain request
	Delete(c *gin.Context)
	//Patch is for the HttpPatch route for the certain request
	Patch(c *gin.Context)
	//Head is for the HttpHead route for the certain request
	Head(c *gin.Context)
	//Options is for the HttpOptions route for the certain request
	Options(c *gin.Context)

	//Mapping is the method that mapping custom request to certain method
	Mapping() map[string]GinHandler
}

//BaseController implement  IBaseController,
// cover the base functions with 404 response
type BaseController struct {
}

// Get default method
func (t *BaseController) Get(c *gin.Context) {
	returnNotResource(c)
	return
}

// Post default method
func (t *BaseController) Post(c *gin.Context) {
	returnNotResource(c)
	return
}

// Put default method
func (t *BaseController) Put(c *gin.Context) {
	returnNotResource(c)
	return
}

// Delete default method
func (t *BaseController) Delete(c *gin.Context) {
	returnNotResource(c)
	return
}

// Patch default method
func (t *BaseController) Patch(c *gin.Context) {
	returnNotResource(c)
	return
}

// Head default method
func (t *BaseController) Head(c *gin.Context) {
	returnNotResource(c)
	return
}

// Options default method
func (t *BaseController) Options(c *gin.Context) {
	returnNotResource(c)
	return
}

// Mapping default method
func (t *BaseController) Mapping() map[string]GinHandler {
	return nil
}

// returnNotResource response 404 to the client
func returnNotResource(c *gin.Context) {
	c.String(http.StatusNotFound, "")
}

// analyseMappingKey, which used in custom mapping logic, seperate the key into http method value
// and path value. The seperator is semicolon, the seperate pattern is 'httpmethod:pathname'.
// method only use the first semicolon to seperate string.
func analyseMappingKey(key string) (method string, pathName string, err error) {
	key = strings.TrimSpace(key)

	if len(key) < 5 {
		return "", "", errors.New("key has error")
	}

	if i := strings.Index(key, KeySeperator); i == -1 {
		return "", "", errors.New("key needs a comma")
	} else {
		method = key[:i]
		if err != nil {
			return "", "", errors.New("key has error")
		}

		pathName = strings.ToLower(strings.TrimSpace(key[i+1:]))
	}

	return method, pathName, nil
}

// ComposeCustomMappingKey  which used in custom mapping logic, join the http method and the path into a string.
// method is http.MethodXXX which defined in http package
func ComposeCustomMappingKey(method string, path string) string {
	return method + KeySeperator + path
}

// GinHandler define a func for gin
type GinHandler func(c *gin.Context)

// RegisterAPIRoute is the main function.use this func can auto register the method to the certain request url
func RegisterAPIRoute(ginEngine *gin.Engine, controllers []IBaseController) {
	routesControllerMapping(ginEngine, controllers)
}

// RegisterGroupAPIRoute as RegisterAPIRout, the only difference between them is group method can
// has pre base url
func RegisterGroupAPIRoute(basePath string, ginEngine *gin.Engine, controllers []IBaseController) {
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	g := ginEngine.Group(basePath)
	{
		routesControllerMapping(g, controllers)
	}
}

func routesControllerMapping(router gin.IRouter, controllers []IBaseController) {
	if controllers == nil || len(controllers) == 0 {
		return
	}
	for _, c := range controllers {
		cname, err := getControllerValidName(c)
		if err != nil {
			panic(err)
		}
		autoMapping(router, cname, c)
		err = autoCustomMapping(router, cname, c)
		if err != nil {
			panic(err)
		}
	}
}

const (
	// ControllerSuffix defin the suffix of controller struct
	ControllerSuffix = "Controller"

	// ErrorControllerName is a message of controller wrong name
	ErrorControllerName = "Controller name must be suffix with 'Controller'"
)

func getControllerValidName(controller IBaseController) (string, error) {
	typeInfo := reflect.TypeOf(controller)
	fullName := typeInfo.Elem().String()
	lastDotIndex := strings.LastIndex(fullName, ".")
	fullName = fullName[lastDotIndex+1:]
	if strings.HasSuffix(fullName, ControllerSuffix) && len(fullName) > len(ControllerSuffix) {
		validName := fullName[0 : len(fullName)-len(ControllerSuffix)]
		return strings.ToLower(strings.TrimSpace(validName)), nil
	}
	return "", errors.New(ErrorControllerName)

}

func autoMapping(router gin.IRouter, controllerName string, controller IBaseController) {
	path := "/" + controllerName
	router.GET(path, func(c *gin.Context) {
		controller.Get(c)
	})
	router.POST(path, func(c *gin.Context) {
		controller.Post(c)
	})
	router.PUT(path, func(c *gin.Context) {
		controller.Post(c)
	})
	router.DELETE(path, func(c *gin.Context) {
		controller.Post(c)
	})
	router.HEAD(path, func(c *gin.Context) {
		controller.Post(c)
	})
	router.OPTIONS(path, func(c *gin.Context) {
		controller.Post(c)
	})
	router.PATCH(path, func(c *gin.Context) {
		controller.Post(c)
	})
}

func autoCustomMapping(router gin.IRouter, controllerName string, controller IBaseController) error {
	route := controller.Mapping()

	for k, v := range route {
		method, path, err := analyseMappingKey(k)
		if err != nil {
			return err
		}
		fullPath := "/" + controllerName + "/" + path
		switch method {
		case http.MethodGet:
			func(handler GinHandler) {
				router.GET(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		case http.MethodPost:
			func(handler GinHandler) {
				router.POST(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		case http.MethodPut:
			func(handler GinHandler) {
				router.PUT(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		case http.MethodDelete:
			func(handler GinHandler) {
				router.DELETE(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		case http.MethodHead:
			func(handler GinHandler) {
				router.HEAD(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		case http.MethodOptions:
			func(handler GinHandler) {
				router.OPTIONS(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		case http.MethodPatch:
			func(handler GinHandler) {
				router.PATCH(fullPath, func(c *gin.Context) {
					handler(c)
				})
			}(v)
		}
	}
	return nil
}
