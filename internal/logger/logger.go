package logger

import (
	"container/ring"
	"fmt"
	"log"
	"os"
)

const logBufferSize = 50

const (
	flags      = log.Ldate | log.Ltime | log.Lshortfile
	boldRed    = "\033[1;31m"
	boldGreen  = "\033[1;32m"
	boldYellow = "\033[1;33m"
	boldPurple = "\033[1;35m"
	boldCyan   = "\033[1;36m"
	intenseRed = "\033[0;91m"
	reset      = "\033[0m"
)

var (
	Debug   = log.New(os.Stdout, fmt.Sprintf("%s[DEBUG]%s ", boldCyan, reset), flags)
	Info    = log.New(os.Stdout, fmt.Sprintf("%s[INFO]%s ", boldGreen, reset), flags)
	Warning = log.New(os.Stdout, fmt.Sprintf("%s[WARNING]%s ", boldYellow, reset), flags)
	Error   = log.New(os.Stderr, fmt.Sprintf("%s[ERROR]%s ", boldRed, reset), flags)
	Fatal   = log.New(os.Stderr, fmt.Sprintf("%s[FATAL]%s", intenseRed, reset), flags)
)

var logBuffer *ring.Ring

type logEntry struct {
	Level   string
	Message string
}

func LogError(log string, message ...interface{}) {
	fullMessage := log + fmt.Sprint(message...)
	Error.Println(fullMessage)
	storeLogEntry("ERROR", fullMessage)
}

func LogWarning(log string, message ...interface{}) {
	fullMessage := log + fmt.Sprint(message...)
	Warning.Println(fullMessage)
	storeLogEntry("WARN", fullMessage)
}

func LogInfo(log string, message ...interface{}) {
	fullMessage := log + fmt.Sprint(message...)
	Info.Println(fullMessage)
	storeLogEntry("INFO", fullMessage)
}

func LogDebug(log string, message ...interface{}) {
	fullMessage := log + fmt.Sprint(message...)
	Debug.Println(fullMessage)
	storeLogEntry("DEBUG", fullMessage)
}

func storeLogEntry(level, message string) {
	logBuffer.Value = logEntry{Level: level, Message: message}
	logBuffer = logBuffer.Next()
}

func RetrieveLogs() []logEntry {
	entries := make([]logEntry, 0, logBufferSize)
	logBuffer.Do(func(p interface{}) {
		if entry, ok := p.(logEntry); ok {
			entries = append(entries, entry)
		}
	})
	return entries
}

func init() {

	logBuffer = ring.New(logBufferSize)
	stdout := os.Stdout
	stderr := os.Stderr

	Debug.SetOutput(stdout)
	Info.SetOutput(stdout)
	Error.SetOutput(stderr)
	Fatal.SetOutput(stderr)

	Info.Println("Initialized Logger")
}
