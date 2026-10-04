package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ANSI Color and Style escape codes
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Italic      = "\033[3m"
	Underline   = "\033[4m"

	// Foreground colors
	Black       = "\033[30m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	Gray        = "\033[90m"

	// High Intensity Foreground
	HiBlack     = "\033[90m"
	HiRed       = "\033[91m"
	HiGreen     = "\033[92m"
	HiYellow    = "\033[93m"
	HiBlue      = "\033[94m"
	HiMagenta   = "\033[95m"
	HiCyan      = "\033[96m"
	HiWhite     = "\033[97m"

	// Background colors
	BgBlack     = "\033[40m"
	BgRed       = "\033[41m"
	BgGreen     = "\033[42m"
	BgYellow    = "\033[43m"
	BgBlue      = "\033[44m"
	BgMagenta   = "\033[45m"
	BgCyan      = "\033[46m"
	BgWhite     = "\033[47m"
)

// Color formatting helpers
func Colorize(color, text string) string {
	return color + text + Reset
}

func BoldText(text string) string   { return Bold + text + Reset }
func DimText(text string) string    { return Dim + text + Reset }
func GreenText(text string) string  { return Green + text + Reset }
func CyanText(text string) string   { return Cyan + text + Reset }
func YellowText(text string) string { return Yellow + text + Reset }
func MagentaText(text string) string { return Magenta + text + Reset }
func RedText(text string) string    { return Red + text + Reset }
func BlueText(text string) string   { return Blue + text + Reset }

// Field represents a key-value pair for banner printing
type Field struct {
	Key   string
	Value string
}

// Logger provides structured, colorized terminal logging for CIPHER nodes
type Logger struct {
	mu        sync.Mutex
	out       io.Writer
	role      string
	roleColor string
	subDomain string
}

// New creates a new role-tagged Logger
func New(role string) *Logger {
	roleColor := Cyan
	switch strings.ToUpper(role) {
	case "PUBLISHER":
		roleColor = Magenta
	case "PROVIDER":
		roleColor = Green
	case "CONSUMER", "CLIENT":
		roleColor = Blue
	case "BOOTSTRAP":
		roleColor = Cyan
	case "RELAY":
		roleColor = Yellow
	case "AVAILABILITY":
		roleColor = HiMagenta
	case "PAYMENT", "PAYMENTS":
		roleColor = HiYellow
	}

	return &Logger{
		out:       os.Stdout,
		role:      strings.ToUpper(role),
		roleColor: roleColor,
	}
}

// Sub returns a sub-logger tagged with a specific module/domain (e.g. DHT, PAYMENT, AVAILABILITY)
func (l *Logger) Sub(domain string) *Logger {
	return &Logger{
		out:       l.out,
		role:      l.role,
		roleColor: l.roleColor,
		subDomain: strings.ToUpper(domain),
	}
}

func (l *Logger) formatPrefix(levelTag, levelColor string) string {
	now := time.Now().Format("15:04:05")
	timeStr := Dim + "[" + now + "]" + Reset

	var roleStr string
	if l.role != "" {
		roleStr = Bold + l.roleColor + "[" + l.role + "]" + Reset
	}

	var domainStr string
	if l.subDomain != "" {
		domainColor := Cyan
		switch l.subDomain {
		case "PAYMENT", "PAYMENTS", "SETTLEMENT":
			domainColor = Yellow
		case "AVAILABILITY", "CHALLENGE":
			domainColor = Magenta
		case "DHT", "DISCOVERY":
			domainColor = Cyan
		case "INGEST", "STORAGE", "CAS":
			domainColor = Blue
		case "TRANSFER", "SWARM":
			domainColor = HiCyan
		}
		domainStr = Bold + domainColor + "[" + l.subDomain + "]" + Reset
	}

	var tagStr string
	if levelTag != "" {
		tagStr = Bold + levelColor + levelTag + Reset
	}

	parts := []string{timeStr}
	if roleStr != "" {
		parts = append(parts, roleStr)
	}
	if domainStr != "" {
		parts = append(parts, domainStr)
	}
	if tagStr != "" {
		parts = append(parts, tagStr)
	}

	return strings.Join(parts, " ")
}

func (l *Logger) logWithPrefix(levelTag, levelColor, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	prefix := l.formatPrefix(levelTag, levelColor)
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	fmt.Fprintf(l.out, "%s %s\n", prefix, msg)
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	l.logWithPrefix("[INFO]", Cyan, format, args...)
}

// Success logs a positive operation outcome
func (l *Logger) Success(format string, args ...interface{}) {
	l.logWithPrefix("[✓]", HiGreen, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logWithPrefix("[!]", HiYellow, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.logWithPrefix("[✗]", HiRed, format, args...)
}

// Fatal logs an error and exits the program with code 1
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logWithPrefix("[FATAL]", BgRed+White, format, args...)
	os.Exit(1)
}

// Fatalf is an alias for Fatal
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(format, args...)
}

// Banner prints an eye-catching startup/status banner box
func (l *Logger) Banner(title string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	border := "======================================================================"
	fmt.Fprintf(l.out, "\n%s%s%s\n", Bold, l.roleColor, border)
	fmt.Fprintf(l.out, " %s%s%s\n", Bold, title, Reset)
	fmt.Fprintf(l.out, "%s%s%s\n", Bold, l.roleColor, border)

	for _, f := range fields {
		if f.Key == "" && f.Value == "" {
			fmt.Fprintln(l.out)
			continue
		}
		if f.Key == "" {
			fmt.Fprintf(l.out, "%s\n", f.Value)
			continue
		}
		fmt.Fprintf(l.out, "%s: %s\n", f.Key, f.Value)
	}

	fmt.Fprintf(l.out, "%s%s%s\n\n", Bold, l.roleColor, border)
}



// Global default loggers for convenient quick usage
var (
	Publisher  = New("PUBLISHER")
	Provider   = New("PROVIDER")
	Consumer   = New("CONSUMER")
	Bootstrap  = New("BOOTSTRAP")
	Relay      = New("RELAY")
	Default    = New("CIPHER")
)
