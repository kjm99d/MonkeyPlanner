package main

// mcpToolDefinitions returns the MCP tools/list schema declarations.
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
