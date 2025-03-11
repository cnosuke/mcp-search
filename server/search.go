package server

import (
	"context"

	bravesearch "github.com/cnosuke/go-brave-search"
	"github.com/cnosuke/mcp-search/config"
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

// SearchServer - Search server structure
type SearchServer struct {
	client             *bravesearch.Client
	cfg                *config.Config
	DefaultCountry     string
	DefaultSearchLang  string
	DefaultUILang      string
}

// NewSearchServer - Create a new Search server
func NewSearchServer(cfg *config.Config) (*SearchServer, error) {
	zap.S().Infow("creating new Search server",
		"default_country", cfg.Search.DefaultCountry,
		"default_search_lang", cfg.Search.DefaultSearchLang,
		"default_ui_lang", cfg.Search.DefaultUILang)

	// Create Brave Search client
	client, err := bravesearch.NewClient(
		cfg.Search.APIKey,
		bravesearch.WithTimeout(cfg.Search.Timeout),
		bravesearch.WithRetries(cfg.Search.MaxRetries),
		bravesearch.WithDefaultCountry(cfg.Search.DefaultCountry),
		bravesearch.WithDefaultSearchLanguage(cfg.Search.DefaultSearchLang),
		bravesearch.WithDefaultUILanguage(cfg.Search.DefaultUILang),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Brave Search client")
	}

	return &SearchServer{
		client:             client,
		cfg:                cfg,
		DefaultCountry:     cfg.Search.DefaultCountry,
		DefaultSearchLang:  cfg.Search.DefaultSearchLang,
		DefaultUILang:      cfg.Search.DefaultUILang,
	}, nil
}

// ExecuteSearch - Execute a web search
func (s *SearchServer) ExecuteSearch(query string, params *bravesearch.WebSearchParams) (*bravesearch.WebSearchResponse, error) {
	ctx := context.Background()

	// Prepare search parameters
	searchParams := params
	if searchParams == nil {
		searchParams = bravesearch.NewWebSearchParams()
	}

	// Set query
	searchParams.Query = query

	// Perform search
	zap.S().Infow("executing search",
		"query", query,
		"country", searchParams.Country,
		"search_lang", searchParams.SearchLang,
		"ui_lang", searchParams.UILang,
		"count", searchParams.Count,
		"offset", searchParams.Offset)

	results, err := s.client.WebSearch(ctx, query, searchParams)
	if err != nil {
		zap.S().Errorw("search failed", "error", err)
		return nil, errors.Wrap(err, "search failed")
	}

	zap.S().Infow("search completed",
		"query", query,
		"result_count", len(results.GetWebResults()))

	return results, nil
}
