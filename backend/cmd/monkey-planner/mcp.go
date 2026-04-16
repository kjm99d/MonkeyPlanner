package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

// runMCP starts the MCP server communicating via stdio JSON-RPC.
func runMCP() {
	// All logging to stderr; stdout is reserved for JSON-RPC protocol.
	log.SetOutput(os.Stderr)

	dsn := getenv("MP_DSN", defaultDSN())

	repo, err := storage.NewRepo(dsn)
	if err != nil {
		log.Fatalf("mcp: storage open: %v", err)
	}
	defer repo.Close()

	svc := service.New(repo, nil)
	ctx := context.Background()

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

		resp := handleMCPRequest(ctx, svc, &req)
		if resp == nil {
			// Notification: no response needed.
			continue
		}

		out, _ := json.Marshal(resp)
		fmt.Fprintln(os.Stdout, string(out))
	}
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

func handleMCPRequest(ctx context.Context, svc *service.Service, req *jsonRPCRequest) *jsonRPCResponse {
	switch req.Method {
	case "initialize":
		return &jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":   map[string]any{"tools": map[string]any{}},
				"serverInfo":     map[string]any{"name": "monkey-planner", "version": version},
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
		return handleMCPToolCall(ctx, svc, req)

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

func handleMCPToolCall(ctx context.Context, svc *service.Service, req *jsonRPCRequest) *jsonRPCResponse {
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

	result, err := mcpCallTool(ctx, svc, params.Name, params.Arguments)
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

func mcpCallTool(ctx context.Context, svc *service.Service, name string, argsRaw json.RawMessage) (any, error) {
	switch name {
	case "get_version":
		return map[string]any{"version": version, "name": "monkey-planner"}, nil

	case "list_boards":
		return svc.ListBoards(ctx)

	case "list_issues":
		var args struct {
			BoardID string `json:"boardId"`
			Status  string `json:"status"`
		}
		_ = json.Unmarshal(argsRaw, &args)

		var f storage.IssueFilter
		if args.BoardID != "" {
			f.BoardID = &args.BoardID
		}
		if args.Status != "" {
			st := domain.Status(args.Status)
			f.Status = &st
		}
		return svc.ListIssues(ctx, f)

	case "get_issue":
		var args struct {
			IssueID string `json:"issueId"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		issue, children, err := svc.GetIssue(ctx, args.IssueID)
		if err != nil {
			return nil, err
		}
		comments, _ := svc.ListComments(ctx, args.IssueID)
		return map[string]any{
			"issue":    issue,
			"children": children,
			"comments": comments,
		}, nil

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
		issue, err := svc.CreateIssue(ctx, service.CreateIssueInput{
			BoardID: args.BoardID,
			Title:   args.Title,
			Body:    args.Body,
		})
		if err != nil {
			return nil, err
		}
		// If instructions were provided, update them separately since CreateIssueInput doesn't have Instructions.
		if args.Instructions != "" {
			issue, err = svc.UpdateIssue(ctx, issue.ID, service.UpdateIssueInput{
				Instructions: &args.Instructions,
			})
			if err != nil {
				return nil, err
			}
		}
		return issue, nil

	case "approve_issue":
		var args struct {
			IssueID string `json:"issueId"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		return svc.ApproveIssue(ctx, args.IssueID)

	case "claim_issue":
		var args struct {
			IssueID string `json:"issueId"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		st := domain.StatusInProgress
		return svc.UpdateIssue(ctx, args.IssueID, service.UpdateIssueInput{
			Status: &st,
		})

	case "submit_qa":
		var args struct {
			IssueID string `json:"issueId"`
			Comment string `json:"comment"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		st := domain.StatusQA
		issue, err := svc.UpdateIssue(ctx, args.IssueID, service.UpdateIssueInput{
			Status: &st,
		})
		if err != nil {
			return nil, err
		}
		if args.Comment != "" {
			_, _ = svc.CreateComment(ctx, args.IssueID, args.Comment)
		}
		return issue, nil

	case "complete_issue":
		var args struct {
			IssueID string `json:"issueId"`
			Comment string `json:"comment"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId is required")
		}
		issue, err := svc.CompleteIssue(ctx, args.IssueID)
		if err != nil {
			return nil, err
		}
		if args.Comment != "" {
			_, _ = svc.CreateComment(ctx, args.IssueID, args.Comment)
		}
		return issue, nil

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
		st := domain.StatusInProgress
		issue, err := svc.UpdateIssue(ctx, args.IssueID, service.UpdateIssueInput{
			Status: &st,
		})
		if err != nil {
			return nil, err
		}
		_, _ = svc.CreateComment(ctx, args.IssueID, "❌ QA 거절: "+args.Reason)
		return issue, nil

	case "add_comment":
		var args struct {
			IssueID string `json:"issueId"`
			Body    string `json:"body"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" || args.Body == "" {
			return nil, fmt.Errorf("issueId and body are required")
		}
		return svc.CreateComment(ctx, args.IssueID, args.Body)

	case "update_criteria":
		var args struct {
			IssueID string `json:"issueId"`
			Index   int    `json:"index"`
			Done    bool   `json:"done"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.IssueID == "" {
			return nil, fmt.Errorf("issueId, index, and done are required")
		}
		// Get current issue to read existing criteria
		issue, _, err := svc.GetIssue(ctx, args.IssueID)
		if err != nil {
			return nil, err
		}
		if args.Index < 0 || args.Index >= len(issue.Criteria) {
			return nil, fmt.Errorf("index %d out of range (issue has %d criteria)", args.Index, len(issue.Criteria))
		}
		criteria := make([]domain.Criterion, len(issue.Criteria))
		copy(criteria, issue.Criteria)
		criteria[args.Index].Done = args.Done
		updated, err := svc.UpdateIssue(ctx, args.IssueID, service.UpdateIssueInput{
			Criteria: &criteria,
		})
		if err != nil {
			return nil, err
		}
		return updated.Criteria, nil

	case "search_issues":
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(argsRaw, &args); err != nil || args.Query == "" {
			return nil, fmt.Errorf("query is required")
		}
		// List all issues and filter by title substring match.
		all, err := svc.ListIssues(ctx, storage.IssueFilter{})
		if err != nil {
			return nil, err
		}
		q := strings.ToLower(args.Query)
		var matched []domain.Issue
		for _, iss := range all {
			if strings.Contains(strings.ToLower(iss.Title), q) {
				matched = append(matched, iss)
			}
		}
		return matched, nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}
