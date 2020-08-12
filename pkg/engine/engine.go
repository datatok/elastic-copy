package engine

import "time"

const READ_FINISH_REASON_ERROR = "error"
const READ_FINISH_REASON_COUNT = "count"
const READ_FINISH_REASON_ALL = "all"

const WRITE_RESULT_OK = "ok"
const WRITE_RESULT_FAILFAST = "fail_fast"

type Datum struct {
	Type, Index, ID string
	Body string
}

type Target interface {
	Write(process ProcessQuery, data [] Datum) WriteResult
}

type WriteResult struct {
	Reason, Error string
	CountEntries, CountAdded, CountUpdated, CountErrors uint64
	Duration time.Duration
}

type ProcessShardResult struct {
	Reason, Error string
	CountEntries, CountAdded, CountUpdated, CountErrors uint64
	Duration time.Duration
}

type ProcessResult struct {
	Index SourceCollection
	CountEntries, CountAdded, CountUpdated, CountErrors uint64
	Duration time.Duration
}

type ReadQuery struct {
	Index        SourceCollection
	Query, Shard string
	Count        uint64
	BatchSize 	 int
}

type ProcessQuery struct {
	FailFast bool
	BatchSize int
	ReadQuery ReadQuery
	TypeOverride string
}

/**
 * ES index
 */
type SourceCollection struct {
	Name           string
	ShardsCount    int
	DocumentsCount uint64
	Shards         []CollectionShard
}

/**
 * ES index shard
 */
type CollectionShard struct {
	ID, Name string
}

func (p *ProcessShardResult) IncrementWithWriteResult(result WriteResult) {
	p.CountAdded += result.CountAdded
	p.CountEntries += result.CountEntries
	p.CountErrors += result.CountErrors
	p.CountUpdated += result.CountUpdated
	p.Duration += result.Duration
}

func (p *ProcessResult) IncrementWithShardResult(result ProcessShardResult) {
	p.CountAdded += result.CountAdded
	p.CountEntries += result.CountEntries
	p.CountErrors += result.CountErrors
	p.CountUpdated += result.CountUpdated
	p.Duration += result.Duration
}
