package util

import (
	"fmt"
	"os"
)

func RuntimeError(err error) {
	fmt.Printf("runtime error: %v\n", err)
	os.Exit(1)
}
