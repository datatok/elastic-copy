package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Client struct {
	Client *elasticsearch7.Client

	URL, Version string

	Indices []gjson.Result
}

/**
 * Build a new elasticsearch client.
 */
func NewClient(URL string) *Client {
	esClient, _ := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{
			URL,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   12,
			ResponseHeaderTimeout: 30 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	cc := &Client{
		Client: esClient,
		URL:    URL,
	}

	cc.getServerVersion()

	cc.fetchIndices()

	log.WithField("version", cc.Version).
		WithField("URL", URL).
		WithField("indices", len(cc.Indices)).
		Info("Connected to elasticsearch server")

	return cc
}

func (c *Client) fetchIndices() {
	res, err := c.Client.Cat.Indices(c.Client.Cat.Indices.WithFormat("json"))

	if err != nil {
		log.Fatalf("Cant run _cat/indices: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Cant run _cat/indices: %s", res.String())
	}

	j := read(res.Body)

	c.Indices = gjson.Parse(j).Array()
}

/**
 * Get server version from API.
 */
func (c *Client) getServerVersion() {
	var (
		r map[string]interface{}
	)

	res, err := c.Client.Info()

	if err != nil {
		log.Fatalf("Error getting server version: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error getting server version: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	c.Version = r["version"].(map[string]interface{})["number"].(string)
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
