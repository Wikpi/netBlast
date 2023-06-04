package pkg

import (
	"fmt"
	"os"
)

func GetWDir() {
	path, err := os.Getwd()
	if err != nil {
		LogError(err)
		fmt.Println(BadDir)
	}

	fmt.Print("Current workind dir: ", path, "\n\n")
}
