package restfulrouter

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetControllerValidName(t *testing.T) {
	tc := TestController{}
	name, err := getControllerValidName(&tc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Test", name)
}

func TestAnalyseMappingKey1(t *testing.T) {
	temp1 := ComposeCustomMappingKey(http.MethodGet, "test/:userName")
	assert.Equal(t, "GET:test/:userName", temp1)
	method, path, err := analyseMappingKey(temp1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "test/:userName", path)
}

func TestAnalyseMappingKey2(t *testing.T) {
	temp1 := ComposeCustomMappingKey(http.MethodGet, "test?_:userName")
	assert.Equal(t, "GET:test?_:userName", temp1)
	method, path, err := analyseMappingKey(temp1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "test?_:userName", path)
}

type TestController struct {
	BaseController
}

func (tc *TestController) Mapping() map[string]GinHandler {
	m := make(map[string]GinHandler, 1)
	m[ComposeCustomMappingKey(http.MethodGet, "customTest")] = CustomMethodTest
	return m
}

func CustomMethodTest(c *gin.Context) {
	c.String(http.StatusNotFound, "")
}

func TestRegisterAPIRoute(t *testing.T) {
	RegisterAPIRoute(gin.Default(), []IBaseController{&TestController{}})
}

func TestRegisterGroupAPIRoute(t *testing.T) {
	RegisterGroupAPIRoute("/test", gin.Default(), []IBaseController{&TestController{}})
}
