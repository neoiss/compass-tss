package memo

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"
)

func toString(v any) string {
	marshal, _ := json.Marshal(v)
	return string(marshal)
}

func Test_parseAffiliates(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []*Affiliate
	}{
		{
			name:  "empty-string",
			input: "",
			want:  []*Affiliate{},
		},
		{
			name:  "string-with-only-spaces",
			input: "   ",
			want:  []*Affiliate{},
		},
		{
			name:  "single-uncompressed-format",
			input: "flashx:15",
			want: []*Affiliate{
				{
					Name:       "flashx",
					Bps:        big.NewInt(15),
					Compressed: false,
				},
			},
		},
		{
			name:  "multiple-uncompressed-format",
			input: "flashx:15,butter:20,imtoken:10",
			want: []*Affiliate{
				{
					Name:       "flashx",
					Bps:        big.NewInt(15),
					Compressed: false,
				},
				{
					Name:       "butter",
					Bps:        big.NewInt(20),
					Compressed: false,
				},
				{
					Name:       "imtoken",
					Bps:        big.NewInt(10),
					Compressed: false,
				},
			},
		},
		{
			name:  "single-compressed-format",
			input: "fx15",
			want: []*Affiliate{
				{
					Name:       "fx",
					Bps:        big.NewInt(15),
					Compressed: true,
				},
			},
		},
		{
			name:  "multiple-compressed-format",
			input: "fx15b20im10",
			want: []*Affiliate{
				{
					Name:       "fx",
					Bps:        big.NewInt(15),
					Compressed: true,
				},
				{
					Name:       "b",
					Bps:        big.NewInt(20),
					Compressed: true,
				},
				{
					Name:       "im",
					Bps:        big.NewInt(10),
					Compressed: true,
				},
			},
		},
		{
			name:  "compressed-format-with-underscore-and-hyphen",
			input: "tk-:100,a_b:200",
			want: []*Affiliate{
				{
					Name:       "tk-",
					Bps:        big.NewInt(100),
					Compressed: false,
				},
				{
					Name:       "a_b",
					Bps:        big.NewInt(200),
					Compressed: false,
				},
			},
		},
		{
			name:  "invalid-uncompressed-format--no-colon",
			input: "flashx15",
			want:  []*Affiliate{},
		},
		{
			name:  "invalid-uncompressed-format--invalid-bps",
			input: "flashx:abc",
			want:  []*Affiliate{},
		},
		{
			name:  "invalid-compressed-format--only-letters",
			input: "abcd",
			want:  []*Affiliate{},
		},
		{
			name:  "invalid-compressed-format--only-numbers",
			input: "1234",
			want:  []*Affiliate{},
		},
		{ // todo
			name:  "compressed-format-with-special-characters",
			input: "a1@~%b2#c3",
			want:  []*Affiliate{},
		},
		{
			name:  "uncompressed-format-with-spaces",
			input: " flashx : 15 , butter : 20 ",
			want: []*Affiliate{
				{
					Name:       "flashx",
					Bps:        big.NewInt(15),
					Compressed: false,
				},
				{
					Name:       "butter",
					Bps:        big.NewInt(20),
					Compressed: false,
				},
			},
		},
		{
			name:  "compressed-format-with-spaces",
			input: " fx  15 b20 ",
			want:  []*Affiliate{},
		},
		{
			name:  "compressed-format-with-zero-bps",
			input: "a0",
			want: []*Affiliate{
				{
					Name:       "a",
					Bps:        big.NewInt(0),
					Compressed: true,
				},
			},
		},
		{
			name:  "compressed-format-with-negative-bps (impossible with digits)",
			input: "a-1",
			want: []*Affiliate{
				{
					Name:       "a-",
					Bps:        big.NewInt(1),
					Compressed: true,
				},
			},
		},
		{
			name:  "large-bps-value-within-limit",
			input: "tp9999",
			want: []*Affiliate{
				{
					Name:       "tp",
					Bps:        big.NewInt(9999),
					Compressed: true,
				},
			},
		},
		{
			name:  "name-too-long",
			input: "test10",
			want:  []*Affiliate{},
		},
	}

	for _, tt := range tests {
		//t.Run(tt.name, func(t *testing.T) {
		//	if got := parseAffiliates(tt.input); !reflect.DeepEqual(got, tt.want) {
		//		t.Errorf("parseAffiliates() = %v, want %v", got, tt.want)
		//	}
		//})

		t.Run(tt.name, func(t *testing.T) {
			got := new(parser).parseAffiliates(tt.input)

			if len(got) != len(tt.want) {
				t.Errorf("parseAffiliates(%q), got %v, want %v", tt.input, toString(got), toString(tt.want))
				return
			}

			for i := range got {
				if got[i].Name != tt.want[i].Name {
					t.Errorf("parseAffiliates(%q)[%d].Name = %s, want %s", tt.input, i, got[i].Name, tt.want[i].Name)
				}

				if got[i].Bps.Cmp(tt.want[i].Bps) != 0 {
					t.Errorf("parseAffiliates(%q)[%d].Bps = %s, want %s", tt.input, i, got[i].Bps.String(), tt.want[i].Bps.String())
				}

				if got[i].Compressed != tt.want[i].Compressed {
					t.Errorf("parseAffiliates(%q)[%d].Compressed = %t, want %t", tt.input, i, got[i].Compressed, tt.want[i].Compressed)
				}
			}
		})
	}
}

func Test_parseMinAmount(t *testing.T) {
	large, _ := new(big.Int).SetString("999999999999999999999999999999", 10)
	// Define test cases with input and expected output
	tests := []struct {
		name     string
		input    string
		expected *big.Int
	}{
		{
			name:     "empty-string",
			input:    "",
			expected: big.NewInt(0),
		},
		{
			name:     "string-with-only-spaces",
			input:    "   ",
			expected: big.NewInt(0),
		},
		{
			name:     "single space mixed with number",
			input:    " 123 ",
			expected: big.NewInt(0),
		},
		{
			name:     "multiple-spaces-in-number-string",
			input:    "1 2 3",
			expected: big.NewInt(0), // Invalid format after removing spaces
		},
		{
			name:     "positive-integer",
			input:    "123",
			expected: big.NewInt(123),
		},
		//{
		//	name:     "negative integer",
		//	input:    "-456",
		//	expected: big.NewInt(-456),
		//},
		{
			name:     "zero",
			input:    "0",
			expected: big.NewInt(0),
		},
		{
			name:     "scientific-notation-with-lowercase-e",
			input:    "1.23e2",
			expected: big.NewInt(123),
		},
		{
			name:     "scientific-notation-with-uppercase-E",
			input:    "1.5E3",
			expected: big.NewInt(1500),
		},
		//{
		//	name:     "scientific notation with negative exponent",
		//	input:    "1.23e-2", // This would result in 0 when converted to int
		//	expected: big.NewInt(0),
		//},
		{
			name:     "large integer",
			input:    "999999999999999999999999999999",
			expected: large,
		},
		{
			name:     "invalid-string-without-e/E",
			input:    "abc",
			expected: big.NewInt(0),
		},
		{
			name:     "invalid-string-with-e-but-not-scientific-notation",
			input:    "abc123e",
			expected: big.NewInt(0),
		},
		{
			name:     "mixed-valid-number-with-spaces",
			input:    " 456 ",
			expected: big.NewInt(0),
		},
		{
			name:     "decimal number without scientific notation (should fail)",
			input:    "123.45",
			expected: big.NewInt(0),
		},
		{
			name:     "negative-decimal-number-without-scientific-notation (should fail)",
			input:    "-123.45",
			expected: big.NewInt(0),
		},
		{
			name:     "1e0",
			input:    "1e0",
			expected: big.NewInt(1),
		},
		{
			name:     "2E1",
			input:    "2E1",
			expected: big.NewInt(20),
		},
		{
			name:     "3e2",
			input:    "3e2",
			expected: big.NewInt(300),
		},
		{
			name:     "12345e8",
			input:    "12345e8",
			expected: big.NewInt(1234500000000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := new(parser).parseMinAmount(tt.input)

			// Compare the values of the big.Ints
			if result.Cmp(tt.expected) != 0 {
				t.Errorf("parseMinAmount(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func Test_parseMinAmountScientificNotation(t *testing.T) {
	// Additional specific tests for scientific notation
	scienceTests := []struct {
		input    string
		expected int64
	}{
		{"1e0", 1},
		{"2E1", 20},
		{"3e2", 300},
		{"1.5e1", 15},
		{"1.78e2", 178},
		{"12345e8", 1234500000000},
	}

	for _, tt := range scienceTests {
		t.Run(tt.input, func(t *testing.T) {
			result := new(parser).parseMinAmount(tt.input)
			expected := big.NewInt(tt.expected)

			if result.Cmp(expected) != 0 {
				t.Errorf("parseMinAmount(%q) = %v, want %v", tt.input, result, expected)
			}
		})
	}
}

// TestParseMinAmountWhitespaceHandling tests whitespace removal functionality
func TestParseMinAmountWhitespaceHandling(t *testing.T) {
	whitespaceTests := []struct {
		input    string
		expected *big.Int
	}{
		{" 123", big.NewInt(123)},
		{"123 ", big.NewInt(123)},
		{" 123 ", big.NewInt(123)},
		{"\t123\n", big.NewInt(123)}, // tabs and newlines are NOT removed by Replace
		{"  456  ", big.NewInt(456)},
		{"   ", big.NewInt(0)},
		{" \t \n ", big.NewInt(0)}, // Only spaces are removed
	}

	for _, tt := range whitespaceTests {
		t.Run(tt.input, func(t *testing.T) {
			result := new(parser).parseMinAmount(tt.input)

			if result.Cmp(tt.expected) != 0 {
				t.Errorf("parseMinAmount(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Additional helper test to verify the flush function behavior
func TestParseAffiliatesEdgeCases(t *testing.T) {
	// Test case for name too long (though this requires modifying maxNameLength or creating a specific scenario)
	veryLongName := strings.Repeat("a", maxNameLength+1) + "123"
	result := new(parser).parseAffiliates(veryLongName)

	// With current implementation, if name exceeds maxNameLength, both builders are reset
	// So we expect an empty result
	if len(result) != 0 {
		t.Errorf("want empty result for name too long, got %d affiliates", len(result))
	}

	// Test consecutive numbers without letters before
	result2 := new(parser).parseAffiliates("123456")
	if len(result2) != 0 {
		t.Errorf("want empty result for numbers without preceding names, got %d affiliates", len(result2))
	}

	// Test single letter followed by number
	result3 := new(parser).parseAffiliates("a1")
	if len(result3) != 1 {
		t.Errorf("want 1 affiliate for 'a1', got %d", len(result3))
	} else if result3[0].Name != "a" || result3[0].Bps.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("want affiliate with name 'a' and bps 1, got name '%s' and bps %s",
			result3[0].Name, result3[0].Bps.String())
	}
}
