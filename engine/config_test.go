package engine

import (
	"bytes"
	"testing"
)

func TestParseConfig(t *testing.T) {
	configContent := `{
	"templates":{
		"template01": "path01",
		"template02": "path02",
		"template03": "path03"
	},
	"pages":[
		{
			"name": "page01",
			"URLPattern": "/page-01",
			"BackendURLPattern": "https://jsonplaceholder.typicode.com/users/1",
			"Template": "template01",
			"CacheTTL": "3600s"
		},
		{
			"name": "page02",
			"URLPattern": "/page-02",
			"BackendURLPattern": "https://jsonplaceholder.typicode.com/users/2",
			"Template": "template02",
			"CacheTTL": "3600s"
		}
	]
}`
	c, err := ParseConfig(bytes.NewBufferString(configContent))
	if err != nil {
		t.Error(err)
	}
	if len(c.Templates) != 3 {
		t.Error("unexpected number of templates:", c.Templates)
	}
	if len(c.Pages) != 2 {
		t.Error("unexpected number of pages:", c.Pages)
	}
}
