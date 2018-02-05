package engine

import (
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
	PublicFolder     *PublicFolder          `json:"public_folder"`
}

type PublicFolder struct {
	Path   string `json:"path_to_folder"`
	Prefix string `json:"url_prefix"`
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

type Backend func(params map[string]string, headers map[string]string) (*http.Response, error)

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
