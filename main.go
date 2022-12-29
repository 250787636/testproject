package main

import (
	"example.com/m/v2/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var CONF = config.LoadConfig()

func main() {
	//r := gin.Default()
	//
	//r.POST("api/ticket/v1/state", IFLOGIN01)
	//r.POST("api/ticket/v1/analysis", IFLOGIN02)
	//
	//r.Run(":8888")
	a := []string{"1","2","3","4","5"}
	_ = append(a, "6")
	fmt.Println(a)

}

type Login01 struct {
	RetCode string `json:"retCode"`
	Msg     string `json:"msg"`
}

type Login02 struct {
	RetCode  string   `json:"retCode"`
	Msg      string   `json:"msg"`
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

type UserInfo struct {
	AccountID string `json:"accountID"`
	Name      string `json:"name"`
	EmpNo     string `json:"empNo"`
	IdCardNum string `json:"idCardNum"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Tenant    string `json:"tenant"`
}

type Login02Param struct {
	AppCode string `json:"appCode"`
	Tenant  string `json:"tenant"`
	Data    string `json:"iamcaspticket" binding:"required"`
}

func IFLOGIN01(c *gin.Context) {
	//如果retCode的值不是1000说明服务不可用；或者根据msg的值不是serverisok说明服务不可用
	c.JSON(http.StatusOK, Login01{
		RetCode: "1000",
		Msg:     "serverisok",
	})
}

func IFLOGIN02(c *gin.Context) {
	var param Login02Param
	err := c.ShouldBind(&param)
	if err != nil {
		fmt.Println(err)
	}
	param.AppCode = CONF.CMIC4AService.AppCode
	param.Tenant = CONF.CMIC4AService.Tenant

	// 有效的票据串retCode是1000，否则是非法票据串，token是4A的会话ID，userInfo是帐号数据
	c.JSON(http.StatusOK, Login02{
		RetCode: "1000",
		Msg:     "serverisok",
		Token:   "",
		UserInfo: UserInfo{
			AccountID: "18860230359",
			Name:      "",
			EmpNo:     "",
			IdCardNum: "",
			Phone:     "",
			Mobile:    "",
			Email:     "",
			Tenant:    "",
		},
	})
}
