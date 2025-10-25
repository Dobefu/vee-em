// The main entrypoint of the application.
package main

import (
	"log"

	"github.com/Dobefu/vee-em/vm"
)

func main() {
	err := vm.New([]byte{}).Run()

	if err != nil {
		log.Fatalf("Failed to run VM: %s\n", err.Error())
	}
}
