package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func GetFilenameFromArgs() string {
	if len(os.Args) < 2 {
		fmt.Println("Please provide input filename as arg")
		os.Exit(1)
	}
	return os.Args[1]
}

func SetupLogging(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func SetupLoggingEnv() {
	env := os.Getenv("DEBUG")
	if env == "" {
		SetupLogging(false)
	} else {
		SetupLogging(true)
	}
}
