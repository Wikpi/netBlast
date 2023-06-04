package pkg

import (
	"fmt"
	"os"
	"time"
)

// Clears logs file
func ClearLogs() {
	err := os.Truncate(Logs, 0)
	if err != nil {
		LogError(err)
		fmt.Println("Failed to clear logs file.")
	}
}

// Handles incoming error and gives out a simplified version
func LogError(incomingErr error) {
	file, err := os.OpenFile(Logs, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	if _, err := file.WriteString(time.Now().Format("2006-01-02 15:04") + " " + incomingErr.Error() + "\n\n"); err != nil {
		fmt.Println(err)
	}
}
