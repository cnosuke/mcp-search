package tools

import (
	"encoding/json"

	bravesearch "github.com/cnosuke/go-brave-search"
	"github.com/cockroachdb/errors"
	mcp "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"
)

// WebSearchArgs - Arguments for web_search tool
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
func RegisterWebSearchTool(server *mcp.Server, searchExecutor SearchExecutor) error {
	zap.S().Debugw("registering web_search tool")
	err := server.RegisterTool("web_search", "Performs a web search using the Brave Search API. Use this for broad information gathering, recent events, or when you need diverse web sources. Supports pagination, content filtering, and freshness controls. Maximum 20 results per request, with offset for pagination.",
		func(args WebSearchArgs) (*mcp.ToolResponse, error) {
			zap.S().Debugw("executing web_search",
				"query", args.Query,
				"count", args.Count,
				"offset", args.Offset,
				"safe_search", args.SafeSearch,
				"freshness", args.Freshness,
				"spellcheck", args.SpellCheck,
				"country", args.Country,
				"search_lang", args.SearchLang,
				"ui_lang", args.UILang)

			// Validate query
			if args.Query == "" {
				return nil, errors.New("query is required")
			}

			// Set up search parameters
			params := bravesearch.NewWebSearchParams()

			// Apply user-provided parameters
			if args.Count > 0 {
				if args.Count > 20 {
					params.Count = 20 // Max 20 results per Brave Search API
				} else {
					params.Count = args.Count
				}
			} else {
				params.Count = 10 // Default count
			}

			if args.Offset >= 0 {
				params.Offset = args.Offset
			}

			// Apply safe search settings
			if args.SafeSearch != "" {
				switch args.SafeSearch {
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
			if args.Freshness != "" {
				switch args.Freshness {
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
			params.Spellcheck = args.SpellCheck

			// Apply country, search language, and UI language
			if args.Country != "" {
				params.Country = args.Country
			}

			if args.SearchLang != "" {
				params.SearchLang = args.SearchLang
			}

			if args.UILang != "" {
				params.UILang = args.UILang
			}

			// Execute search
			results, err := searchExecutor.ExecuteSearch(args.Query, params)
			if err != nil {
				zap.S().Errorw("failed to execute search",
					"query", args.Query,
					"params", params,
					"error", err)
				return nil, errors.Wrap(err, "failed to execute search")
			}

			// Convert search results to JSON
			jsonContent, err := json.Marshal(results)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal search results to JSON")
			}

			// Return the search results as text content
			return mcp.NewToolResponse(mcp.NewTextContent(string(jsonContent))), nil
		})

	if err != nil {
		zap.S().Errorw("failed to register web_search tool", "error", err)
		return errors.Wrap(err, "failed to register web_search tool")
	}

	return nil
}
