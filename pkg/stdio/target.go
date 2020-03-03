package stdio

import (
	"fmt"
	"github.com/ebuildy/elastic-copy/pkg/engine"
)

type StdIO struct {

}

func NewClient() *StdIO {
	return &StdIO{}
}

func (s *StdIO) Write(process engine.ProcessQuery, data [] engine.Datum) engine.WriteResult {
	for _, datum := range data {
		fmt.Println(datum.Body)
	}

	return engine.WriteResult{
		Reason:       engine.WRITE_RESULT_OK,
		Error:        "",
		CountEntries: uint64(len(data)),
		CountAdded:   uint64(len(data)),
		CountUpdated: 0,
		CountErrors:  0,
		Duration:     0,
	}
}
