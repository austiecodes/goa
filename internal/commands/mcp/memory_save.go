package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/austiecodes/gomor/internal/memory/store"
	"github.com/austiecodes/gomor/internal/provider"
	"github.com/austiecodes/gomor/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MemorySaveInput defines the input schema for the memory save tool
type MemorySaveInput struct {
	Text string `json:"text" jsonschema:"the preference or fact to save"`
	Tags string `json:"tags,omitempty" jsonschema:"comma-separated tags for categorization"`
}

// MemorySaveOutput defines the output schema for the memory save tool
type MemorySaveOutput struct {
	Message string `json:"message" jsonschema:"success message with memory ID"`
	ID      string `json:"id" jsonschema:"the ID of the saved memory"`
}

// handleMemorySave handles the memory_save tool call
func handleMemorySave(ctx context.Context, request *mcp.CallToolRequest, input MemorySaveInput) (*mcp.CallToolResult, MemorySaveOutput, error) {
	// Validate text (required)
	text := strings.TrimSpace(input.Text)
	if text == "" {
		return nil, MemorySaveOutput{}, fmt.Errorf("parameter 'text' must be a non-empty string")
	}

	// Extract tags (optional)
	var tags []string
	if input.Tags != "" {
		for _, t := range strings.Split(input.Tags, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	// Load config for embedding
	config, err := utils.LoadConfig()
	if err != nil {
		return nil, MemorySaveOutput{}, fmt.Errorf("failed to load config: %w", err)
	}

	if config.Model.EmbeddingModel == nil {
		return nil, MemorySaveOutput{}, fmt.Errorf("embedding model not configured. Run 'gomor set' to configure")
	}

	// Create embedding client
	embeddingModel := *config.Model.EmbeddingModel
	embClient, err := provider.NewEmbeddingClient(config, embeddingModel.Provider)
	if err != nil {
		return nil, MemorySaveOutput{}, fmt.Errorf("failed to create embedding client: %w", err)
	}

	// Generate embedding
	embedding, err := embClient.Embed(ctx, embeddingModel, text)
	if err != nil {
		return nil, MemorySaveOutput{}, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Normalize embedding for cosine similarity
	normalizedEmbedding := store.NormalizeVector(embedding)

	// Open memory store
	memStore, err := store.NewStore()
	if err != nil {
		return nil, MemorySaveOutput{}, fmt.Errorf("failed to open memory store: %w", err)
	}
	defer memStore.Close()

	// Save memory
	item := &store.MemoryItem{
		Text:      text,
		Tags:      tags,
		Source:    store.SourceExplicit,
		Provider:  embeddingModel.Provider,
		ModelID:   embeddingModel.ModelID,
		Dim:       len(normalizedEmbedding),
		Embedding: normalizedEmbedding,
	}

	if err := memStore.SaveMemory(item); err != nil {
		return nil, MemorySaveOutput{}, fmt.Errorf("failed to save memory: %w", err)
	}

	return nil, MemorySaveOutput{
		Message: fmt.Sprintf("Memory saved successfully (id: %s)", item.ID),
		ID:      item.ID,
	}, nil
}
