package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	err := exec.Command("bash", "-c", fmt.Sprintf("go run scripts/at/update_order_status/main.go %s | at now", os.Args[1])).Run()
	if err != nil {
		fmt.Println("Error run exec: ", err.Error())
	}
}
