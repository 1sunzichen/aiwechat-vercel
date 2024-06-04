package api

import (
	"fmt"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, req *http.Request) {
	time.AfterFunc(1*time.Second, func() {
		fmt.Fprintf(w, "Hello, world2222!")
	})

	fmt.Fprintf(w, "Hello, world!")

}
