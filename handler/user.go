package handler

import (
	dblayer "DistributedMemory/db"
	mydb "DistributedMemory/db/mysql"
	"DistributedMemory/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler 处理用户注册请求
func SignupHandler(c *gin.Context) {
	//if r.Method == http.MethodGet {
	//	data, err := ioutil.ReadFile("./static/view/signup.html")
	//	if err != nil {
	//		fmt.Println(http.StatusInternalServerError)
	//		return
	//	}
	//w.Write(data)
	//http.Redirect(w,r,"/static/view/signup.html",http.StatusFound)
	c.Redirect(http.StatusFound, "/static/view/signup.html")
	return
}

func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	if len(username) < 3 || len(password) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Invailed parameter",
			"code": -1,
		})
		return
	}

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	if dblayer.UserSignup(username, enc_passwd) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Success",
			"code": 0,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "failed",
			"code": -2,
		})
	}
}

// DoSinginHandler 处理登录post请求
func DoSinginHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	enc_passwd := util.Sha1([]byte(password + pwd_salt))

	// 1校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, enc_passwd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "login failed",
			"code": -1,
		})
		return
	}

	// 2生成访问凭证(token)
	token := GenToken(username)
	upRres := dblayer.UpdateToken(username, token)
	if !upRres {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "login failed",
			"code": -2,
		})
		return
	}

	// 3登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  username,
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	//w.Write(resp.JSONBytes())
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// SignInHandler 登录接口
func SignInHandler(c *gin.Context) {
	//if r.Method == http.MethodGet {
	// data, err := ioutil.ReadFile("./static/view/signin.html")
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(data)
	//http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
	//return
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

//
//r.ParseForm()
//username := r.Form.Get("username")
//password := r.Form.Get("password")
//enc_passwd := util.Sha1([]byte(password + pwd_salt))
//
//// 1校验用户名及密码
//pwdChecked := dblayer.UserSignin(username, enc_passwd)
//if !pwdChecked {
//	w.Write([]byte("FAILED"))
//	return
//}
//
//// 2生成访问凭证(token)
//token := GenToken(username)
//upRres := dblayer.UpdateToken(username, token)
//if !upRres {
//	w.Write([]byte("FAILED"))
//	return
//}
//
//// 3登录成功后重定向到首页
////w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
//resp := util.RespMsg{
//	Code: 0,
//	Msg:  username,
//	Data: struct {
//		Location string
//		Username string
//		Token    string
//	}{
//		Location: "http://" + r.Host + "/static/view/home.html",
//		Username: username,
//		Token:    token,
//	},
//}
//w.Write(resp.JSONBytes())
//}

// UserInfoHandler 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")

	// 2验证token是否有效
	//isValidToken := IsTokenValid(username, token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}

	// 3查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user.SignupAt,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid token是否有效
func IsTokenValid(username string, token string) bool {
	if len(token) != 40 {
		return false
	}
	// todo:判断token的时效性
	ts, _ := strconv.Atoi(token[:8])
	now := time.Now().Unix()
	if now-int64(ts) > 3600 {
		return false
	}

	// todo:从数据库表tbl_user_token查询username对应的token信息
	stmt, err := mydb.DBConn().Prepare("select user_token from tbl_user_token where user_name = ?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var t string
	err = stmt.QueryRow(username).Scan(&t)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// todo:对比两个token是否一致
	if t != token {
		return false
	}
	return true
}
