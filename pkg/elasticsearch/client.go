package elasticsearch

import (
	"coin-feed/pkg/logger"

	"github.com/elastic/go-elasticsearch/v9"
	"go.uber.org/zap"
)

type Client struct {
	*elasticsearch.Client
}

func NewClient(addresses []string, username, password string, caCert []byte) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
		//CACert:    caCert,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		logger.Logger.Error("Error connecting to Elasticsearch", zap.Error(err))
		return nil, err
	}

	logger.Logger.Info("Successfully connected to Elasticsearch")
	return &Client{es}, nil
}
