package api

import (
	"fmt"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, req *http.Request) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for t := range ticker.C {
			s := t.Format("2006-01-02 15:04:05")
			fmt.Fprintf(w, "<h1>Current time2: %s</h1>\n", s)
		}
	}()
	for t := range ticker.C {
		s := t.Format("2006-01-02 15:04:05")
		fmt.Fprintf(w, "<h1>Current time1: %s</h1>\n", s)
	}
}
