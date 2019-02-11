package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var (
	minStringLength = 12
	maxStringLength = 24
	parallelism     = runtime.NumCPU()
	interval        = 5
	printMutex      sync.Mutex
)

func init() {
	cli.HelpFlag = cli.BoolFlag{Name: "help, h"}
	cli.VersionFlag = cli.BoolFlag{Name: "version"}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", Version)
	}
}

func main() {
	app := cli.NewApp()

	app.Writer = color.Output
	app.ErrWriter = color.Output

	cli.AppHelpTemplate = fmt.Sprintf("%s\n%s", colorLogo, cli.AppHelpTemplate)
	app.Name = "hasherbasher"
	app.Usage = "Bruteforce SQL Injection MD5 Hashes"
	app.Description = "Tool to bruteforce MD5 hashes that contain SQL injection strings capable of bypassing authentication when not properly santizied or generated (see PHP's MD5() function)."

	app.Version = Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}

	app.Copyright = `(c) 2018 Alex Levinson`
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "bruteforce",
			Usage:     "Start a bruteforce attack.",
			UsageText: "hasherbasher bruteforce",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "min-string-length",
					Usage:       "Minimum length of generated input strings",
					Value:       minStringLength,
					Destination: &minStringLength,
				},
				cli.IntFlag{
					Name:        "max-string-length",
					Usage:       "Maximum length of generated input strings",
					Value:       maxStringLength,
					Destination: &maxStringLength,
				},
				cli.IntFlag{
					Name:        "parallelism",
					Usage:       "Number of parallel brute force workers",
					Value:       parallelism,
					Destination: &parallelism,
				},
				cli.IntFlag{
					Name:        "interval",
					Usage:       "Interval to print statistics in seconds",
					Value:       interval,
					Destination: &interval,
				},
			},
			Action: runCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		Logger.Fatalf("Terminated due to error: %v", err)
	}
}

func runCommand(c *cli.Context) error {
	fmt.Fprintf(color.Output, "%s\n", colorLogo)
	lines := []string{BoldBrightYellow("Configuration")}
	lines = append(lines, fmt.Sprintf("\n %s: %s", BoldWhite("Minimum Length"), BoldBrightRed("%d", minStringLength)))
	lines = append(lines, fmt.Sprintf(" %s: %s", BoldWhite("Maximum Length"), BoldBrightRed("%d", maxStringLength)))
	lines = append(lines, fmt.Sprintf(" %s: %s", BoldWhite("   Parallelism"), BoldBrightRed("%d", parallelism)))
	lines = append(lines, fmt.Sprintf(" %s: %s\n", BoldWhite("Stats Interval"), BoldBrightRed("%d", interval)))
	Logger.Infof(strings.Join(lines, "\n"))
	Logger.Infof("Beginning brute force...")
	currentTime := time.Now()

	finished := false
	resChan := make(chan Result, 1)

	counters := make([]int64, parallelism)

	for i := 0; i < parallelism; i++ {
		counters[i] = int64(1)
		go worker(resChan, &counters[i], minStringLength, maxStringLength, &finished)
	}

	ticker := time.NewTicker(time.Duration(int64(interval)) * time.Second)
	defer ticker.Stop()

	calcTotalFunc := func() int64 {
		total := int64(0)
		for _, x := range counters {
			total = total + x
		}
		return total
	}

	for {
		select {
		case result := <-resChan:
			total := calcTotalFunc()
			printMatch(result, total, parallelism, currentTime)
			return nil
		case <-ticker.C:
			printStats(calcTotalFunc(), parallelism, currentTime)
		}
	}
}

func commandNotImplemented(c *cli.Context) error {
	return fmt.Errorf("%s command not implemented", c.Command.FullName())
}
