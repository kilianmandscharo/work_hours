package datetime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidRFC3339(t *testing.T) {
	tests := []struct {
		name       string
		dateString string
		expected   bool
	}{
		{
			name:       "Valid RFC3339 Date",
			dateString: "2023-07-23T12:34:56Z",
			expected:   true,
		},
		{
			name:       "Invalid RFC3339 Date (Missing Timezone)",
			dateString: "2023-07-23T12:34:56",
			expected:   false,
		},
		{
			name:       "Invalid RFC3339 Date (Incorrect Format)",
			dateString: "2023-07-23T12:34:56.789.111Z",
			expected:   false,
		},
		{
			name:       "Invalid RFC3339 Date (Invalid Characters)",
			dateString: "2023-07-23T12:34:56Zinvalid",
			expected:   false,
		},
		{
			name:       "Invalid RFC3339 Date (Empty String)",
			dateString: "",
			expected:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsValidRFC3339(test.dateString))
		})
	}
}
