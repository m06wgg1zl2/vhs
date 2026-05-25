package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Tape represents a VHS tape file containing a sequence of commands
// to be executed for terminal recording.
type Tape struct {
	Commands []Command
}

// Command represents a single instruction parsed from a tape file.
type Command struct {
	Type    CommandType
	Args    string
	Options map[string]string
}

// CommandType enumerates the supported tape command types.
type CommandType int

const (
	CommandUnknown CommandType = iota
	CommandOutput
	CommandSet
	CommandType_Type
	CommandSleep
	CommandEnter
	CommandBackspace
	CommandDelete
	CommandCtrl
	CommandAlt
	CommandScreenshot
	CommandHide
	CommandShow
	CommandWaitFor
	CommandSource
)

// commandNames maps string tokens to CommandType values.
// Note: keys are all lowercase; token matching is case-insensitive.
var commandNames = map[string]CommandType{
	"output":     CommandOutput,
	"set":        CommandSet,
	"type":       CommandType_Type,
	"sleep":      CommandSleep,
	"enter":      CommandEnter,
	"backspace":  CommandBackspace,
	"delete":     CommandDelete,
	"ctrl":       CommandCtrl,
	"alt":        CommandAlt,
	"screenshot": CommandScreenshot,
	"hide":       CommandHide,
	"show":       CommandShow,
	"waitfor":    CommandWaitFor,
	"source":     CommandSource,
	// Personal note: keeping "key" as an alias for common single-key presses
	// would be nice here eventually, but for now this map covers all upstream commands.
}

// ParseTape reads a .tape file from disk and returns a Tape with
// all parsed commands. Lines beginning with '#' are treated as
// comments and ignored. Empty lines are skipped.
//
// Lines may also use '--' as an inline comment delimiter in addition
// to the ' #' style, e.g.: "Sleep 500ms -- wait for prompt"
func ParseTape(path string) (*Tape, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open tape %q: %w", path, err)
	}
	defer f.Close()

	tape := &Tape{}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and blank lines.
		// Also skip inline comments: if a line contains " #" or " --", treat
		// everything from that point onward as a comment.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		if idx := strings.Index(line, " --"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		if line == "" {
			continue
		}

		cmd, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("tape %q line %d: %w", path, lineNum, err)
		}
		tape.Commands = append(tape.Commands, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading tape %q: %w", path, err)
	}

	return tape, nil
}

// parseLine converts a single non-empty, non-comment line into a Command.
func parseLine(line string) (Command, error) {
	parts := strings.SplitN(line, " ", 2)
	token := strings.ToLower(parts[0])

	cmdType, ok := commandNames[token]
	if !ok {
		return Command{}, fmt.Errorf("unknown command %q", parts[0])
	}

	args := ""
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}

	return Command{
		Type: cmdType,
		Args: args,
	}, nil
}
