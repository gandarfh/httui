package utils

import (
	"fmt"
	"log"
)

func MsgError(err error, msgs ...string) {
	for _, item := range msgs {
		fmt.Printf(item)
	}

	log.Fatal(err)
}
