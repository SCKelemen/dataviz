package mcp

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the DataViz MCP server
type Server struct {
	server *mcp.Server
}

// NewServer creates a new DataViz MCP server
func NewServer() (*Server, error) {
	// Create MCP server with name and version
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "dataviz-mcp",
			Version: "0.1.0",
		},
		nil, // ServerOptions
	)

	s := &Server{
		server: mcpServer,
	}

	// Register tools
	s.RegisterTools()

	return s, nil
}

// Run starts the MCP server using stdio transport
func (s *Server) Run(ctx context.Context) error {
	log.Println("Starting DataViz MCP server...")
	log.Println("Listening on stdio...")

	// Run server with stdio transport
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

// GetMCPServer returns the underlying MCP server instance
func (s *Server) GetMCPServer() *mcp.Server {
	return s.server
}
