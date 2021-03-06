package web

// 此Controller用于 无需鉴权 即可访问的cgi
import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// index 子站的 路由映射
var HttpuriIndex = map[string]func(s *webServer, c *gin.Context){
	"login":     IndexLogin,
	"do_login":  IndexDoLogin,
	"do_logout": IndexDoLogout,
}

// 登录页面
func IndexLogin(s *webServer, c *gin.Context) {
	checkLogin(s, c)
	c.HTML(http.StatusOK, "index/login.html", gin.H{})
}

//处理登录
func IndexDoLogin(s *webServer, c *gin.Context) {
	conf := GetConf()
	user := conf.WebUi.User
	password := conf.WebUi.Password
	str1 := user + password
	str2 := c.PostForm("user") + c.PostForm("password")
	h := md5.New()
	h.Write([]byte(str1))
	md51 := hex.EncodeToString(h.Sum(nil))
	h = md5.New()
	h.Write([]byte(str2))
	md52 := hex.EncodeToString(h.Sum(nil))
	if md51 == md52 {
		cookie := &http.Cookie{
			Name:     "userinfo",
			Value:    md51,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * 30 * time.Hour), //默认30天内cookie 有效期
		}
		http.SetCookie(c.Writer, cookie)
		c.JSON(200, gin.H{"code": 0, "msg": "登录成功"})
		return
	}
	c.JSON(200, gin.H{"code": -1, "msg": "登录失败"})
}

func IndexDoLogout(s *webServer, c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "userinfo",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(c.Writer, cookie)
	c.HTML(http.StatusOK, "index/jump.html", gin.H{
		"url":     "/index/login",
		"timeout": "3",
		"code":    1, //1为success,0为error
		"msg":     "已退出登录，即将转跳到登录页",
	})
	c.Abort()
	return
}

func checkLogin(s *webServer, c *gin.Context) {
	conf := GetConf()
	user := conf.WebUi.User
	password := conf.WebUi.Password
	if user == "" || password == "" {
		//无需登录
		c.Redirect(302, "/admin/index")
		return
	}
	str1 := user + password
	h := md5.New()
	h.Write([]byte(str1))
	md51 := hex.EncodeToString(h.Sum(nil))
	if cookie, err := c.Request.Cookie("userinfo"); err == nil {
		value := cookie.Value
		if value == md51 {
			c.Redirect(302, "/admin/index")
			return
		}
	}
	return
}
