package main

import (
	"log"

	"YuriyLisovskiy/borsch-runner-service/internal/cli"
)

func main() {
	err := cli.ExecuteApp()
	if err != nil {
		log.Fatalln(err)
	}
}
