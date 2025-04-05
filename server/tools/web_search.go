package tools

import (
	"context"
	"encoding/json"

	bravesearch "github.com/cnosuke/go-brave-search"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// WebSearchArgs - Arguments for web_search tool (保持しておく旧APIの型)
type WebSearchArgs struct {
	Query      string `json:"query" jsonschema:"description=Search query"`
	Count      int    `json:"count,omitempty" jsonschema:"description=Number of results (default: 10, max: 20)"`
	Offset     int    `json:"offset,omitempty" jsonschema:"description=Pagination offset (default: 0)"`
	SafeSearch string `json:"safe_search,omitempty" jsonschema:"description=Safe search mode (off, moderate, strict)"`
	Freshness  string `json:"freshness,omitempty" jsonschema:"description=Freshness filter (pd: past day, pw: past week, pm: past month, py: past year)"`
	SpellCheck bool   `json:"spellcheck,omitempty" jsonschema:"description=Enable spellcheck"`
	Country    string `json:"country,omitempty" jsonschema:"description=Country code (e.g., US, JP)"`
	SearchLang string `json:"search_lang,omitempty" jsonschema:"description=Search language (e.g., en, jp)"`
	UILang     string `json:"ui_lang,omitempty" jsonschema:"description=UI language (e.g., en-US, ja-JP)"`
}

// RegisterWebSearchTool - Register the web_search tool
func RegisterWebSearchTool(mcpServer *server.MCPServer, searchExecutor SearchExecutor) error {
	zap.S().Debugw("registering web_search tool")

	// ツールの定義
	tool := mcp.NewTool("web_search",
		mcp.WithDescription("Performs a web search using the Brave Search API. Use this for broad information gathering, recent events, or when you need diverse web sources. Supports pagination, content filtering, and freshness controls. Maximum 20 results per request, with offset for pagination."),
		mcp.WithString("query",
			mcp.Description("Search query"),
			mcp.Required(),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of results (default: 10, max: 20)"),
		),
		mcp.WithNumber("offset",
			mcp.Description("Pagination offset (default: 0)"),
		),
		mcp.WithString("safe_search",
			mcp.Description("Safe search mode (off, moderate, strict)"),
		),
		mcp.WithString("freshness",
			mcp.Description("Freshness filter (pd: past day, pw: past week, pm: past month, py: past year)"),
		),
		mcp.WithBoolean("spellcheck",
			mcp.Description("Enable spellcheck"),
		),
		mcp.WithString("country",
			mcp.Description("Country code (e.g., US, JP)"),
		),
		mcp.WithString("search_lang",
			mcp.Description("Search language (e.g., en, jp)"),
		),
		mcp.WithString("ui_lang",
			mcp.Description("UI language (e.g., en-US, ja-JP)"),
		),
	)

	// ツールハンドラーの登録
	mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 引数の取得
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}

		// その他のパラメータを取得
		count := 10
		if countVal, ok := request.Params.Arguments["count"].(float64); ok {
			count = int(countVal)
		}

		offset := 0
		if offsetVal, ok := request.Params.Arguments["offset"].(float64); ok {
			offset = int(offsetVal)
		}

		safeSearch := ""
		if safeSearchVal, ok := request.Params.Arguments["safe_search"].(string); ok {
			safeSearch = safeSearchVal
		}

		freshness := ""
		if freshnessVal, ok := request.Params.Arguments["freshness"].(string); ok {
			freshness = freshnessVal
		}

		spellcheck := false
		if spellcheckVal, ok := request.Params.Arguments["spellcheck"].(bool); ok {
			spellcheck = spellcheckVal
		}

		country := ""
		if countryVal, ok := request.Params.Arguments["country"].(string); ok {
			country = countryVal
		}

		searchLang := ""
		if searchLangVal, ok := request.Params.Arguments["search_lang"].(string); ok {
			searchLang = searchLangVal
		}

		uiLang := ""
		if uiLangVal, ok := request.Params.Arguments["ui_lang"].(string); ok {
			uiLang = uiLangVal
		}

		zap.S().Debugw("executing web_search",
			"query", query,
			"count", count,
			"offset", offset,
			"safe_search", safeSearch,
			"freshness", freshness,
			"spellcheck", spellcheck,
			"country", country,
			"search_lang", searchLang,
			"ui_lang", uiLang)

		// Set up search parameters
		params := bravesearch.NewWebSearchParams()

		// Apply user-provided parameters
		if count > 0 {
			if count > 20 {
				params.Count = 20 // Max 20 results per Brave Search API
			} else {
				params.Count = count
			}
		} else {
			params.Count = 10 // Default count
		}

		if offset >= 0 {
			params.Offset = offset
		}

		// Apply safe search settings
		if safeSearch != "" {
			switch safeSearch {
			case "off":
				params.SafeSearch = bravesearch.SafeSearchOff
			case "moderate":
				params.SafeSearch = bravesearch.SafeSearchModerate
			case "strict":
				params.SafeSearch = bravesearch.SafeSearchStrict
			default:
				// Use default (moderate) if invalid
				params.SafeSearch = bravesearch.SafeSearchModerate
			}
		}

		// Apply freshness filter
		if freshness != "" {
			switch freshness {
			case "pd":
				params.Freshness = bravesearch.FreshnessDay
			case "pw":
				params.Freshness = bravesearch.FreshnessWeek
			case "pm":
				params.Freshness = bravesearch.FreshnessMonth
			case "py":
				params.Freshness = bravesearch.FreshnessYear
			default:
				// No default freshness
			}
		}

		// Apply spellcheck
		params.Spellcheck = spellcheck

		// Apply country, search language, and UI language
		if country != "" {
			params.Country = country
		}

		if searchLang != "" {
			params.SearchLang = searchLang
		}

		if uiLang != "" {
			params.UILang = uiLang
		}

		// Execute search
		results, err := searchExecutor.ExecuteSearch(query, params)
		if err != nil {
			zap.S().Errorw("failed to execute search",
				"query", query,
				"params", params,
				"error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert search results to JSON
		jsonContent, err := json.Marshal(results)
		if err != nil {
			return mcp.NewToolResultError("failed to marshal search results to JSON"), nil
		}

		// Return the search results as text content
		return mcp.NewToolResultText(string(jsonContent)), nil
	})

	return nil
}
