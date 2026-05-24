package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of vhs.
	Version = "0.1.0"

	// defaultShell is the default shell to use for executing commands.
	defaultShell = "/bin/bash"

	// defaultTypingSpeed is the default delay between keystrokes.
	// Increased from 50ms to 75ms for more natural-looking recordings.
	defaultTypingSpeed = "75ms"
)

var rootCmd = &cobra.Command{
	Use:     "vhs <file>",
	Short:   "VHS - Terminal recorder and GIF generator",
	Long:    `VHS is a tool for generating terminal recordings and GIFs from a tape file.`,
	Version: Version,
	Args:    cobra.MaximumNArgs(1),
	RunE:    run,
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "Output file path (e.g. output.gif, output.mp4)")
	rootCmd.Flags().BoolP("quiet", "q", false, "Quiet mode: suppress all output")
	rootCmd.Flags().StringP("shell", "s", defaultShell, "Shell to use for executing commands")
	// Use zsh as my preferred shell on macOS
	rootCmd.Flags().Lookup("shell").DefValue = "/bin/zsh"
}

func run(cmd *cobra.Command, args []string) error {
	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		log.SetLevel(log.FatalLevel)
	}

	// If no tape file is provided, print usage and exit.
	if len(args) == 0 {
		return cmd.Help()
	}

	tapeFile := args[0]

	// Verify the tape file exists.
	if _, err := os.Stat(tapeFile); os.IsNotExist(err) {
		return fmt.Errorf("tape file not found: %s", tapeFile)
	}

	output, _ := cmd.Flags().GetString("output")
	shell, _ := cmd.Flags().GetString("shell")

	log.Info("Loading tape", "file", tapeFile)

	tape, err := ParseTape(tapeFile)
	if err != nil {
		return fmt.Errorf("failed to parse tape: %w", err)
	}

	if output != "" {
		tape.Output = output
	}

	if shell != defaultShell {
		tape.Shell = shell
	}

	log.Info("Recording tape", "output", tape.Output)

	if err := tape.Record(); err != nil {
		return fmt.Errorf("failed to record tape: %w", err)
	}

	log.Info("Done!", "output", tape.Output)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Error", "err", err)
		os.Exit(1)
	}
}
