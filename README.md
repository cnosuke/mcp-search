# MCP Search Server

MCP Search Server is a Go-based MCP server implementation that provides web search functionality using the Brave Search API, allowing MCP clients (e.g., Claude Desktop) to perform web searches.

## Features

* MCP Compliance: Provides a JSON‐RPC based interface for tool execution according to the MCP specification.
* Web Search: Supports web search using Brave Search API with various parameters for customization.

## Requirements

* Go 1.24 or later
* Brave Search API key

## Configuration

The server is configured via a YAML file (default: config.yml). For example:

```yaml
log: 'path/to/mcp-search.log' # Log file path, if empty no log will be produced
debug: false # Enable debug mode for verbose logging

search:
  api_key: "your-brave-search-api-key"
  timeout: 30
  max_retries: 2
```

You can override configurations using environment variables:
- `LOG_PATH`: Path to log file
- `DEBUG`: Enable debug mode (true/false)
- `BRAVE_SEARCH_API_KEY`: Your Brave Search API key
- `BRAVE_SEARCH_TIMEOUT`: Request timeout in seconds
- `BRAVE_SEARCH_MAX_RETRIES`: Number of retries for failed requests
- `BRAVE_SEARCH_COUNTRY`: Default country code for search
- `BRAVE_SEARCH_LANGUAGE`: Default search language
- `BRAVE_SEARCH_UI_LANGUAGE`: Default UI language

## Logging

Logging behavior is controlled through configuration:

- If `log` is set in the config file, logs will be written to the specified file
- If `log` is empty, no logs will be produced
- Set `debug: true` for more verbose logging

## MCP Server Usage

MCP clients interact with the server by sending JSON‐RPC requests to execute various tools. The following MCP tools are supported:

* `web_search`: Performs a web search using the Brave Search API with configurable parameters.

### Tool Parameters

The `web_search` tool accepts the following parameters:

* `query` (required): The search query string.
* `count` (optional): Number of results to return (default: 10, max: 20).
* `offset` (optional): Pagination offset for results (default: 0).
* `safe_search` (optional): Safe search mode ("off", "moderate", "strict").
* `freshness` (optional): Freshness filter ("pd" for past day, "pw" for past week, "pm" for past month, "py" for past year).
* `spellcheck` (optional): Enable or disable spellcheck.
* `country` (optional): Country code for search results (e.g., "US", "JP").
* `search_lang` (optional): Search language (e.g., "en", "jp").
* `ui_lang` (optional): UI language (e.g., "en-US", "ja-JP").

### Using with Claude Desktop

To integrate with Claude Desktop, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "search": {
      "command": "./bin/mcp-search",
      "args": ["server"],
      "env": {
        "LOG_PATH": "mcp-search.log",
        "DEBUG": "false",
        "BRAVE_SEARCH_API_KEY": "your-brave-search-api-key",
        "BRAVE_SEARCH_COUNTRY": "US",
        "BRAVE_SEARCH_LANGUAGE": "en",
        "BRAVE_SEARCH_UI_LANGUAGE": "en-US"
      }
    }
  }
}
```

This configuration registers the MCP Search Server with Claude Desktop, ensuring that all logs are directed to the specified log file.

## Command-Line Parameters

When starting the server, you can specify various settings:

```
./bin/mcp-search server [options]
```

Options:
* `--config`, `-c`: Path to the configuration file (default: "config.yml").

## Building and Running

```bash
# Build the server
make build

# Run the server
make run

# or run with specific options
./bin/mcp-search server --config=config.yml
```

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for improvements or bug fixes. For major changes, open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License.

Author: cnosuke ( x.com/cnosuke )
