package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/elastic/go-elasticsearch/v7"

	"github.com/tinrab/meower/schema"
)

type ElasticRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}
	_, err = client.Info()
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) InsertMeow(ctx context.Context, meow schema.Meow) error {
	body, _ := json.Marshal(meow)
	_, err := r.client.Index(
		"meows",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(meow.ID),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (r *ElasticRepository) SearchMeows(ctx context.Context, query string, skip uint64, take uint64) (result []schema.Meow, err error) {
	var buf bytes.Buffer
	reqBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"body"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		},
	}
	if err = json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("meows"),
		r.client.Search.WithFrom(int(skip)),
		r.client.Search.WithSize(int(take)),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			result = nil
		}
	}()
	if res.IsError() {
		return nil, errors.New("search failed")
	}

	type Response struct {
		Took int64
		Hits struct {
			Total struct {
				Value int64
			}
			Hits []*struct {
				Source schema.Meow `json:"_source"`
			}
		}
	}
	resBody := Response{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	var meows []schema.Meow
	for _, hit := range resBody.Hits.Hits {
		meows = append(meows, hit.Source)
	}
	return meows, nil
}
