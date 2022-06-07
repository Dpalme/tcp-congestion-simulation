/*
Copyright © 2022 Diego Palmerin dpalme@me.com

*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	congestionalgos "tcp-congestion/pkg/congestionAlgos"
	congestiontest "tcp-congestion/pkg/congestionTest"
	"tcp-congestion/pkg/logger"
	"tcp-congestion/pkg/simulation"
	"tcp-congestion/pkg/utils"

	"github.com/spf13/cobra"
)

type TCPOptions struct {
	logLevel string
	output   string
	verbose  bool
	debug    bool
	runs     int
}

var (
	opts = TCPOptions{
		logLevel: "",
		output:   "",
		debug:    false,
		verbose:  false,
		runs:     1,
	}
	appLogger *logger.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "tcp-congestion {simulate|test} {TCP|BIC|CTCP} {cc}",
	TraverseChildren: true,
	Args:             cobra.MinimumNArgs(3),
	ValidArgs:        []string{"simulate", "test"},
	Short:            "Test TCP Congestion Algorithms",
	Long:             `Simulate or stress test different implementations of TCP congestion algorithms.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func determineVerbosity() int {
	if opts.verbose {
		return 2
	}
	if opts.debug {
		return 3
	}
	switch opts.logLevel {
	case "DEBUG":
		return 3
	case "INFO":
		return 2
	case "DEFAULT":
		return 1
	default:
		return 0
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.output, "output", "o", "", "Output file. Default is stdout")
	rootCmd.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "Alias for --log-level=INFO")
	rootCmd.PersistentFlags().BoolVarP(&opts.debug, "debug", "d", false, "Alias for --log-level=DEBUG")
	rootCmd.PersistentFlags().StringVarP(&opts.logLevel, "log-level", "l", "DEFAULT", "Specifies the loggers sensitivity. Can be one of: CRITICAL, DEFAULT, INFO, DEBUG")
	rootCmd.ParseFlags(os.Args)

	simulateCmd := &cobra.Command{
		Use:       "simulate {TCP|BIC|CTCP} {cc}",
		Short:     "Run the algorithm to your hearts content and log as you go",
		Long:      `Run the algorithm to your hearts content and log as you go`,
		ValidArgs: congestionalgos.RegisteredAlgorithms,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("missing required arguments")
			}

			_, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(args)
			clientCount, _ := strconv.ParseInt(args[1], 10, 32)
			simulation.RunFromCLI(args[0], int(clientCount), BuildLogger())
		},
	}
	rootCmd.AddCommand(simulateCmd)

	testCmd := &cobra.Command{
		Use:   "test {TCP|BIC|CTCP} {cc} {duration}",
		Short: "Stress-test the algorithm",
		Long: `Run the algorithm for a specified duration a given
		number of times and analyze its behaviour at scale`,
		ValidArgs: congestionalgos.RegisteredAlgorithms,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("missing required arguments")
			}

			_, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				return err
			}
			_, err = strconv.ParseFloat(args[2], 64)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.LocalFlags().Parse(args)
			cc, err := congestionalgos.GetByName(args[0])
			if err != nil {
				appLogger.Critical(err.Error())
				return
			}
			clientCount, _ := strconv.ParseInt(args[1], 10, 32)
			duration, _ := strconv.ParseFloat(args[2], 64)
			congestiontest.New(cc, clientCount, duration, BuildLogger()).RunNTimes(opts.runs)
		},
	}

	testCmd.LocalFlags().IntVarP(&opts.runs, "runs", "r", 1, "Number of times to run the test for.")
	testCmd.LocalFlags().Parse(os.Args)
	rootCmd.AddCommand(testCmd)
}

func BuildLogger() *logger.Logger {
	verbosity := determineVerbosity()
	appLogger = logger.New("", verbosity)
	if opts.output != "" {
		appLogger.File = utils.NewLogFile(opts.output)
	}
	return appLogger
}
