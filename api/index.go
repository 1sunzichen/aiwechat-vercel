package api

import (
	"fmt"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, req *http.Request) {
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("One second later...")
	})

	fmt.Fprintf(w, "Hello, world!")

}
