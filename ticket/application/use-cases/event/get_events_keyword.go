package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Indexer interface {
	Search(ctx context.Context, index, keyword string) ([]byte, error)
}

type EvenKeyWordtOutputModel struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Date        time.Time `json:"date"`
}

type GetEventsKeyword struct {
	elasticSearchClient Indexer
}

func NewGetEventsKeyword(elasticSearch Indexer) *GetEventsKeyword {
	return &GetEventsKeyword{
		elasticSearchClient: elasticSearch,
	}
}

func (g *GetEventsKeyword) Execute(ctx context.Context, index, keyword string) ([]EvenKeyWordtOutputModel, error) {
	body, err := g.elasticSearchClient.Search(ctx, index, keyword)
	if err != nil {
		return nil, fmt.Errorf("erro ao pesquisar no Elasticsearch: %w", err)
	}

	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source EvenKeyWordtOutputModel `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(body, &esResponse); err != nil {
		return nil, fmt.Errorf("erro ao unmarshal resultado da pesquisa: %w", err)
	}

	results := make([]EvenKeyWordtOutputModel, 0, len(esResponse.Hits.Hits))
	for _, hit := range esResponse.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}
