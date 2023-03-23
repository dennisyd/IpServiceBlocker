package base

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	zlog "github.com/rs/zerolog/log"
)

// Returns the desired environment variable as string
// If the environment variable does not exist, returns the value in `def`
func GetEnvAsStringOrDefault(key string, def string) string {
	s, exists := os.LookupEnv(key)
	if !exists {
		return def
	}
	return s
}

// Returns the desired environment variable as int
// If the environment variable does not exist, returns the value in `def`
// On error, returns 0
func GetEnvAsIntOrDefault(key string, def int) int {
	s, exists := os.LookupEnv(key)
	if !exists {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Println(fmt.Sprintf("Could not convert environment variable \"%v\" with value \"%v\" to int", key, s))
		return 0
	}
	return v
}

// Returns the desired environment variable as int
// If the environment variable does not exist, returns the value in `def`
// On error, returns 0
func GetEnvAsBoolOrDefault(key string, def bool) bool {
	s, exists := os.LookupEnv(key)
	if !exists {
		return def
	}
	if strings.ToLower(s) == "true" {
		return true
	}

	return false
}

func GetFunctionName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.Function
}

// Helper function to log function name
func HandleNewMessage(ctx context.Context) (func(context.Context, string), string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	PrintLogLine(ctx, fmt.Sprintf("--> %v", frame.Function))
	return LogWhenFinished, frame.Function
}

func LogWhenFinished(ctx context.Context, functionName string) {
	PrintLogLine(ctx, fmt.Sprintf("<-- %v", functionName))
}

func LogResponse(ctx context.Context, res any) {
	PrintLogLine(ctx, fmt.Sprintf("Response %+v", res))
}

type structuredLog struct {
	MsgId   string `json:"msgId"`
	Message string
}

// Helper function for ensuring log messages follow the same format
func PrintLogLine(ctx context.Context, logMessage string) {
	if ctx != nil {
		// log.Println(fmt.Sprintf("msgId: %v : %v", messageId, logMessage))
		val := ctx.Value(CONTEXT_MESSAGE_ID)
		zlog.Print(fmt.Sprintf("msgId: %v : %v", val, logMessage))
		return
	}

	// log.Println(logMessage)
	zlog.Print(logMessage)
}
