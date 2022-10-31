package main

import (
	"fmt"
	"kamari.com/db"
	"kamari.com/repository"
	"log"
	"os"
	"time"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Println("Shutting down server: ", err.Error())
		os.Exit(1)
	}
	productPromoRepository := repository.NewProductPromoRepository(repository.PPRConfig{
		DB: db.Get(),
	})
	deleteExpiredPromotion(productPromoRepository)
}

func deleteExpiredPromotion(productPromoRepository repository.ProductPromoRepository) {
	fmt.Println("/>>---------------")
	fmt.Println("START EXECUTE CRON: DELETE EXPIRED PROMOTION", time.Now())
	promos, count, err := productPromoRepository.CronDeleteExpiredPromotion()
	if err != nil {
		fmt.Println("ERROR EXECUTE CRON: DELETE EXPIRED PROMOTION")
		fmt.Println(err.Error())
		fmt.Println("---------------<</")
		return
	}
	fmt.Println("SUCCESS EXECUTE CRON: DELETE EXPIRED PROMOTION, rows affected: ", count)
	for i, promo := range promos {
		fmt.Println(fmt.Sprintf("%d. promo id: %d", i+1, promo.Id))
	}
	fmt.Println("---------------<</")
}
