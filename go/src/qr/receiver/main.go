package main

import (
	"fmt"

	"github.com/tuvirz/qr/go/src/qr/receiver/receiver"
)

func main() {
	filaName, err := receiver.RunReceiver()
	if err != nil {
		fmt.Println("--- ERROR:", err)
	} else {
		fmt.Println("receiver: file was transfered successfully to", filaName)
	}
}
