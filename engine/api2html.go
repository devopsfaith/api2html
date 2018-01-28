package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Pages            []Page                 `json:"pages"`
	StaticTXTContent []string               `json:"static_txt_content"`
	Robots           bool                   `json:"robots"`
	Sitemap          bool                   `json:"sitemap"`
	Templates        map[string]string      `json:"templates"`
	Layouts          map[string]string      `json:"layouts"`
	Extra            map[string]interface{} `json:"extra"`
}

type Page struct {
	Name              string
	URLPattern        string
	BackendURLPattern string
	Template          string
	Layout            string
	CacheTTL          string
	Header            string
	IsArray           bool
	Extra             map[string]interface{}
}

type Backend func(map[string]string, map[string]string) (*http.Response, error)

type Decoder func(io.Reader) (map[string]interface{}, error)

type Renderer interface {
	Render(io.Writer, interface{}) error
}

type Subscription struct {
	Name string
	In   chan Renderer
}

func ErrorPlaceHolder(_ map[string]string) (*http.Response, error) {
	return nil, ErrNoBackendDefined
}

func JSONDecoder(r io.Reader) (map[string]interface{}, error) {
	var target map[string]interface{}
	err := json.NewDecoder(r).Decode(&target)
	return map[string]interface{}{"data": target}, err
}

func JSONArrayDecoder(r io.Reader) (map[string]interface{}, error) {
	var target []map[string]interface{}
	err := json.NewDecoder(r).Decode(&target)
	return map[string]interface{}{"data": target}, err
}

type ErrorRenderer struct {
	Error error
}

func (r ErrorRenderer) Render(_ io.Writer, _ interface{}) error {
	return r.Error
}

var (
	ErrNoResponseGeneratorDefined = fmt.Errorf("no response generator defined")
	ErrNoBackendDefined           = fmt.Errorf("no backend defined")
	ErrNoRendererDefined          = fmt.Errorf("no rendered defined")
	EmptyRenderer                 = ErrorRenderer{ErrNoRendererDefined}
)
