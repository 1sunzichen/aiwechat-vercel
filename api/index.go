package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/db"
)

// 手动录入
func Index(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	param := query.Get("name")
	url := query.Get("url")
	fmt.Fprintf(w, "<h1>Hello Aiwechat-Vercel!"+param+":"+url+"</h1>")
	db.ChatDbInstance.SetVideoValue(param, url)
}
