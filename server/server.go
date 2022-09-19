package server

import (
	"log"
)

func Init() {
	r := NewRouter(&RouterConfig{})

	log.Fatal(r.Run())
}
