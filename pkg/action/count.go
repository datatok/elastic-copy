package action

import (
	"github.com/ebuildy/elastic-copy/pkg/cli"
	"github.com/ebuildy/elastic-copy/pkg/elasticsearch"
	"github.com/ebuildy/elastic-copy/pkg/engine"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type CountAction struct {
	Source, Query string
	Indices []string
}

func NewCountAction(settings *cli.EnvSettings) *CountAction {
	return &CountAction{
		Source: "",
	}
}

func (a *CountAction) Run() {
	sourceClient := elasticsearch.NewClient(a.Source)

	defer ants.Release()

	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(10, func(payload interface{}) {
		index := payload.(engine.SourceCollection)

		index.DocumentsCount = sourceClient.CountDocuments(index.Name, a.Query)

		log.WithField("collection", index.Name).
			WithField("documents", index.DocumentsCount).
			WithField("shards", index.ShardsCount).
			Info("count")

		wg.Done()
	})

	defer p.Release()

	for _, indexExp := range a.Indices {
		resolveIndices := sourceClient.ResolveIndex(indexExp)

		for _, index := range resolveIndices {
			wg.Add(1)

			_ = p.Invoke(index)
		}
	}

	wg.Wait()
}