package handler

import (
	dblayer "DistributedMemory/db"
	"DistributedMemory/util"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			fmt.Println(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invailed parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	if dblayer.UserSignup(username, enc_passwd) {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}
