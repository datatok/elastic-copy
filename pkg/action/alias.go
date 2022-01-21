package action

import (
	"bytes"
	"regexp"

	log "github.com/sirupsen/logrus"

	"encoding/json"

	"github.com/ebuildy/elastic-copy/pkg/cli"
	"github.com/ebuildy/elastic-copy/pkg/elasticsearch"
)

type AliasAction struct {
	Source, Target, IndicesFilter string
}

type getAliasesResultItem struct {
	Aliases map[string]map[string]interface{}
}

func NewAliasAction(settings *cli.EnvSettings) *AliasAction {
	return &AliasAction{
		Source: "",
	}
}

// https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-aliases.html
func (a *AliasAction) Run() {

	sourceClient := elasticsearch.NewClient(a.Source)
	targetClient := elasticsearch.NewClient(a.Target)

	allAliases := getAliases(sourceClient, a.IndicesFilter)

	log.Infof("found %d aliases to import", len(allAliases))

	var writeActions = make([]map[string]interface{}, 1)

	for aliasIndexName, alias := range allAliases {
		for aliasName, aliasOptions := range alias.Aliases {
			aliasOptions["index"] = aliasIndexName
			aliasOptions["alias"] = aliasName

			action := map[string]interface{}{"add": aliasOptions}

			log.WithField("index", aliasIndexName).
				WithField("alias", aliasName).
				Debug("add")

			writeActions = append(writeActions, action)
		}
	}

	apiBody := map[string]interface{}{"actions": writeActions}
	apiBodyJSON, _ := json.Marshal(apiBody)

	res, err := targetClient.Client.Indices.UpdateAliases(bytes.NewReader(apiBodyJSON))

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	var (
		r map[string]getAliasesResultItem
	)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Info("ok!")
}

func getAliases(sourceClient *elasticsearch.Client, indicesFilter string) map[string]getAliasesResultItem {
	var (
		r map[string]getAliasesResultItem
	)

	res, err := sourceClient.Client.Indices.GetAlias()

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	if len(indicesFilter) == 0 {
		return r
	}

	var filtered = make(map[string]getAliasesResultItem)

	filterExpr := regexp.MustCompile(indicesFilter)

	for indexName, data := range r {
		if filterExpr.MatchString(indexName) {
			filtered[indexName] = data
		}
	}

	return filtered
}
