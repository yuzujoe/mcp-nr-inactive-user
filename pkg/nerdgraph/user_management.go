package nerdgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdgraph"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type InactiveUser struct {
	Email      string          `json:"email"`
	Name       string          `json:"name"`
	LastActive nrtime.DateTime `json:"last_active"`
}

func GetInactiveUsers(mcpServer *server.MCPServer, client *newrelic.NewRelic) error {
	tool := mcp.NewTool(
		"get_inactive_users",
		mcp.WithDescription("Returns the number of FullPlatform Users who have not logged in this month."),
		mcp.WithString("authentication_domain_id",
			mcp.Required(),
			mcp.Description("the newrelic authentication domain id"),
		),
	)

	mcpServer.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		authenticationDomainID := req.Params.Arguments["authentication_domain_id"].(string)
		if authenticationDomainID == "" {
			return mcp.NewToolResultError("authentication_domain_id is required"), nil
		}

		authenticationDomainIDs := []string{authenticationDomainID}

		vars := map[string]interface{}{
			"authenticationDomainIDs": authenticationDomainIDs,
		}

		res, err := client.NerdGraph.QueryWithContext(ctx, getUsersQuery, vars)
		if err != nil {
			return mcp.NewToolResultError("failed to get users: " + err.Error()), nil
		}

		respData, ok := res.(nerdgraph.QueryResponse)
		if !ok {
			return mcp.NewToolResultError("failed to parse response"), nil
		}

		// Navigate the response structure to get to users array
		actor, ok := respData.Actor.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("actor not found in response"), nil
		}

		organization, ok := actor["organization"].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("organization not found in response"), nil
		}

		userManagement, ok := organization["userManagement"].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("userManagement not found in response"), nil
		}

		authDomains, ok := userManagement["authenticationDomains"].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("authenticationDomains not found in response"), nil
		}

		domains, ok := authDomains["authenticationDomains"].([]interface{})
		if !ok || len(domains) == 0 {
			return mcp.NewToolResultError("no authentication domains found"), nil
		}

		// Process the first domain since we're effectively working with just one
		domainMap, ok := domains[0].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid authentication domain format"), nil
		}

		usersObj, ok := domainMap["users"].(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("users data not found in domain"), nil
		}

		usersArray, ok := usersObj["users"].([]interface{})
		if !ok {
			return mcp.NewToolResultError("invalid users array format"), nil
		}

		// Process all users in the domain
		allUsers := make([]map[string]interface{}, 0, len(usersArray))
		for _, user := range usersArray {
			userMap, ok := user.(map[string]interface{})
			if ok {
				// Add active this month information
				fmt.Println(userMap)
				lastActive, ok := userMap["lastActive"].(string)
				if ok {
					if !isActiveThisMonth(lastActive) {
						allUsers = append(allUsers, userMap)
					} else {
						fmt.Printf("User %v is active this month\n", userMap)
					}
				}
			}
		}

		if len(allUsers) == 0 {
			return mcp.NewToolResultText("No inactive users found"), nil
		}

		result, err := json.Marshal(allUsers)
		if err != nil {
			return mcp.NewToolResultError("failed to marshal result: " + err.Error()), nil
		}

		fmt.Printf("Inactive users: %s\n", result)

		return mcp.NewToolResultText(string(result)), nil
	})

	return nil
}

// isActiveThisMonth checks if the lastActive timestamp is in the current month
func isActiveThisMonth(lastActive string) bool {
	// Parse the lastActive timestamp
	lastActiveTime, err := time.Parse(time.RFC3339, lastActive)
	if err != nil {
		fmt.Printf("Error parsing lastActive timestamp '%s': %v\n", lastActive, err)
		return false
	}

	// Get current time
	now := time.Now()

	// Compare year and month
	return lastActiveTime.Year() == now.Year() && lastActiveTime.Month() == now.Month()
}

const getUsersQuery = `query(
	  $authenticationDomainIDs: [ID!],
	) 
	{
	  actor {
		user {
		  name
		}
		organization {
		  userManagement {
			authenticationDomains(id: $authenticationDomainIDs) {
			  authenticationDomains {
				users(filter: {type: {eq: FULL_PLATFORM}}) {
				  users {
					lastActive
					email
					name
				  }
				  totalCount
				  nextCursor
				}
			  }
			}
		  }
		}
	  }
	}`
