package main

import (
	"fmt"
	"log"

	"github.com/josh/datajar-server/internal/datajar/sqlite"
)

func main() {
	store, err := sqlite.FetchStore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(store)
}
