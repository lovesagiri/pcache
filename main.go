package main

import (
	"fmt"
	"log"
	"net/http"
	"pcache/pcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	pcache.NewGroup("test", 2<<10, pcache.GetterFunc(
		func(key string) ([]byte, error) {
			if val, ok := db[key]; ok {
				return []byte(val), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9090"
	httpPool := pcache.NewHTTPPool(addr)
	log.Println("pcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, httpPool))
}
