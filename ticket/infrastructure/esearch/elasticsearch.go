package esearch

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

func NewElasticsearchClient(address string) (*ElasticsearchClient, error) {
	client, err := elasticsearch.New(
		elasticsearch.WithAddresses(address),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente Elasticsearch: %w", err)
	}

	return &ElasticsearchClient{
		client: client,
	}, nil
}

func (ec *ElasticsearchClient) CreateIndex(ctx context.Context, index string) error {
	response, error := ec.client.Indices.Create(index, ec.client.Indices.Create.WithContext(ctx))

	if error != nil {
		return fmt.Errorf("erro ao criar índice no Elasticsearch: %w", error)
	}
	defer response.Body.Close()

	if response.IsError() {
		return fmt.Errorf("erro ao criar índice no Elasticsearch: %s", response.String())
	}

	return nil
}

func (ec *ElasticsearchClient) IndexDocument(ctx context.Context, index, documentID string, document []byte) error {
	response, error := ec.client.Index(
		index,
		bytes.NewReader(document),
		ec.client.Index.WithDocumentID(documentID),
		ec.client.Index.WithContext(ctx),
		ec.client.Index.WithRefresh("true"),
	)
	if error != nil {
		return fmt.Errorf("erro ao indexar documento no Elasticsearch: %w", error)
	}
	defer response.Body.Close()

	if response.IsError() {
		return fmt.Errorf("erro ao indexar documento no Elasticsearch: %s", response.String())
	}
	return nil
}

func (ec *ElasticsearchClient) Search(ctx context.Context, index, keyword string) ([]byte, error) {

	query := fmt.Sprintf(`{"query":{"multi_match":{"query":%q,"fields":["Name","Description"]}}}`, keyword)

	response, error := ec.client.Search(
		ec.client.Search.WithContext(ctx),
		ec.client.Search.WithIndex(index),
		ec.client.Search.WithBody(bytes.NewReader([]byte(query))),
	)

	if error != nil {
		return nil, fmt.Errorf("erro ao pesquisar no Elasticsearch: %w", error)
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf("erro ao pesquisar no Elasticsearch: %s", response.String())
	}

	body, error := io.ReadAll(response.Body)
	if error != nil {
		return nil, fmt.Errorf("erro ao ler resposta do Elasticsearch: %w", error)
	}

	return body, nil
}
