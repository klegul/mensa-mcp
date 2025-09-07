package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const akkBaseUrl = "https://mensa.akk.org/json"

var akkMensaApi = NewAkkMensaApi(akkBaseUrl)

func main() {
	s := server.NewMCPServer("Mensa Karlsruhe", "0.1.1",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, true),
		server.WithLogging(),
	)

	s.AddTool(
		mcp.NewTool("get_available_dates",
			mcp.WithDescription("Get available dates"),
		),
		handleMenuTool,
	)

	s.AddResource(
		mcp.NewResource(
			"mensa-ka://menu/{date}",
			"Mensa Menu",
			mcp.WithResourceDescription("Mensa menu for a specific date"),
			mcp.WithMIMEType("application/json"),
		),
		handleMenuResource,
	)

	log.Println("Starting StreamableHTTP server on :8080")
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/api/v1/mcp"),
		server.WithHeartbeatInterval(30*time.Second))
	if err := httpServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func handleMenuResource(_ context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	dateString := req.Params.URI[:16]
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, fmt.Errorf("invalid date. Use YYYY-MM-DD format")
	}

	result, err := akkMensaApi.GetMenuForDate(date)
	if err != nil {
		return nil, fmt.Errorf("error fetching menu: %v", err)
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("error marshalling menu to JSON: %v", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     string(jsonData),
		},
	}, nil

}

func handleMenuTool(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dates, err := akkMensaApi.GetAvailableDates()
	if err != nil {
		return nil, fmt.Errorf("error fetching available dates: %v", err)
	}

	var dateStrings []string
	for _, date := range dates {
		dateStrings = append(dateStrings, date.Format("2006-01-02"))
	}

	resultObject := map[string]interface{}{
		"available_dates": dateStrings,
	}

	return mcp.NewToolResultStructured(resultObject, "Available dates: "+strings.Join(dateStrings, ", ")), nil
}
