package elasticsearch

import (
	"bytes"
	"fmt"

	"github.com/ebuildy/elastic-copy/pkg/engine"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func (c *Client) Write(process engine.ProcessQuery, data []engine.Datum) engine.WriteResult {

	var (
		buf bytes.Buffer
		res *esapi.Response
		err error

		numItems   uint64
		numErrors  uint64
		numNOOP    uint64
		numIndexed uint64
		numBatches int
		currBatch  int

		meta []byte
	)

	count := len(data)
	batch := process.BatchSize
	bEOL := []byte("\n")
	es := c.Client
	index := data[0].Index
	indexType := data[0].Type

	if count%batch == 0 {
		numBatches = count / batch
	} else {
		numBatches = (count / batch) + 1
	}

	log.WithField("index", index).
		Debugf("bulk start: %d batches", numBatches)

	for i, datum := range data {
		numItems++

		currBatch = i / batch

		if i == count-1 {
			currBatch++
		}

		if len(datum.ID) == 0 {
			meta = []byte(`{ "index" : { } }`)
		} else {
			meta = []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }`, datum.ID))
		}

		lineData := []byte(datum.Body)

		buf.Grow(len(meta) + len(data) + 2)
		buf.Write(meta)
		buf.Write(bEOL)
		buf.Write(lineData)
		buf.Write(bEOL)

		if i > 0 && i%batch == 0 || i == count-1 {
			/*log.WithField("index", index).
			WithField("batch_id", currBatch).
			Debug("bulk sending")
			*/

			retryCounter := 0

			for retryCounter < 5 {
				bulkReader := bytes.NewReader(buf.Bytes())
				err = nil

				if len(indexType) == 0 || process.TypeOverride == "_" {
					res, err = es.Bulk(bulkReader, es.Bulk.WithIndex(index))
				} else if len(process.TypeOverride) > 0 {
					res, err = es.Bulk(bulkReader, es.Bulk.WithIndex(index), es.Bulk.WithDocumentType(process.TypeOverride))
				} else {
					res, err = es.Bulk(bulkReader, es.Bulk.WithIndex(index), es.Bulk.WithDocumentType(indexType))
				}

				if err == nil {
					break
				}

				log.Debugf("retrying %d/5 because error: %s", retryCounter, err)

				retryCounter++
			}

			if err != nil {
				log.WithField("index", index).
					WithField("batch_id", currBatch).
					Warnf("bulk error: %s", err)

				if process.FailFast {
					log.Fatal("error detected, fail fast!")
				}
			}

			json := read(res.Body)
			res.Body.Close()

			if res.IsError() {
				numErrors += numItems

				log.WithField("index", index).
					WithField("batch_id", currBatch).
					Warnf("bulk error: %s", json)

				if process.FailFast {
					log.Fatal("error detected, fail fast!")
				}
			} else {
				gjson.Get(json, "items").ForEach(func(k gjson.Result, v gjson.Result) bool {
					switch status := v.Get("index.status").Int(); status {
					case 200:
						numNOOP++
					case 201:
						numIndexed++
					default:
						numErrors++

						log.WithField("index", index).
							WithField("status", status).
							WithField("batch_id", currBatch).
							Warnf("bulk error: %s", v.Get("index.error.reason"))

						if process.FailFast {
							log.Fatal("error detected, fail fast!")
						}
					}

					return true
				})
			}

			buf.Reset()
			numItems = 0
		}
	}

	return engine.WriteResult{
		Reason:       engine.WRITE_RESULT_OK,
		Error:        "",
		CountEntries: uint64(count),
		CountAdded:   numIndexed,
		CountUpdated: numNOOP,
		CountErrors:  numErrors,
		Duration:     0,
	}
}
