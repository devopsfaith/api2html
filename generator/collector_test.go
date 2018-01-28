package generator

import (
	"fmt"
	"os"
	"testing"
)

func TestNewCollector(t *testing.T) {
	pwd := os.Getenv("PWD")

	collector := NewCollector("unknownConfig1", "unknownI18N1")
	availableISOs := collector.AvailableISOs()
	expectedSize := 0
	if len(availableISOs) != expectedSize {
		t.Errorf("unexpected size for availableISOs. have %d, want %d", len(availableISOs), expectedSize)
	}
	_, err := collector.Collect("unknown_iso_1")
	if err == nil {
		t.Error("expecting error!")
		return
	}
	if err.Error() != "open unknownI18N1/unknown_iso_1.ini: no such file or directory" {
		t.Error("collecting from unknown translation folder", err)
		return
	}

	collector = NewCollector("unknownConfig2", fmt.Sprintf("%s/test/i18n", pwd))
	availableISOs = collector.AvailableISOs()
	expectedSize = 3
	if len(availableISOs) != expectedSize {
		t.Errorf("unexpected size for availableISOs. have %d, want %d", len(availableISOs), expectedSize)
	}
	_, err = collector.Collect("en-US")
	if err != ErrNoConfig {
		t.Error("expecting ErrNoConfig. got:", err)
		return
	}

	collector = NewCollector(fmt.Sprintf("%s/test/config", pwd), fmt.Sprintf("%s/test/i18n", pwd))
	availableISOs = collector.AvailableISOs()
	expectedSize = 3
	if len(availableISOs) != expectedSize {
		t.Errorf("unexpected size for availableISOs. have %d, want %d", len(availableISOs), expectedSize)
	}

	for _, iso := range collector.AvailableISOs() {
		data, err := collector.Collect(iso)
		if err != nil {
			t.Error("collecting", iso, err)
			return
		}
		fmt.Println(iso, data)
	}
}

func TestData_String(t *testing.T) {
	data := Data{
		Config: map[string]Map{
			"cfg1": {
				"key1": "value1",
				"key2": "value2",
			},
		},
		I18N: map[string]Map{
			"cfg1": {
				"literal1": "literal_value_1",
			},
		},
	}
	if data.String() != `{"I18N":{"cfg1":{"literal1":"literal_value_1"}},"Config":{"cfg1":{"key1":"value1","key2":"value2"}}}` {
		t.Error("unexpected data serialization result:", data)
	}
	if data.Config["cfg1"].String() != `{"key1":"value1","key2":"value2"}` {
		t.Error("unexpected config serialization result:", data)
	}
}
