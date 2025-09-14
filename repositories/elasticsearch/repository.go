package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"coin-feed/internal/domain/repository"
	"coin-feed/pkg/logger"

	"github.com/elastic/go-elasticsearch/v9/esutil"

	"coin-feed/pkg/elasticsearch"
)

type Repository struct {
	client *elasticsearch.Client
	index  string
}

func NewRepository(client *elasticsearch.Client, index string) *Repository {
	return &Repository{
		client: client,
		index:  index,
	}
}

func (r *Repository) SaveLatestCryptoCurrency(ctx context.Context, data []*repository.CryptoCurrencyData) error {
	if len(data) == 0 {
		return nil
	}

	bulk, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         r.index,
		Client:        r.client,
		NumWorkers:    4,
		FlushBytes:    5 << 20,
		FlushInterval: 2 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("bulk indexer init: %w", err)
	}

	type docWithTS struct {
		*repository.CryptoCurrencyData
		Timestamp time.Time `json:"@timestamp"`
	}

	for _, cryptoInfo := range data {
		doc := docWithTS{
			CryptoCurrencyData: cryptoInfo,
			Timestamp:          time.Now().UTC(),
		}

		body, mErr := json.Marshal(doc)
		if mErr != nil {
			return fmt.Errorf("marshal doc id=%v: %w", cryptoInfo.Id, mErr)
		}

		item := esutil.BulkIndexerItem{
			Action: "create",
			Body:   bytes.NewReader(body),

			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, e error) {
				if e != nil {
					logger.Logger.Error(fmt.Sprintf("[bulk] transport error: %v", e))
					return
				}
				logger.Logger.Error(fmt.Sprintf("[bulk] index error type=%s reason=%s", res.Error.Type, res.Error.Reason))
			},
		}

		if addErr := bulk.Add(ctx, item); addErr != nil {
			return fmt.Errorf("bulk add: %w", addErr)
		}
	}

	if err := bulk.Close(ctx); err != nil {
		return fmt.Errorf("bulk close: %w", err)
	}

	stats := bulk.Stats()
	if stats.NumFailed > 0 {
		return fmt.Errorf("bulk finished with %d failures (ok=%d)", stats.NumFailed, stats.NumFlushed)
	}

	logger.Logger.Info(fmt.Sprintf("Bulk appended successfully: %d docs (flushed=%d)", len(data), stats.NumFlushed))
	return nil
}

func (r *Repository) GetLatestCryptoCurrencyDataById(ctx context.Context, id string) (*repository.CryptoCurrencyData, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}

	query := map[string]interface{}{
		"size": 1,
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": id,
			},
		},
		"sort": []map[string]interface{}{
			{"@timestamp": map[string]string{"order": "desc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error search: %s", res.String())
	}

	var sr struct {
		Hits struct {
			Hits []struct {
				Source repository.CryptoCurrencyData `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	if len(sr.Hits.Hits) == 0 {
		return nil, nil
	}

	return &sr.Hits.Hits[0].Source, nil
}
