package main

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	options, err := redis.ParseURL("redis://default:5e27d347869141eeb77127ccbe40b5ad@in-maggot-33849.upstash.io:33849")
	if err != nil {
		fmt.Println(err)
		return
	}
	options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	client := redis.NewClient(options)
	client.Set(context.Background(), "test", "sjq", 0)
	fmt.Println("success")
}
