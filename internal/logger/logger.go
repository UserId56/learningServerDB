package logger

import "fmt"

const (
	Info    = "Info"
	Warning = "Warning"
	Error   = "Error"
)

func LogError(err error, message, typeNotification string) {
	fmt.Printf("%s. %s: %v\n", typeNotification, message, err)
}
