package tools

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
)

func ResisterTools(mcpServer *server.MCPServer, nrClinet *newrelic.NewRelic) error {
	if err := RegisterNerdGraphTool(mcpServer, nrClinet); err != nil {
		return fmt.Errorf("failed to register nerdgraph tool: %w", err)
	}
	return nil
}
