package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/gscript/logger/standard"
	"github.com/sirupsen/logrus"
)

var (
	// Logger is a global singleton logger
	Logger logger.Logger

	// BoldBrightGreen defines a color printer
	BoldBrightGreen = color.New(color.FgHiGreen, color.Bold).SprintfFunc()

	// BoldBrightWhite defines a color printer
	BoldBrightWhite = color.New(color.FgHiWhite, color.Bold).SprintfFunc()

	// BoldBrightRed defines a color printer
	BoldBrightRed = color.New(color.FgHiRed, color.Bold).SprintfFunc()

	// BoldBrightYellow defines a color printer
	BoldBrightYellow = color.New(color.FgHiYellow, color.Bold).SprintfFunc()

	// BoldBrightCyan defines a color printer
	BoldBrightCyan = color.New(color.FgHiCyan, color.Bold).SprintfFunc()

	// BoldBrightBlue defines a color printer
	BoldBrightBlue = color.New(color.FgHiBlue, color.Bold).SprintfFunc()

	// BoldBrightMagenta defines a color printer
	BoldBrightMagenta = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()

	// BrightGreen defines a color printer
	BrightGreen = color.New(color.FgHiGreen).SprintfFunc()

	// BrightWhite defines a color printer
	BrightWhite = color.New(color.FgHiWhite).SprintfFunc()

	// BrightRed defines a color printer
	BrightRed = color.New(color.FgHiRed).SprintfFunc()

	// BrightYellow defines a color printer
	BrightYellow = color.New(color.FgHiYellow).SprintfFunc()

	// BrightCyan defines a color printer
	BrightCyan = color.New(color.FgHiCyan).SprintfFunc()

	// BrightBlue defines a color printer
	BrightBlue = color.New(color.FgHiBlue).SprintfFunc()

	// BrightMagenta defines a color printer
	BrightMagenta = color.New(color.FgHiMagenta).SprintfFunc()

	// BoldGreen defines a color printer
	BoldGreen = color.New(color.FgGreen, color.Bold).SprintfFunc()

	// BoldWhite defines a color printer
	BoldWhite = color.New(color.FgWhite, color.Bold).SprintfFunc()

	// BoldRed defines a color printer
	BoldRed = color.New(color.FgRed, color.Bold).SprintfFunc()

	// BoldYellow defines a color printer
	BoldYellow = color.New(color.FgYellow, color.Bold).SprintfFunc()

	// BoldCyan defines a color printer
	BoldCyan = color.New(color.FgCyan, color.Bold).SprintfFunc()

	// BoldBlue defines a color printer
	BoldBlue = color.New(color.FgBlue, color.Bold).SprintfFunc()

	// BoldMagenta defines a color printer
	BoldMagenta = color.New(color.FgMagenta, color.Bold).SprintfFunc()

	// Green defines a color printer
	Green = color.New(color.FgGreen).SprintfFunc()

	// White defines a color printer
	White = color.New(color.FgWhite).SprintfFunc()

	// Red defines a color printer
	Red = color.New(color.FgRed).SprintfFunc()

	// Yellow defines a color printer
	Yellow = color.New(color.FgYellow).SprintfFunc()

	// Cyan defines a color printer
	Cyan = color.New(color.FgCyan).SprintfFunc()

	// Blue defines a color printer
	Blue = color.New(color.FgBlue).SprintfFunc()

	// Magenta defines a color printer
	Magenta = color.New(color.FgMagenta).SprintfFunc()

	// NoColor defines a color printer
	NoColor = color.New(color.Reset).SprintfFunc()
)

var (
	intLogger    *internalLogger
	globalProg   = `HASHERBASHER`
	startName    = `cli`
	defaultLevel = logrus.InfoLevel
)

type internalLogger struct {
	internal *logrus.Logger
	writer   *logWriter
	prog     string
	context  string
}

type logWriter struct {
	Name string
	Prog string
}

func init() {
	base := standard.NewStandardLogger(nil, globalProg, startName, false, false)
	baseSL := base.Logger
	writer := &logWriter{Prog: globalProg, Name: startName}
	baseSL.Out = writer
	baseSL.SetLevel(defaultLevel)
	logger := &internalLogger{
		internal: baseSL,
		writer:   writer,
	}
	intLogger = logger
	Logger = base
}

// SetLogLevel allows you to override the logging level for the Laforge global logger
func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		intLogger.internal.SetLevel(logrus.DebugLevel)
	case "info":
		intLogger.internal.SetLevel(logrus.InfoLevel)
	case "warn":
		intLogger.internal.SetLevel(logrus.WarnLevel)
	case "error":
		intLogger.internal.SetLevel(logrus.ErrorLevel)
	case "fatal":
		intLogger.internal.SetLevel(logrus.FatalLevel)
	}
}

// SetLogName allows you to override the log name parameter (part after LAFORGE in log output)
func SetLogName(name string) {
	intLogger.writer.Name = name
}

func (w *logWriter) Write(p []byte) (int, error) {
	output := fmt.Sprintf(
		"%s%s%s%s%s %s",
		BrightWhite("["),
		BoldBrightRed(w.Prog),
		BrightWhite(":"),
		BrightRed(strings.ToLower(w.Name)),
		BrightWhite("]"),
		string(p),
	)
	written, err := io.Copy(color.Output, strings.NewReader(output))
	return int(written), err
}

func printMatch(result Result, count int64, workerCount int, begin time.Time) {
	printMutex.Lock()
	defer printMutex.Unlock()
	curr := time.Now()
	duration := curr.Sub(begin)
	perSec := float64(count) / duration.Seconds()
	lines := []string{}
	lines = append(lines, BoldBrightCyan("Statistics"))
	lines = append(lines, fmt.Sprintf("\n       %s: %s", BoldWhite("Start Time"), BoldYellow(begin.Format(time.RFC822Z))))
	lines = append(lines, fmt.Sprintf(" %s: %s", BoldWhite("Elapsed Duration"), BoldYellow(humanize.RelTime(begin, curr, "", ""))))

	lines = append(lines, fmt.Sprintf("   %s: %s", BoldWhite("Total Attempts"), BoldYellow(humanize.Comma(count))))
	lines = append(lines, fmt.Sprintf("       %s: %s per second", BoldWhite("Crack Rate"), BoldYellow(humanize.CommafWithDigits(perSec, 2))))
	lines = append(lines, fmt.Sprintf("       %s: %s per worker per second\n", BoldWhite("Per Worker"), BoldYellow(humanize.CommafWithDigits(perSec/float64(workerCount), 2))))
	Logger.Infof(strings.Join(lines, "\n"))

	matchedStringAsBytes := []byte(result.Output)
	inputString := result.Input

	Logger.Infof("%s %s %s", BoldBrightWhite("====="), BoldBrightGreen("Match Found"), BoldBrightWhite("====="))
	Logger.Infof("Cracked In: %s", BoldBrightCyan("%s seconds", humanize.Commaf(duration.Seconds())))
	Logger.Infof(" -- BEGIN RAW BYTES --")
	fmt.Printf("%s\n", matchedStringAsBytes)
	Logger.Infof(" -- END RAW BYTES --")
	lines = []string{fmt.Sprintf("%s %s %s", BoldBrightWhite("====="), BoldBrightGreen("Results"), BoldBrightWhite("====="))}

	lines = append(lines, fmt.Sprintf("\n %s: %s", BoldBrightWhite("Located String"), BoldBrightGreen("%s", inputString)))
	lines = append(lines, fmt.Sprintf("    %s: %s", BoldBrightWhite("Result Size"), BoldBrightGreen("%d", len(matchedStringAsBytes))))
	lines = append(lines, fmt.Sprintf("   %s: %s", BoldBrightWhite("Result Bytes"), BoldBrightGreen("%v", matchedStringAsBytes)))
	lines = append(lines, fmt.Sprintf("     %s: %s\n", BoldBrightWhite("Result Hex"), BoldBrightGreen("%x", matchedStringAsBytes)))
	Logger.Infof(strings.Join(lines, "\n"))

}

func printStats(count int64, workerCount int, begin time.Time) {
	printMutex.Lock()
	defer printMutex.Unlock()
	curr := time.Now()
	duration := curr.Sub(begin)
	perSec := float64(count) / duration.Seconds()

	lines := []string{}
	lines = append(lines, BoldBrightCyan("Statistics"))
	lines = append(lines, fmt.Sprintf("\n       %s: %s", BoldWhite("Start Time"), BoldYellow(begin.Format(time.RFC822Z))))
	lines = append(lines, fmt.Sprintf(" %s: %s", BoldWhite("Elapsed Duration"), BoldYellow(humanize.RelTime(begin, curr, "", ""))))

	lines = append(lines, fmt.Sprintf("   %s: %s", BoldWhite("Total Attempts"), BoldYellow(humanize.Comma(count))))
	lines = append(lines, fmt.Sprintf("       %s: %s per second", BoldWhite("Crack Rate"), BoldYellow(humanize.CommafWithDigits(perSec, 2))))
	lines = append(lines, fmt.Sprintf("       %s: %s per worker per second\n", BoldWhite("Per Worker"), BoldYellow(humanize.CommafWithDigits(perSec/float64(workerCount), 2))))
	Logger.Infof(strings.Join(lines, "\n"))
}
