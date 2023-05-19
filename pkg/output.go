package pkg

import (
	"fmt"
	"os"
)

func GetWDir() {
	path, err := os.Getwd()
	HandleError("WorkingDir: couldnt get directory.", err, 1)

	fmt.Print("Current workind dir: ", path, "\n\n")
}
