package action

import (
	"fmt"
	"github.com/ebuildy/elastic-copy/pkg/cli"
	"github.com/ebuildy/elastic-copy/pkg/elasticsearch"
	"github.com/ebuildy/elastic-copy/pkg/engine"
	"github.com/ebuildy/elastic-copy/pkg/sample"
	"math/rand"
)

type SampleAction struct {
	Target, Index string
	Count int
}

func NewSampleAction(settings *cli.EnvSettings) *SampleAction {
	return &SampleAction{}
}

func (a *SampleAction) Run() {
	targetClient := elasticsearch.NewClient(a.Target)

	dd := make([]engine.Datum, a.Count)

	for i := 0; i < a.Count ; i++ {
		dd[i] =engine.Datum{
			Index: a.Index,
			Body:  fmt.Sprintf("{\"i\" : %d, \"r\" : %d, \"p\" : %d, \"text\" : \"%s\"}", i, rand.Intn(10000), rand.Int() % 2, sample.Lorem(2)),
		}
	}

	targetClient.Write(engine.ProcessQuery{ FailFast:true}, dd)
}

