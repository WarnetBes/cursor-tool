package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var quiet bool

func SetQuiet(q bool) { quiet = q }

func Info(f string, a ...interface{}) {
	if !quiet { fmt.Printf("[INFO]  "+f+"\n", a...) }
}
func Success(f string, a ...interface{}) {
	if !quiet { fmt.Printf("[OK]    "+f+"\n", a...) }
}
func Warn(f string, a ...interface{}) {
	if !quiet { fmt.Fprintf(os.Stderr, "[WARN]  "+f+"\n", a...) }
}
func Error(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] "+f+"\n", a...)
}

// MaskUUID shows only first 8 and last 4 chars.
func MaskUUID(u string) string {
	if len(u) < 12 {
		return "****"
	}
	parts := strings.Split(u, "-")
	if len(parts) < 5 {
		return u[:8] + "-xxxx-xxxx-xxxx-..." + u[len(u)-4:]
	}
	return parts[0] + "-xxxx-xxxx-xxxx-..." + parts[4][len(parts[4])-4:]
}

// MaskPath returns only the base filename.
func MaskPath(p string) string { return filepath.Base(p) }

func Timestamp() string { return time.Now().Format("20060102_150405") }
