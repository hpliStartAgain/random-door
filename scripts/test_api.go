package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:8080/api/cities")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body length: %d\n", len(body))
	if len(body) > 200 {
		fmt.Printf("Preview: %s...\n", string(body[:200]))
	} else {
		fmt.Printf("Body: %s\n", string(body))
	}
}
