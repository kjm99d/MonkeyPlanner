package main

import (
	"encoding/json"
)

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

// ---- Tool call dispatcher ----

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
		return toolGetVersion()
	case "list_boards":
		return toolListBoards(c)
	case "list_issues":
		return toolListIssues(c, argsRaw)
	case "get_issue":
		return toolGetIssue(c, argsRaw)
	case "create_issue":
		return toolCreateIssue(c, argsRaw)
	case "approve_issue":
		return toolApproveIssue(c, argsRaw)
	case "claim_issue":
		return toolClaimIssue(c, argsRaw)
	case "submit_qa":
		return toolSubmitQA(c, argsRaw)
	case "complete_issue":
		return toolCompleteIssue(c, argsRaw)
	case "reject_issue":
		return toolRejectIssue(c, argsRaw)
	case "add_comment":
		return toolAddComment(c, argsRaw)
	case "update_criteria":
		return toolUpdateCriteria(c, argsRaw)
	case "search_issues":
		return toolSearchIssues(c, argsRaw)
	}
	return nil, unknownToolError(name)
}
