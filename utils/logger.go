package utils

import (
	"fmt"
	"os"
	"time"
)

func LogAction(userID int, action, status, details string) {
	logEntry := fmt.Sprintf("%d,%s,%s,%s,%s\n", userID, action, time.Now().Format(time.RFC3339), status, details)
	file, err := os.OpenFile("utils/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(logEntry)
}
