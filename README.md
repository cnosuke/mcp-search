# MCP Search Server

MCP Search Server is a Go-based MCP server implementation that provides web search functionality using the Brave Search API, allowing MCP clients (e.g., Claude Desktop) to perform web searches.

## Features

- MCP Compliance: Provides a JSON‐RPC based interface for tool execution according to the MCP specification.
- Web Search: Supports web search using Brave Search API with various parameters for customization.

## Requirements

- Docker (recommended)
- Brave Search API key

For local development:

- Go 1.24 or later

## Using with Docker (Recommended)

```bash
docker pull cnosuke/mcp-search:latest

docker run \
  -i \
  -v /path/to/your/config.yml:/app/config.yml \
  -e BRAVE_SEARCH_API_KEY=your-brave-search-api-key \
  cnosuke/mcp-search:latest
```

### Using with Claude Desktop (Docker)

To integrate with Claude Desktop using Docker, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "search": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "BRAVE_SEARCH_API_KEY=your-brave-search-api-key",
        "-e",
        "BRAVE_SEARCH_COUNTRY=US",
        "-e",
        "BRAVE_SEARCH_LANGUAGE=en",
        "-e",
        "BRAVE_SEARCH_UI_LANGUAGE=en-US",
        "cnosuke/mcp-search:latest"
      ]
    }
  }
}
```

## Building and Running (Go Binary)

Alternatively, you can build and run the Go binary directly:

```bash
# Build the server
make bin/mcp-search

# Run the server
./bin/mcp-search server --config=config.yml
```

### Using with Claude Desktop (Go Binary)

To integrate with Claude Desktop using the Go binary, add an entry to your `claude_desktop_config.json` file:

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

## Configuration

The server is configured via a YAML file (default: config.yml). For example:

```yaml
log: 'path/to/mcp-search.log' # Log file path, if empty no log will be produced
debug: false # Enable debug mode for verbose logging

search:
  api_key: 'your-brave-search-api-key'
  timeout: 30
  max_retries: 2
  default_country: JP
  default_search_lang: jp
  default_ui_lang: ja-JP
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

- `web_search`: Performs a web search using the Brave Search API with configurable parameters.

### Tool Parameters

The `web_search` tool accepts the following parameters:

- `query` (required): The search query string.
- `count` (optional): Number of results to return (default: 10, max: 20).
- `offset` (optional): Pagination offset for results (default: 0).
- `safe_search` (optional): Safe search mode ("off", "moderate", "strict").
- `freshness` (optional): Freshness filter ("pd" for past day, "pw" for past week, "pm" for past month, "py" for past year).
- `spellcheck` (optional): Enable or disable spellcheck.
- `country` (optional): Country code for search results (e.g., "US", "JP").
- `search_lang` (optional): Search language (e.g., "en", "jp").
- `ui_lang` (optional): UI language (e.g., "en-US", "ja-JP").

## Command-Line Parameters

When starting the server, you can specify various settings:

```
./bin/mcp-search server [options]
```

Options:

- `--config`, `-c`: Path to the configuration file (default: "config.yml").

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for improvements or bug fixes. For major changes, open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License.

Author: cnosuke ( x.com/cnosuke )
