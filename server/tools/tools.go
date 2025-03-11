package tools

import (
	bravesearch "github.com/cnosuke/go-brave-search"
	mcp "github.com/metoro-io/mcp-golang"
)

type ResultList struct {
	Results []WebResult `json:"results"`
}

type WebResult struct {
	Title          string `json:"title"`
	URL            string `json:"url"`
	Description    string `json:"description"`
	PageAge        string `json:"page_age"`
	Language       string `json:"language"`
	FamilyFriendly bool   `json:"family_friendly"`
}

// SearchExecutor defines the interface for search execution
type SearchExecutor interface {
	ExecuteSearch(query string, params *bravesearch.WebSearchParams) (*ResultList, error)
}

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(mcpServer *mcp.Server, searchExecutor SearchExecutor) error {
	// Register web_search tool
	if err := RegisterWebSearchTool(mcpServer, searchExecutor); err != nil {
		return err
	}

	return nil
}
