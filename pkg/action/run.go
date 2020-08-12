package action

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/panjf2000/ants/v2"
	"github.com/ebuildy/elastic-copy/pkg/cli"
	"github.com/ebuildy/elastic-copy/pkg/elasticsearch"
	"github.com/ebuildy/elastic-copy/pkg/engine"
	"github.com/ebuildy/elastic-copy/pkg/stdio"
	"github.com/ebuildy/elastic-copy/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type RunAction struct {
	Source, Target, Query, TargetIndexType string
	Indices               []string
	Count                 uint64

	Results map[string]*engine.ProcessResult

	WriteBatchSize, ReadBatchSize, Threads int

	ForceType string
}

type ThreadCopyPayload struct {
	index engine.SourceCollection
	shard engine.CollectionShard
}

func NewRunAction(settings *cli.EnvSettings) *RunAction {
	return &RunAction{
		Results: make(map[string] *engine.ProcessResult),
	}
}

func (a *RunAction) countAction(sourceClient *elasticsearch.Client) {
	var wg sync.WaitGroup

	pCoutDocuments, _ := ants.NewPoolWithFunc(a.Threads, func(i interface{}) {
		index := i.(engine.SourceCollection)

		index.DocumentsCount = sourceClient.CountDocuments(index.Name, a.Query)

		log.WithField("collection", index.Name).
			WithField("documents", index.DocumentsCount).
			WithField("shards", index.ShardsCount).
			Info("start copy")
		
		a.Results[index.Name] = &engine.ProcessResult{
			Index: index,
			CountEntries: index.DocumentsCount,
		}

		wg.Done()
	})

	defer pCoutDocuments.Release()

	for _, indexExp := range a.Indices {
		resolveIndices := sourceClient.ResolveIndex(indexExp)

		for _, index := range resolveIndices {
			wg.Add(1)

			_ = pCoutDocuments.Invoke(index)
		}
	}

	wg.Wait()
}

func (a *RunAction) copyAction(sourceClient *elasticsearch.Client, targetClient engine.Target) {
	var wg sync.WaitGroup

	barCount := uint64(0)

	for _, res := range a.Results {
		barCount = barCount + res.Index.DocumentsCount
	}

	progressBar := pb.Start64(int64(barCount))

	p, _ := ants.NewPoolWithFunc(a.Threads, func(payload interface{}) {
		index := payload.(ThreadCopyPayload).index
		shard := payload.(ThreadCopyPayload).shard

		readQuery := engine.ReadQuery{
			BatchSize: a.ReadBatchSize,
			Query: a.Query,
			Index: index,
			Count: a.Count,
			Shard: shard.ID,
		}

		process := engine.ProcessQuery{
			FailFast:  true,
			BatchSize: a.WriteBatchSize,
			ReadQuery: readQuery,
			TypeOverride: a.ForceType,
		}

		startDate := time.Now().UTC()

		report := sourceClient.Read(process, targetClient, progressBar)

		finishDate := time.Now().UTC()

		report.Duration = finishDate.Sub(startDate)

		log.WithField("reason", report.Reason).
			WithField("count_entries", report.CountEntries).
			WithField("count_added", report.CountAdded).
			WithField("count_errors", report.CountErrors).
			WithField("duration", report.Duration.String()).
			WithField("shard", shard.ID).
			WithField("index", index.Name).
			Info("finish shard copy")

		a.Results[index.Name].IncrementWithShardResult(report)

		wg.Done()
	})

	defer p.Release()

	for _, res := range a.Results {
		index := res.Index

		for _, shard := range index.Shards {

			wg.Add(1)

			_ = p.Invoke(ThreadCopyPayload{
				index: index,
				shard: shard,
			})
		}
	}

	progressBar.SetWriter(os.Stdout)

	progressBar.Start()

	wg.Wait()

	progressBar.Finish()
}

func (a *RunAction) Run() {
	sourceClient := elasticsearch.NewClient(a.Source)

	targetClient := getTarget(a)

	defer ants.Release()

	a.countAction(sourceClient)

	a.copyAction(sourceClient, targetClient)

	for indexName, report := range a.Results {
		log.WithField("collection", indexName).
			WithField("count_entries", report.CountEntries).
			WithField("count_added", report.CountAdded).
			WithField("count_errors", report.CountErrors).
			WithField("duration", report.Duration.String()).
			Info("finish collection copy")
	}
}

func getTarget(a *RunAction) engine.Target {
	var (
		targetClient engine.Target
	)

	if a.Target == utils.TARGET_STDOUT {
		targetClient = stdio.NewClient()
	} else {
		targetClient = elasticsearch.NewClient(a.Target)
	}

	return targetClient
}

