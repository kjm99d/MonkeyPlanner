package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// toolListBoards handles the list_boards MCP tool.
func toolListBoards(c *mcpHTTPClient) (any, error) {
	raw, err := c.get("/api/boards", nil)
	if err != nil {
		return nil, err
	}
	return unmarshalAny(raw), nil
}

// toolListIssues handles the list_issues MCP tool.
func toolListIssues(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		BoardID string `json:"boardId"`
		Status  string `json:"status"`
	}](argsRaw)
	if err != nil {
		return nil, err
	}
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
	return unmarshalAny(raw), nil
}

// toolGetIssue handles the get_issue MCP tool, returning issue + children + comments.
func toolGetIssue(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
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
	var issueData map[string]any
	var comments any
	_ = json.Unmarshal(issueRaw, &issueData)
	_ = json.Unmarshal(commentsRaw, &comments)
	issueData["comments"] = comments
	return issueData, nil
}

// toolCreateIssue handles the create_issue MCP tool. When instructions are supplied,
// a follow-up PATCH sets them since POST /api/issues does not accept the field.
func toolCreateIssue(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		BoardID      string `json:"boardId"`
		Title        string `json:"title"`
		Body         string `json:"body"`
		Instructions string `json:"instructions"`
	}](argsRaw)
	if err != nil {
		return nil, err
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
	return unmarshalAny(issueRaw), nil
}

// toolApproveIssue handles the approve_issue MCP tool.
func toolApproveIssue(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
		return nil, fmt.Errorf("issueId is required")
	}
	raw, err := c.post("/api/issues/"+args.IssueID+"/approve", map[string]any{})
	if err != nil {
		return nil, err
	}
	return unmarshalAny(raw), nil
}

// toolClaimIssue handles the claim_issue MCP tool.
func toolClaimIssue(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
		return nil, fmt.Errorf("issueId is required")
	}
	raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "InProgress"})
	if err != nil {
		return nil, err
	}
	return unmarshalAny(raw), nil
}

// toolSubmitQA handles the submit_qa MCP tool.
func toolSubmitQA(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
		Comment string `json:"comment"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
		return nil, fmt.Errorf("issueId is required")
	}
	raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "QA"})
	if err != nil {
		return nil, err
	}
	if args.Comment != "" {
		_, _ = c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": args.Comment})
	}
	return unmarshalAny(raw), nil
}

// toolCompleteIssue handles the complete_issue MCP tool.
func toolCompleteIssue(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
		Comment string `json:"comment"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
		return nil, fmt.Errorf("issueId is required")
	}
	raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "Done"})
	if err != nil {
		return nil, err
	}
	if args.Comment != "" {
		_, _ = c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": args.Comment})
	}
	return unmarshalAny(raw), nil
}

// toolRejectIssue handles the reject_issue MCP tool.
func toolRejectIssue(c *mcpHTTPClient, argsRaw json.RawMessage) (any, error) {
	args, err := decodeArgs[struct {
		IssueID string `json:"issueId"`
		Reason  string `json:"reason"`
	}](argsRaw)
	if err != nil || args.IssueID == "" {
		return nil, fmt.Errorf("issueId is required")
	}
	if args.Reason == "" {
		return nil, fmt.Errorf("reason is required for rejection")
	}
	raw, err := c.patch("/api/issues/"+args.IssueID, map[string]any{"status": "InProgress"})
	if err != nil {
		return nil, err
	}
	_, _ = c.post("/api/issues/"+args.IssueID+"/comments", map[string]any{"body": "❌ QA rejected: " + args.Reason})
	return unmarshalAny(raw), nil
}
