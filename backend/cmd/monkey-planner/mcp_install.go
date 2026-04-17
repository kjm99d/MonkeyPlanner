package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// runMCPInstall handles `monkey-planner mcp install [--for CLIENT] [...]`.
// It writes (or merges) a MonkeyPlanner MCP server entry into the target
// client's JSON config so users do not have to hand-edit paths.
//
// Usage:
//
//	monkey-planner mcp install --for claude-code
//	monkey-planner mcp install --for claude-desktop
//	monkey-planner mcp install --for cursor
//	monkey-planner mcp install --for claude-code --base-url http://localhost:8080
//	monkey-planner mcp install --for claude-code --scope user   # ~/.mcp.json
func runMCPInstall(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	client := fs.String("for", "", "client to configure: claude-code | claude-desktop | cursor")
	baseURL := fs.String("base-url", "http://localhost:8080", "MonkeyPlanner HTTP base URL the MCP server should talk to")
	scope := fs.String("scope", "project", "where to install (claude-code/cursor only): project | user")
	name := fs.String("name", "monkey-planner", "MCP server entry name")
	dryRun := fs.Bool("dry-run", false, "print the resulting config without writing it")
	force := fs.Bool("force", false, "overwrite an existing entry with the same --name")
	_ = fs.Parse(args)

	if *client == "" {
		fmt.Fprintln(os.Stderr, "error: --for is required (claude-code | claude-desktop | cursor)")
		fs.Usage()
		os.Exit(2)
	}

	binPath, err := os.Executable()
	if err != nil {
		exitf("resolve executable path: %v", err)
	}
	binPath, _ = filepath.Abs(binPath)

	configPath, err := mcpConfigPath(*client, *scope)
	if err != nil {
		exitf("%v", err)
	}

	entry := map[string]any{
		"command": binPath,
		"args":    []string{"mcp"},
		"env": map[string]string{
			"MP_BASE_URL": *baseURL,
		},
	}

	merged, err := mergeMCPEntry(configPath, *name, entry, *force)
	if err != nil {
		exitf("%v", err)
	}

	pretty, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		exitf("encode config: %v", err)
	}
	pretty = append(pretty, '\n')

	if *dryRun {
		fmt.Printf("# would write %s\n%s", configPath, pretty)
		return
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		exitf("mkdir %s: %v", filepath.Dir(configPath), err)
	}
	if err := os.WriteFile(configPath, pretty, 0o644); err != nil {
		exitf("write %s: %v", configPath, err)
	}

	fmt.Printf("✓ Installed MonkeyPlanner MCP server as %q in:\n  %s\n", *name, configPath)
	fmt.Printf("  command: %s mcp\n", binPath)
	fmt.Printf("  env.MP_BASE_URL: %s\n", *baseURL)
	fmt.Println()
	fmt.Println("Next: make sure the HTTP server is running on the same address, e.g.:")
	fmt.Printf("  %s\n", binPath)
	fmt.Println()
	fmt.Printf("Then restart %s to pick up the new MCP server.\n", *client)
}

// mcpConfigPath resolves the correct config file for the requested client.
func mcpConfigPath(client, scope string) (string, error) {
	home, _ := os.UserHomeDir()

	switch client {
	case "claude-code":
		// Claude Code looks for .mcp.json in the current working directory (project
		// scope) or in $HOME (user scope).
		if scope == "user" {
			if home == "" {
				return "", fmt.Errorf("cannot resolve $HOME for user scope")
			}
			return filepath.Join(home, ".mcp.json"), nil
		}
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getwd: %w", err)
		}
		return filepath.Join(cwd, ".mcp.json"), nil

	case "cursor":
		// Cursor mirrors the Claude Code layout under .cursor/.
		if scope == "user" {
			if home == "" {
				return "", fmt.Errorf("cannot resolve $HOME for user scope")
			}
			return filepath.Join(home, ".cursor", "mcp.json"), nil
		}
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getwd: %w", err)
		}
		return filepath.Join(cwd, ".cursor", "mcp.json"), nil

	case "claude-desktop":
		// Claude Desktop always uses the OS-native app config directory.
		switch runtime.GOOS {
		case "darwin":
			return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"), nil
		case "windows":
			appData := os.Getenv("APPDATA")
			if appData == "" {
				appData = filepath.Join(home, "AppData", "Roaming")
			}
			return filepath.Join(appData, "Claude", "claude_desktop_config.json"), nil
		case "linux":
			xdg := os.Getenv("XDG_CONFIG_HOME")
			if xdg == "" {
				xdg = filepath.Join(home, ".config")
			}
			return filepath.Join(xdg, "Claude", "claude_desktop_config.json"), nil
		default:
			return "", fmt.Errorf("claude-desktop on %s is not supported yet", runtime.GOOS)
		}

	default:
		return "", fmt.Errorf("unknown client %q (expected: claude-code | claude-desktop | cursor)", client)
	}
}

// mergeMCPEntry loads an existing config (if any) and inserts/updates the
// given server entry under "mcpServers".
func mergeMCPEntry(path, name string, entry map[string]any, force bool) (map[string]any, error) {
	cfg := map[string]any{}
	if data, err := os.ReadFile(path); err == nil && len(strings.TrimSpace(string(data))) > 0 {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parse existing %s: %w", path, err)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	servers, _ := cfg["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}
	if _, exists := servers[name]; exists && !force {
		return nil, fmt.Errorf("entry %q already exists in %s; pass --force to overwrite", name, path)
	}
	servers[name] = entry
	cfg["mcpServers"] = servers
	return cfg, nil
}

func exitf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}
