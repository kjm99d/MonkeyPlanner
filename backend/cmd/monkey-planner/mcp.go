package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// runMCP starts the MCP server communicating via stdio JSON-RPC.
func runMCP() {
	// All logging to stderr; stdout is reserved for JSON-RPC protocol.
	log.SetOutput(os.Stderr)

	baseURL := strings.TrimRight(getenv("MP_BASE_URL", "http://localhost:8080"), "/")
	client := &mcpHTTPClient{base: baseURL, hc: &http.Client{}}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer

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
			// Notification: no response needed.
			continue
		}

		out, _ := json.Marshal(resp)
		fmt.Fprintln(os.Stdout, string(out))
	}
}

// mcpHTTPClient forwards MCP tool calls to the running HTTP server.
// This ensures SSE events are published through the server's broker.
type mcpHTTPClient struct {
	base string
	hc   *http.Client
}

func (c *mcpHTTPClient) get(path string, q url.Values) (json.RawMessage, error) {
	u := c.base + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	resp, err := c.hc.Get(u)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	return readHTTPBody(resp)
}

func (c *mcpHTTPClient) post(path string, body any) (json.RawMessage, error) {
	return c.doJSON("POST", path, body)
}

func (c *mcpHTTPClient) patch(path string, body any) (json.RawMessage, error) {
	return c.doJSON("PATCH", path, body)
}

func (c *mcpHTTPClient) doJSON(method, path string, body any) (json.RawMessage, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, c.base+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	return readHTTPBody(resp)
}

func readHTTPBody(resp *http.Response) (json.RawMessage, error) {
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return json.RawMessage(b), nil
}

// ---- JSON-RPC types ----

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// ---- Request router ----

func handleMCPRequest(client *mcpHTTPClient, req *jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{}},
				"serverInfo":      map[string]any{"name": "monkey-planner", "version": version},
			},
		}

	case "notifications/initialized":
		return nil // notification, no response

	case "tools/list":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]any{"tools": mcpToolDefinitions()},
		}

	case "tools/call":
		return handleMCPToolCall(client, req)

	default:
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   map[string]any{"code": -32601, "message": "method not found"},
		}
	}
}

// ---- Tool definitions ----

func mcpToolDefinitions() []map[string]any {
	return []map[string]any{
		{
			"name":        "get_version",
			"description": "Get the MCP server version. Use this to check if the server needs a restart after an update.",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			"name":        "list_boards",
			"description": "List all boards",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			"name":        "list_issues",
			"description": "List issues. Filter by boardId and/or status (Pending, Approved, InProgress, QA, Done, Rejected)",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"boardId": map[string]any{"type": "string", "description": "Filter by board ID"},
					"status":  map[string]any{"type": "string", "description": "Filter by status: Pending, Approved, InProgress, Done, Rejected"},
				},
			},
		},
		{
			"name":        "get_issue",
			"description": "Get full issue detail including instructions, criteria, and comments",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID"},
				},
				"required": []string{"issueId"},
			},
		},
		{
			"name":        "create_issue",
			"description": "Create a new issue on a board",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"boardId":      map[string]any{"type": "string", "description": "Board ID"},
					"title":        map[string]any{"type": "string", "description": "Issue title"},
					"body":         map[string]any{"type": "string", "description": "Issue body/description"},
					"instructions": map[string]any{"type": "string", "description": "Implementation instructions"},
				},
				"required": []string{"boardId", "title"},
			},
		},
		{
			"name":        "approve_issue",
			"description": "Approve a pending issue (moves from Pending to Approved)",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID to approve"},
				},
				"required": []string{"issueId"},
			},
		},
		{
			"name":        "claim_issue",
			"description": "Claim an approved issue and move it to InProgress",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID to claim"},
				},
				"required": []string{"issueId"},
			},
		},
		{
			"name":        "submit_qa",
			"description": "Submit an in-progress issue for QA review (moves from InProgress to QA)",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID to submit for QA"},
					"comment": map[string]any{"type": "string", "description": "Optional summary of what was done"},
				},
				"required": []string{"issueId"},
			},
		},
		{
			"name":        "complete_issue",
			"description": "Complete a QA-reviewed issue (moves from QA to Done). Optionally add a completion summary as a comment.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID to complete"},
					"comment": map[string]any{"type": "string", "description": "Optional completion summary comment"},
				},
				"required": []string{"issueId"},
			},
		},
		{
			"name":        "reject_issue",
			"description": "Reject a QA issue back to InProgress with a required reason (moves from QA to InProgress)",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID to reject"},
					"reason":  map[string]any{"type": "string", "description": "Reason for rejection (required)"},
				},
				"required": []string{"issueId", "reason"},
			},
		},
		{
			"name":        "add_comment",
			"description": "Add a comment to an issue",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID"},
					"body":    map[string]any{"type": "string", "description": "Comment body"},
				},
				"required": []string{"issueId", "body"},
			},
		},
		{
			"name":        "update_criteria",
			"description": "Check or uncheck a success criterion on an issue",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"issueId": map[string]any{"type": "string", "description": "Issue ID"},
					"index":   map[string]any{"type": "integer", "description": "Zero-based index of the criterion"},
					"done":    map[string]any{"type": "boolean", "description": "Whether the criterion is done"},
				},
				"required": []string{"issueId", "index", "done"},
			},
		},
		{
			"name":        "search_issues",
			"description": "Search issues by title (case-insensitive substring match)",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{"type": "string", "description": "Search query to match against issue titles"},
				},
				"required": []string{"query"},
			},
		},
	}
}

// ---- Tool call handler ----

func handleMCPToolCall(client *mcpHTTPClient, req *jsonRPCRequest) *jsonRPCResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   map[string]any{"code": -32602, "message": "invalid params"},
		}
	}

	result, err := mcpCallTool(client, params.Name, params.Arguments)
	if err != nil {
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"content": []map[string]any{{"type": "text", "text": "Error: " + err.Error()}},
				"isError": true,
			},
		}
	}

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	return &jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]any{
			"content": []map[string]any{{"type": "text", "text": string(jsonBytes)}},
		},
	}
}

func mcpCallTool(c *mcpHTTPClient, name string, argsRaw json.RawMessage) (any, error) {
	switch name {
	case "get_version":
		return map[string]any{"version": version, "name": "monkey-planner"}, nil

	case "list_boards":
		raw, err := c.get("/api/boards", nil)
		if err != nil {
			return nil, err
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "list_issues":
		var args struct {
			BoardID string `json:"boardId"`
			Status  string `json:"status"`
		}
		_ = json.Unmarshal(argsRaw, &args)
		q := url.Values{}
		if args.BoardID != "" {
			q.Set("board_id", args.BoardID)
		}
		if args.Status != "" {
			q.Set("status", args.Status)
		}
		raw, err := c.get("/api/issues", q)
		if err != nil {
			return nil, err
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "get_issue":
		var args struct {
			IssueID string `json:"issueId"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		issueRaw, err := c.get("/api/issues/"+args.IssueID, nil)
		if err != nil {
			return nil, err
		}
		commentsRaw, err := c.get("/api/issues/"+args.IssueID+"/comments", nil)
		if err != nil {
			return nil, err
		}
		// issueRaw is {"issue": {...}, "children": [...]}
		var issueData map[string]any
		var comments any
		_ = json.Unmarshal(issueRaw, &issueData)
		_ = json.Unmarshal(commentsRaw, &comments)
		issueData["comments"] = comments
		return issueData, nil

	case "create_issue":
		var args struct {
			BoardID      string `json:"boardId"`
			Title        string `json:"title"`
			Body         string `json:"body"`
			Instructions string `json:"instructions"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil {
			return nil, fmt.Errorf("invalid arguments: %w", err)
		}
		issueRaw, err := c.post("/api/issues", map[string]any{
			"boardId": args.BoardID,
			"title":   args.Title,
			"body":    args.Body,
		})
		if err != nil {
			return nil, err
		}
		if args.Instructions != "" {
			var issue struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(issueRaw, &issue); err == nil && issue.ID != "" {
				issueRaw, err = c.patch("/api/issues/"+issue.ID, map[string]any{
					"instructions": args.Instructions,
				})
				if err != nil {
					return nil, err
				}
			}
		}
		var result any
		_ = json.Unmarshal(issueRaw, &result)
		return result, nil

	case "approve_issue":
		var args struct {
			IssueID string `json:"issueId"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		raw, err := c.post("/api/issues/"+args.IssueID+"/approve", map[string]any{})
		if err != nil {
			return nil, err
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "claim_issue":
		var args struct {
			IssueID string `json:"issueId"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "InProgress"})
		if err != nil {
			return nil, err
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "submit_qa":
		var args struct {
			IssueID string `json:"issueId"`
			Comment string `json:"comment"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "QA"})
		if err != nil {
			return nil, err
		}
		if args.Comment != "" {
			_, _ = c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": args.Comment})
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "complete_issue":
		var args struct {
			IssueID string `json:"issueId"`
			Comment string `json:"comment"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "Done"})
		if err != nil {
			return nil, err
		}
		if args.Comment != "" {
			_, _ = c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": args.Comment})
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "reject_issue":
		var args struct {
			IssueID string `json:"issueId"`
			Reason  string `json:"reason"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		if args.Reason == "" {
			return nil, fmt.Errorf("reason is required for rejection")
		}
		raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "InProgress"})
		if err != nil {
			return nil, err
		}
		_, _ = c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": "❌ QA 거절: " + args.Reason})
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "add_comment":
		var args struct {
			IssueID string `json:"issueId"`
			Body    string `json:"body"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" || args.Body == "" {
			return nil, fmt.Errorf("issueId and body are required")
		}
		raw, err := c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": args.Body})
		if err != nil {
			return nil, err
		}
		var result any
		_ = json.Unmarshal(raw, &result)
		return result, nil

	case "update_criteria":
		var args struct {
			IssueID string `json:"issueId"`
			Index   int    `json:"index"`
			Done    bool   `json:"done"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId, index, and done are required")
		}
		issueRaw, err := c.get("/api/issues/"+args.IssueID, nil)
		if err != nil {
			return nil, err
		}
		var issueResp struct {
			Issue struct {
				Criteria []map[string]any `json:"criteria"`
			} `json:"issue"`
		}
		if err := json.Unmarshal(issueRaw, &issueResp); err != nil {
			return nil, fmt.Errorf("failed to parse issue: %w", err)
		}
		criteria := issueResp.Issue.Criteria
		if args.Index < 0 || args.Index >= len(criteria) {
			return nil, fmt.Errorf("index %d out of range (issue has %d criteria)", args.Index, len(criteria))
		}
		criteria[args.Index]["done"] = args.Done
		raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"criteria": criteria})
		if err != nil {
			return nil, err
		}
		var updated struct {
			Criteria any `json:"criteria"`
		}
		_ = json.Unmarshal(raw, &updated)
		return updated.Criteria, nil

	case "search_issues":
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.Query == "" {
			return nil, fmt.Errorf("query is required")
		}
		raw, err := c.get("/api/issues", nil)
		if err != nil {
			return nil, err
		}
		var all []map[string]any
		if err := json.Unmarshal(raw, &all); err != nil {
			return nil, err
		}
		q := strings.ToLower(args.Query)
		var matched []map[string]any
		for _, iss := range all {
			if title, ok := iss["title"].(string); ok && strings.Contains(strings.ToLower(title), q) {
				matched = append(matched, iss)
			}
		}
		return matched, nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}
