package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/austiecodes/gomor/internal/client"
	"github.com/austiecodes/gomor/internal/memory/retrieval"
	"github.com/austiecodes/gomor/internal/memory/store"
	"github.com/austiecodes/gomor/internal/provider"
	"github.com/austiecodes/gomor/internal/types"
	"github.com/austiecodes/gomor/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MemoryRetrieveInput defines the input schema for the memory retrieve tool
type MemoryRetrieveInput struct {
	Query string `json:"query" jsonschema:"the query to search for related memories"`
}

// MemoryRetrieveOutput defines the output schema for the memory retrieve tool
type MemoryRetrieveOutput struct {
	Results string `json:"results" jsonschema:"formatted text containing retrieved memories"`
}

// handleMemoryRetrieve handles the goa_memory_retrieve tool call (unified hybrid search)
func handleMemoryRetrieve(ctx context.Context, request *mcp.CallToolRequest, input MemoryRetrieveInput) (*mcp.CallToolResult, MemoryRetrieveOutput, error) {
	// Validate query (required)
	query := strings.TrimSpace(input.Query)
	if query == "" {
		return nil, MemoryRetrieveOutput{}, fmt.Errorf("parameter 'query' must be a non-empty string")
	}

	// Load config
	config, err := utils.LoadConfig()
	if err != nil {
		return nil, MemoryRetrieveOutput{}, fmt.Errorf("failed to load config: %w", err)
	}

	if config.Model.EmbeddingModel == nil {
		return nil, MemoryRetrieveOutput{}, fmt.Errorf("embedding model not configured. Run 'gomor set' to configure")
	}

	// Open memory store
	memStore, err := store.NewStore()
	if err != nil {
		return nil, MemoryRetrieveOutput{}, fmt.Errorf("failed to open memory store: %w", err)
	}
	defer memStore.Close()

	// Create embedding client
	embeddingModel := *config.Model.EmbeddingModel
	embClient, err := provider.NewEmbeddingClient(config, embeddingModel.Provider)
	if err != nil {
		return nil, MemoryRetrieveOutput{}, fmt.Errorf("failed to create embedding client: %w", err)
	}

	// Create query client for LLM transformations (optional, may be nil)
	var queryClient client.QueryClient
	toolModel := types.Model{}
	if config.Model.ToolModel != nil {
		toolModel = *config.Model.ToolModel
		queryClient, _ = provider.NewQueryClient(config, toolModel.Provider)
	}

	// Create retriever
	ret := retrieval.NewRetriever(
		memStore,
		embClient,
		queryClient,
		embeddingModel,
		toolModel,
		config.Memory,
	)

	// Perform retrieval
	response, err := ret.Retrieve(ctx, query)
	if err != nil {
		return nil, MemoryRetrieveOutput{}, fmt.Errorf("retrieval failed: %w", err)
	}

	// Format results
	result := retrieval.FormatAsText(response)
	return nil, MemoryRetrieveOutput{
		Results: result,
	}, nil
}
