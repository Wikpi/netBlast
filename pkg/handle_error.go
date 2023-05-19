package pkg

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Handles incoming error and gives out a simplified version
func HandleError(errMsg string, incomingErr error, action ...int) {
	if incomingErr == nil {
		return
	}

	file, err := os.OpenFile(Logs, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	// Writes error to logs file
	if _, err := file.WriteString(time.Now().Format("2006-01-02 15:04") + " " + incomingErr.Error() + "\n\n"); err != nil {
		fmt.Println(err)
	}

	// Exits program and gives message where error occured
	switch action[0] {
	case 0:
		log.Fatal("FATAL - ", errMsg)
	case 1:
		fmt.Println(errMsg)
	default:
		return
	}
}
