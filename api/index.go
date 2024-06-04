package api

import (
	"fmt"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, req *http.Request) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Current time: ", t)
			}
		}
	}()

	time.Sleep(10 * time.Second)
	done <- true
	fmt.Fprintf(w, "<h1>请关注公众号，电视剧资源已更新</h1>")
}
