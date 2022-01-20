package elasticsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractQuery(t *testing.T) {
	a := assert.New(t)

	t.Run("empty query", func(t *testing.T) {
		query := "{}"

		r := extractQuery(query)

		a.Equal(r, "")
	})

	t.Run("only query", func(t *testing.T) {
		query := "{\"query\":\"toto\"}"

		r := extractQuery(query)

		a.Equal(query, r)
	})

	t.Run("query + sort", func(t *testing.T) {
		query := "{\"query\":\"toto\", \"sort\" : \"toto\"}"

		r := extractQuery(query)

		a.Equal("{\"query\":\"toto\"}", r)
	})

	t.Run("only sort", func(t *testing.T) {
		query := "{\"sort\" : \"toto\"}"

		r := extractQuery(query)

		a.Equal("", r)
	})
}

func __TestIndexSelection(t *testing.T) {

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
