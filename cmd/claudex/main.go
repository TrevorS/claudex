// ABOUTME: Entry point for claudex TUI application
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TrevorS/claudex/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define command-line flags
	configPath := flag.String("config", "", "Path to config file")
	projectsPath := flag.String("projects", "", "Override Claude projects path")
	version := flag.Bool("version", false, "Print version information")
	help := flag.Bool("help", false, "Print help information")

	flag.Parse()

	// Handle --help flag
	if *help {
		printHelp()
		os.Exit(0)
	}

	// Handle --version flag
	if *version {
		printVersion()
		os.Exit(0)
	}

	// Load configuration
	config, err := app.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override with command-line flags if provided
	if *projectsPath != "" {
		config.ClaudeProjectsPath = *projectsPath
	}

	_ = configPath // For future use with config file loading

	// Initialize Bubble Tea debug logging if DEBUG is set
	if os.Getenv("DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "claudex")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not setup debug log: %v\n", err)
		} else {
			defer f.Close()
		}
	}

	// Initialize application
	application, err := app.New(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Run the application
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("claudex version %s\n", Version)
	fmt.Printf("Built: %s\n", BuildDate)
	fmt.Printf("Commit: %s\n", GitCommit)
}

func printHelp() {
	fmt.Println("claudex - Claude conversation history browser")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  claudex [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --projects PATH    Override Claude projects path")
	fmt.Println("  --config PATH      Path to config file")
	fmt.Println("  --version          Print version information")
	fmt.Println("  --help             Print this help message")
	fmt.Println()
	fmt.Println("Keyboard Shortcuts:")
	fmt.Println("  ^C                 Quit application")
	fmt.Println("  Tab                Switch focus between panes")
	fmt.Println("  Escape             Return to conversation list")
	fmt.Println("  ^P                 Open command palette (Phase 2+)")
}
