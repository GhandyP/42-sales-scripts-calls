package playbook

import (
	"errors"

	"github.com/GhandyP/42-sales-scripts-calls/internal/sales"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type searchPlaybookArgs struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type searchPlaybookResult struct {
	Passages []sales.EvidencePassage `json:"passages"`
}

// NewSearchTool creates a read-only ADK tool bound to one request's approved
// playbook index. It accepts search terms only; playbook content is not mutable
// through the tool interface.
func NewSearchTool(index *Index) (tool.Tool, error) {
	if index == nil {
		return nil, errors.New("playbook index is required")
	}

	return functiontool.New(functiontool.Config{
		Name: "search_playbook",
		Description: "Search only the approved playbook supplied for this request. " +
			"Return exact passages and stable citation IDs as evidence; do not invent or rewrite claims. " +
			"Limit is optional, defaults to three results, and cannot exceed five.",
	}, func(_ agent.Context, args searchPlaybookArgs) (searchPlaybookResult, error) {
		passages, err := index.Search(args.Query, args.Limit)
		if err != nil {
			return searchPlaybookResult{}, err
		}
		return searchPlaybookResult{Passages: passages}, nil
	})
}
