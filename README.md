# MCP NR Inactive User

## Overview
A tool to manage and identify inactive fullplatform users in New Relic account.

## Features
- Returns a list of inactive users who have not logged in to New Relic this month

## Requirements
- Go 1.24 or later
- Docker (for containerized deployment)
- Make (optional, for using the provided Makefile)

## Using with Docker(Recommended)

```shell
make docker-build
```

### Using for Claude Desktop

```json
{
  "mcpServers": {
    "nr-inactive-users": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "NEW_RELIC_API_KEY",
        "mcp-nr-inactive-user"
      ],
      "env": {
        "NEW_RELIC_API_KEY": "<set newrelic user key>"
      }
    }
  }
}
```

## Using with Go binary

```shell
make build
```

### Using for Claude Desktop

```json
{
  "mcpServers": {
    "nr-inactive-users": {
      "command": "./dist/mcp-server",
      "args": ["studio"],
      "env": {
        "NEW_RELIC_API_KEY": "<set newrelic user key>"
      }
    }
  }
}
```

## Usage

## Contributing
Welcome to contributions！ Please create an issue first, then we look forward to your pull request.

## License
This project is licensed under the terms of the MIT open source license. Please refer to [MIT](./LICENSE) for the full terms.
