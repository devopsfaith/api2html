package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

type BackendData struct {
	Obj map[string]interface{}
	Arr []map[string]interface{}
}

// String implements the Stringer interface
func (b *BackendData) String() string {
	d, err := json.MarshalIndent(b, "", "\t")
	log.Println("decoding", b, "as", string(d))
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(d)
}

type Backend func(params map[string]string, headers map[string]string) (*http.Response, error)

type Decoder func(io.Reader) (BackendData, error)

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

func JSONDecoder(r io.Reader) (BackendData, error) {
	var target map[string]interface{}
	err := json.NewDecoder(r).Decode(&target)
	return BackendData{Obj: target}, err
}

func JSONArrayDecoder(r io.Reader) (BackendData, error) {
	var target []map[string]interface{}
	err := json.NewDecoder(r).Decode(&target)
	return BackendData{Arr: target}, err
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
