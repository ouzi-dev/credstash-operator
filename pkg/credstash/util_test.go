package credstash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	inputVersion      string
	expectedVersion   string
	expectedErrorText string
}

func TestCredstashSecretGetter_GetCredstashSecretsForCredstashSecretDefs(t *testing.T) {
	tests := []testItem{
		{
			inputVersion:      "1",
			expectedVersion:   "0000000000000000001",
			expectedErrorText: "",
		},
		{
			inputVersion:      "001",
			expectedVersion:   "0000000000000000001",
			expectedErrorText: "",
		},
		{
			inputVersion:      "0000000000000000001",
			expectedVersion:   "0000000000000000001",
			expectedErrorText: "",
		},
		{
			inputVersion:      "00000000000000000000001",
			expectedVersion:   "",
			expectedErrorText: "Version string is longer than supported.",
		},
		{
			inputVersion:      "this is not a number",
			expectedVersion:   "",
			expectedErrorText: "invalid syntax",
		},
	}

	for _, v := range tests {
		actualVersion, actualError := formatCredstashVersion(v.inputVersion)
		assert.Equal(t, v.expectedVersion, actualVersion)
		if actualError != nil {
			assert.Contains(t, actualError.Error(), v.expectedErrorText)
		}
	}
}
