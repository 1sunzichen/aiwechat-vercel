package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/db"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello Aiwechat-Vercel!</h1>")
	query := req.URL.Query()
	param := query.Get("name")
	url := query.Get("url")
	db.ChatDbInstance.SetVideoValue(param, url)
}
