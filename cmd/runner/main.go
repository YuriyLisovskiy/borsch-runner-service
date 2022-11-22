package main

import (
	"log"

	"github.com/YuriyLisovskiy/borsch-runner-service/internal/cli"
)

func main() {
	err := cli.ExecuteApp()
	if err != nil {
		log.Fatalln(err)
	}
}
