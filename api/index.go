package api

import (
	"fmt"
	"net/http"

	Videourl "github.com/pwh-pwh/aiwechat-vercel/chat/videourl"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello Aiwechat-Vercel!</h1>")
	go func() {

		Videourl.VideoConvert("哈尔滨一九四四")
	}()

}
