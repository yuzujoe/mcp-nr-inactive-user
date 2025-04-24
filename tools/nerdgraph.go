package tools

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/yuzujoe/mcp-nr-inactive-user/pkg/nerdgraph"
)

func RegisterNerdGraphTool(mcpServer *server.MCPServer, nrClient *newrelic.NewRelic) error {
	if err := nerdgraph.GetInactiveUsers(mcpServer, nrClient); err != nil {
		return fmt.Errorf("failed to register nerdgraph tool: %w", err)
	}
	return nil
}

func NewNerdGraphServer(licenseKey string) (*newrelic.NewRelic, error) {
	nrClient, err := newrelic.New(newrelic.ConfigPersonalAPIKey(licenseKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create newrelic client: %w", err)
	}

	return nrClient, nil
}
