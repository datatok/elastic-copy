package elasticsearch

import (
	"github.com/panjf2000/ants/v2"
	"github.com/ebuildy/elastic-copy/pkg/engine"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
	"sync"
)

func (c *Client) ResolveIndex(indexExpression string) []engine.SourceCollection {
	ret := make([]engine.SourceCollection, 0)

	indexExpressionReg := regexp.MustCompile(indexExpression)

	defer ants.Release()

	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(10, func(payload interface{}) {
		indexData := payload.(gjson.Result)
		indexFoundInfo := c.getIndexInfo(indexData)
		ret = append(ret, *indexFoundInfo)

		wg.Done()
	})
	defer p.Release()

	for _, indexData := range c.Indices {
		testIndexName := indexData.Get("index").String()

		if testIndexName == indexExpression ||
			(indexExpressionReg != nil && indexExpressionReg.MatchString(testIndexName)) {
			wg.Add(1)
			_ = p.Invoke(indexData)
		}
	}

	wg.Wait()

	return ret
}

/**
 * Get index details (shards, nodes, count).
 */
func (c *Client) getIndexInfo(index gjson.Result) *engine.SourceCollection {
	indexName := index.Get("index").String()

	shards := c.getIndexShardsInfo(indexName)

	return &engine.SourceCollection{
		Name:           indexName,
		Shards:         shards,
		ShardsCount:	len(shards),
		DocumentsCount: index.Get("docs.count").Uint(),
	}
}

func (c *Client) CountDocuments(index string, query string) uint64 {
	es := c.Client

	res, err := es.Count(
		es.Count.WithIndex(index),
		es.Count.WithQuery(query),
	)

	if err != nil {
		log.WithField("index", index).Fatalf("Error getting index shards info: %s", err)
	}

	if res.IsError() {
		log.WithField("index", index).Fatalf("Error getting index shards info: %s", res.String())
	}

	json := read(res.Body)

	res.Body.Close()

	return gjson.Get(json, "count").Uint()
}

func (c *Client) getIndexShardsInfo(index string) []engine.CollectionShard {
	res, err := c.Client.Indices.ShardStores(
		c.Client.Indices.ShardStores.WithIndex(index),
		c.Client.Indices.ShardStores.WithStatus("green,yellow"),
	)

	if err != nil {
		log.WithField("index", index).Fatalf("Error getting index shards info: %s", err)
	}

	json := read(res.Body)

	res.Body.Close()

	if res.IsError() {
		log.WithField("index", index).Fatalf("Error getting index shards info: %s", res.String())
	}

	shards := make([]engine.CollectionShard, 0)

	indexEscaped := strings.Replace(index, ".", "\\.", -1)

	gjson.Get(json, "indices." + indexEscaped + ".shards").ForEach(func(k gjson.Result, v gjson.Result) bool {
		shards = append(shards, engine.CollectionShard{
			ID:   k.String(),
			Name: k.String(),
		})
		return true
	})

	return shards
}
