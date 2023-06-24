package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
	time.Sleep(1 * time.Second)
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
}
