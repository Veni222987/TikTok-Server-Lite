package test

import (
	"DoushengABCD/controller"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	router := gin.New()
	router.POST("/douyin/user/login", controller.Login)

	// 创建一个模拟的HTTP请求
	req, _ := http.NewRequest("POST", "/douyin/user/login?username=Veni&password=asdfghjkl", nil)
	resp := httptest.NewRecorder()

	// 将请求发送到路由引擎处理
	router.ServeHTTP(resp, req)

	// 验证响应
	if resp.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", resp.Code)
	}
}
