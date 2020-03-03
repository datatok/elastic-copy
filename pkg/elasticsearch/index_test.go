package elasticsearch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndexSelection(t *testing.T) {

	a := assert.New(t)

	c := NewClient("http://localhost:9200")
	es := c.Client

	_, _ = es.Indices.Create("test-1")
	_, _ = es.Indices.Create("test-2")
	_, _ = es.Indices.Create("test-3")

	t.Run("resolve expression", func(t *testing.T) {
		indicesInfo := c.ResolveIndex("test-(.*)")

		a.Len(indicesInfo, 3)
	})



}
