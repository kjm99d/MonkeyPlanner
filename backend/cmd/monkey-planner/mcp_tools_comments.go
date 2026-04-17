package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// toolGetVersion handles the get_version MCP tool.
func toolGetVersion() (any, error) {
	return map[string]any{"version": version, "name": "monkey-planner"}, nil
}

// toolAddComment handles the add_comment MCP tool.
func toolAddComment(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
		Body    string `json:"body"`
	}](argsRaw)
	if err != nil || args.IssueID == "" || args.Body == "" {
		return nil, fmt.Errorf("issueId and body are required")
	}
	raw, err := c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": args.Body})
	if err != nil {
		return nil, err
	}
	return unmarshalAny(raw), nil
}

// toolUpdateCriteria toggles a single criterion on an issue by index.
func toolUpdateCriteria(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
		Index   int    `json:"index"`
		Done    bool   `json:"done"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
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
}

// toolSearchIssues returns issues whose title contains the query (case-insensitive).
func toolSearchIssues(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		Query string `json:"query"`
	}](argsRaw)
	if err != nil || args.Query == "" {
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
}
