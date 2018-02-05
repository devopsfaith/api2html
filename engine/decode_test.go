package engine

import (
	"bytes"
	"testing"
)

func TestJSONDecoder(t *testing.T) {
	r := ResponseContext{}
	if err := JSONDecoder(bytes.NewBufferString(`{"a":"b"}`), &r); err != nil {
		t.Error(err)
		return
	}
	if len(r.Array) != 0 {
		t.Errorf("unexpected array value: %v", r.Array)
	}
	if v, ok := r.Data["a"]; !ok || "b" != v.(string) {
		t.Errorf("unexpected obj value: %v", r.Data)
	}
}

func TestJSONArrayDecoder(t *testing.T) {
	r := ResponseContext{}
	if err := JSONArrayDecoder(bytes.NewBufferString(`[{"a":"b"}]`), &r); err != nil {
		t.Error(err)
		return
	}
	if len(r.Array) != 1 {
		t.Errorf("unexpected array value: %v", r.Array)
	}
	if len(r.Data) != 0 {
		t.Errorf("unexpected obj value: %v", r.Data)
	}
}
