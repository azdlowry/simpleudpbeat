package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/azdlowry/simpleudpbeat/beater"
)

func main() {
	err := beat.Run("simpleudpbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
