package tools

import (
	bravesearch "github.com/cnosuke/go-brave-search"
	mcp "github.com/metoro-io/mcp-golang"
)

// SearchExecutor defines the interface for search execution
type SearchExecutor interface {
	ExecuteSearch(query string, params *bravesearch.WebSearchParams) (*bravesearch.WebSearchResponse, error)
}

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(mcpServer *mcp.Server, searchExecutor SearchExecutor) error {
	// Register web_search tool
	if err := RegisterWebSearchTool(mcpServer, searchExecutor); err != nil {
		return err
	}

	return nil
}
