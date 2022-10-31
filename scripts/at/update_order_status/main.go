package main

import (
	"fmt"
	"kamari.com/db"
	"kamari.com/repository"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Println("Shutting down server: ", err.Error())
		os.Exit(1)
	}
	orderRepository := repository.NewOrderRepository(repository.ORConfig{DB: db.Get()})
	if len(os.Args) == 2 {
		orderId, _ := strconv.Atoi(os.Args[1])
		updateFromOnDeliveryToDelivered(orderRepository, orderId)
		//testOnly(orderId)
	}
}

func testOnly(orderId int) {
	command := fmt.Sprintf("at now+10 minutes <<< 'echo order id %d is updated > ~/text.txt'", orderId)
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("Error run exec: ", err.Error())
	}
}

func updateFromOnDeliveryToDelivered(orderRepository repository.OrderRepository, orderId int) {
	fmt.Println("/>>---------------")
	fmt.Println("START EXECUTE AT: UPDATE ORDER STATUS FROM ON DELIVERY TO DELIVERED", time.Now())
	err := orderRepository.CronUpdateStatusFromOnDeliveryToDelivered(orderId)
	if err != nil {
		fmt.Println("ERROR EXECUTE AT: UPDATE ORDER STATUS FROM ON DELIVERY TO DELIVERED")
		fmt.Println(err.Error())
		fmt.Println("---------------<</")
		return
	}
	fmt.Println("SUCCESS EXECUTE AT: UPDATE ORDER STATUS FROM ON DELIVERY TO DELIVERED, rows affected: 1")
	fmt.Println(fmt.Sprintf("1. order id: %d", orderId))
	fmt.Println("---------------<</")
}
