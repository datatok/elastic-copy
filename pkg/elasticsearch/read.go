package elasticsearch

import (
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/ebuildy/elastic-copy/pkg/engine"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func (c *Client) Read(process engine.ProcessQuery, writer engine.Target, progressReporter *pb.ProgressBar) engine.ProcessShardResult {
	var (
		batchNum      int
		foundItems    uint64
		scrollID      string
		index         = process.ReadQuery.Index
		es            = c.Client
		err           error
		res           *esapi.Response
		processResult engine.ProcessShardResult
		query         = process.ReadQuery
	)

	log.Info(strings.NewReader(query.Query))

	for {
		batchNum++

		// First, get scroll ID
		if len(scrollID) == 0 {
			res, err = es.Search(
				es.Search.WithIndex(index.Name),
				es.Search.WithSize(process.BatchSize),
				//es.Search.WithQuery(query.Query),
				es.Search.WithBody(strings.NewReader(query.Query)),
				es.Search.WithScroll(time.Minute),
				es.Search.WithPreference("_shards:"+query.Shard),
			)
		} else {
			res, err = es.Scroll(
				es.Scroll.WithScrollID(scrollID),
				es.Scroll.WithScroll(time.Minute),
			)
		}

		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		if res.IsError() {
			log.Fatalf("Error response: %s", res)
		}

		json := read(res.Body)
		res.Body.Close()

		// Extract the scrollID from response
		//
		scrollID = gjson.Get(json, "_scroll_id").String()

		// Extract the search results
		//
		hits := gjson.Get(json, "hits.hits").Array()

		// Break out of the loop when there are no results
		//
		if len(hits) < 1 {
			break
		} else {
			ret := make([]engine.Datum, 0)

			for _, hit := range hits {
				if query.Count != 0 && foundItems >= query.Count {
					break
				}

				ret = append(ret, engine.Datum{
					Index: index.Name,
					ID:    hit.Get("_id").String(),
					Body:  hit.Get("_source").String(),
				})

				foundItems++
			}

			writeResult := writer.Write(process, ret)

			if progressReporter != nil {
				progressReporter.Add64(int64(writeResult.CountEntries))
			}

			processResult.IncrementWithWriteResult(writeResult)

			if query.Count != 0 && foundItems >= query.Count {
				processResult.Reason = engine.READ_FINISH_REASON_COUNT

				return processResult
			}
		}
	}

	processResult.Reason = engine.READ_FINISH_REASON_ALL

	return processResult
}
