package engine

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseConfigFromFile_noFile(t *testing.T) {
	c, err := ParseConfigFromFile("unknown")
	if err == nil {
		t.Error("error expected")
	}
	if len(c.Templates) != 0 {
		t.Error("unexpected number of templates:", c.Templates)
	}
	if len(c.Pages) != 0 {
		t.Error("unexpected number of pages:", c.Pages)
	}
}

func TestParseConfigFromFile_wrongConfig(t *testing.T) {
	f, err := ioutil.TempFile(".", "wrong_config")
	if err != nil {
		t.Error(err)
		return
	}
	if _, err := f.WriteString("{"); err != nil {
		t.Error(err)
		return
	}
	f.Close()

	c, err := ParseConfigFromFile(f.Name())
	if err == nil {
		t.Error("error expected")
	}
	if len(c.Templates) != 0 {
		t.Error("unexpected number of templates:", c.Templates)
	}
	if len(c.Pages) != 0 {
		t.Error("unexpected number of pages:", c.Pages)
	}
	os.Remove(f.Name())
}

func TestParseConfig_wrongConfig(t *testing.T) {
	for i, subject := range []string{
		`{
	"templates":{
}`,
	} {
		c, err := ParseConfig(bytes.NewBufferString(subject))
		if err == nil {
			t.Error("error expected in", i)
		}
		if len(c.Templates) != 0 {
			t.Error("unexpected number of templates:", c.Templates)
		}
		if len(c.Pages) != 0 {
			t.Error("unexpected number of pages:", c.Pages)
		}
	}
}

func TestParseConfig_ok(t *testing.T) {
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

func TestParseConfig_extra(t *testing.T) {
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
			"CacheTTL": "3600s",
			"extra": {
				"b": true
			}
		}
	],
	"extra": {
		"a": {
			"a1": 42
		}
	}
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
	for i, p := range c.Pages {
		if tmp, ok := p.Extra["a"].(map[string]interface{}); !ok {
			t.Errorf("the page #%d has a wrong extra. have: %v", i, p.Extra)
		} else if v, ok := tmp["a1"].(float64); !ok || v != 42 {
			t.Errorf("the page #%d has a wrong extra['a']. have: %v (%f)", i, tmp, v)
		}
		if p.Name == "page02" {
			if tmp, ok := p.Extra["b"].(bool); !ok || !tmp {
				t.Errorf("the page #%d has a wrong extra. have: %v", i, p.Extra)
			}
		}
	}
}
