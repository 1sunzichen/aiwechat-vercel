package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/db"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello Aiwechat-Vercel!</h1>")
	db.ChatDbInstance.SetVideoValue("哈尔滨一九四四", "test")
}
