package engine

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseContext is the struct ready to rendered and returned to the Handler
type ResponseContext struct {
	// Data cotains the backend data
	Data  BackendData
	Extra map[string]interface{}
	// Params stores the params of the request
	Params map[string]string
	// Helper is a struct containing a few basic template helpers
	Helper interface{}
	// 	Context is a reference to the gin context for the request
	Context *gin.Context
}

// ResponseGenerator is a function that, given a gin request, returns a response struc and an error
type ResponseGenerator func(*gin.Context) (ResponseContext, error)

// NoopResponse is a ResponseGenerator that always returns an empty response and the
// ErrNoResponseGeneratorDefined error
func NoopResponse(_ *gin.Context) (ResponseContext, error) {
	return ResponseContext{}, ErrNoResponseGeneratorDefined
}

// StaticResponseGenerator is a ResponseGenerator that creates a response just by adding the
// default response values to the ResponseContext and a zero value BackendData
type StaticResponseGenerator struct {
	Page Page
}

// ResponseGenerator implements the ResponseGenerator interface
func (s *StaticResponseGenerator) ResponseGenerator(c *gin.Context) (ResponseContext, error) {
	params := map[string]string{}
	for _, v := range c.Params {
		params[v.Key] = v.Value
	}
	target := ResponseContext{
		Extra:   s.Page.Extra,
		Context: c,
		Params:  params,
		Helper:  &tplHelper{},
	}
	return target, nil
}

// DynamicResponseGenerator is a ResponseGenerator that creates a response by adding the decoded data
// returned by the Backend wo the default response values. Depending on the selected decoder,
// the generated responses may have the backend data stored at the `Obj` or at the `Arr` part
type DynamicResponseGenerator struct {
	Page    Page
	Backend Backend
	Decoder Decoder
}

// ResponseGenerator implements the ResponseGenerator interface
func (drg *DynamicResponseGenerator) ResponseGenerator(c *gin.Context) (ResponseContext, error) {
	params := map[string]string{}
	for _, v := range c.Params {
		params[v.Key] = v.Value
	}
	headers := map[string]string{}
	h := c.Request.Header.Get(drg.Page.Header)
	if h != "" {
		headers[drg.Page.Header] = h
	}
	result := ResponseContext{
		Extra:   drg.Page.Extra,
		Context: c,
		Params:  params,
		Helper:  &tplHelper{},
	}
	resp, err := drg.Backend(params, headers)
	if err != nil {
		return result, err
	}

	result.Data, err = drg.Decoder(resp.Body)
	resp.Body.Close()
	return result, err
}

type tplHelper struct {
}

func (tplHelper) Now() string {
	return time.Now().String()
}
