package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	command := fmt.Sprintf("at now+10 minutes <<< 'go run scripts/at/update_order_status/main.go %d'", os.Args[1])
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("Error run exec: ", err.Error())
	}
}
