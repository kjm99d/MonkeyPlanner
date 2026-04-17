package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// runMCP starts the MCP server communicating via stdio JSON-RPC.
// All logging goes to stderr; stdout is reserved for the JSON-RPC protocol.
func runMCP() {
	log.SetOutput(os.Stderr)

	baseURL := strings.TrimRight(getenv("MP_BASE_URL", "http://localhost:8080"), "/")
	client := &mcpHTTPClient{base: baseURL, hc: &http.Client{}}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req jsonRPCRequest
		if err := json.Unmarshal(line, &req); err != nil {
			continue
		}
		resp := handleMCPRequest(client, &req)
		if resp == nil {
			continue
		}
		out, _ := json.Marshal(resp)
		fmt.Fprintln(os.Stdout, string(out))
	}
}
