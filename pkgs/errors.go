package pkgs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"restaurant/pkgs/sferror"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Unauthorized = sferror.NewError(401, errors.New("token已过期"))
var ParamError = sferror.NewError(400, errors.New("参数错误"))
var NoPermission = sferror.NewError(422, errors.New("无权限"))
var SignError = sferror.NewError(422, errors.New("签名错误"))
var RetryAgainError = sferror.NewError(1, errors.New("服务超时, 请重试"))
var TooManyFiles = sferror.NewError(1, errors.New("最多支持上传20个文件"))
var ServerError = sferror.NewError(1, errors.New("服务错误"))
var TokenParseError = sferror.NewError(10001, errors.New("token解析错误"))
var LoginInfoGetError = sferror.NewError(10002, errors.New("获取登陆信息失败"))
var LoginPasswordError = sferror.NewError(10003, errors.New("登陆密码错误"))
var PageError = sferror.NewError(10004, errors.New("分页错误"))
var UserNotFound = sferror.NewError(10005, errors.New("用户不存在"))
var AppNotFound = sferror.NewError(10006, errors.New("商家不存在"))
var AccountNotFound = sferror.NewError(10007, errors.New("账户不存在"))
var SwitchNoPermission = sferror.NewError(10008, errors.New("未开通插件权限"))
var NotWhiteIp = sferror.NewError(10009, errors.New("未配置白名单IP"))
var VarifyCodeError = sferror.NewError(10010, errors.New("验证码错误"))

func wrapErr(c *gin.Context, err error) {
	// 错误处理
	c.Errors = append(c.Errors, &gin.Error{
		Err: fmt.Errorf("%+v, %s", err, zap.Stack("stackstrace").String),
	})
	switch nerr := sferror.Cause(err).(type) {
	case *sferror.Error:
		c.JSON(200, gin.H{"msg": nerr.Err.Error(), "code": nerr.Code, "body": ""})
	default:
		c.JSON(200, gin.H{"msg": err.Error(), "code": 1, "body": ""})
	}
}

func wrapRecoverErr(c *gin.Context, err interface{}) {
	ginErr := &gin.Error{
		Err: fmt.Errorf("%+v", err),
	}
	switch e := err.(type) {
	case *sferror.Error:
		c.JSON(200, gin.H{"msg": e.Err.Error(), "code": e.Code, "body": ""})
	default:
		// 生成堆栈信息
		stack := zap.Stack("stackstrace")
		ginErr.Err = fmt.Errorf("%+v, %s", err, stack.String)
		c.JSON(200, gin.H{"msg": fmt.Sprintf("%+v", err), "code": 1, "body": ""})
	}
	c.Errors = append(c.Errors, ginErr)
}

// handlerFunc 函数
type handlerFunc func(*gin.Context) (interface{}, error)

// WrapperHandler 请求包装器
func WrapperHandler(handler handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				wrapRecoverErr(c, err)
			}
		}()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Max-Age", "86400")

		body, err := handler(c)
		if err != nil {
			wrapErr(c, err)
		} else {
			c.JSON(200, gin.H{"code": 0, "body": body, "msg": "ok"})
		}
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "*")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}
