/**
 * @Time: 2020/3/28 17:47
 * @Author: solacowa@gmail.com
 * @File: client
 * @Software: GoLand
 */

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func main() {
	resp, err := http.Get("http://127.0.0.1:8080/get/hello")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var res response

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Fatalf("json.NewDecoder: %v", err)
	}

	log.Println(res.Data)
}
