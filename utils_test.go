// Copyright (c) 2020-2026 Bryan Frimin <bryan@frimin.fr>.
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
// REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
// INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
// LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
// OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

package privatebin

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBtoi(t *testing.T) {
	tests := []struct {
		name  string
		input bool
		want  int
	}{
		{
			name:  "true converts to 1",
			input: true,
			want:  1,
		},
		{
			name:  "false converts to 0",
			input: false,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := btoi(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestItob(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  bool
	}{
		{
			name:  "zero converts to false",
			input: 0,
			want:  false,
		},
		{
			name:  "positive integer converts to true",
			input: 1,
			want:  true,
		},
		{
			name:  "negative integer converts to true",
			input: -1,
			want:  true,
		},
		{
			name:  "large positive integer converts to true",
			input: 42,
			want:  true,
		},
		{
			name:  "large negative integer converts to true",
			input: -42,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := itob(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDecode64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:  "standard base64 with padding",
			input: "SGVsbG8gV29ybGQ=",
			want:  []byte("Hello World"),
		},
		{
			name:  "standard base64 without padding",
			input: "SGVsbG8gV29ybGQ",
			want:  []byte("Hello World"),
		},
		{
			name:  "empty string",
			input: "",
			want:  []byte{},
		},
		{
			name:  "single character encoded with padding",
			input: "QQ==",
			want:  []byte("A"),
		},
		{
			name:  "single character encoded without padding",
			input: "QQ",
			want:  []byte("A"),
		},
		{
			name:  "two characters encoded with padding",
			input: "QUI=",
			want:  []byte("AB"),
		},
		{
			name:  "two characters encoded without padding",
			input: "QUI",
			want:  []byte("AB"),
		},
		{
			name:  "three characters encoded with padding",
			input: "QUJD",
			want:  []byte("ABC"),
		},
		{
			name:  "binary data with padding",
			input: base64.StdEncoding.EncodeToString([]byte{0x00, 0x01, 0x02, 0x03}),
			want:  []byte{0x00, 0x01, 0x02, 0x03},
		},
		{
			name:  "binary data without padding",
			input: base64.RawStdEncoding.EncodeToString([]byte{0x00, 0x01, 0x02, 0x03}),
			want:  []byte{0x00, 0x01, 0x02, 0x03},
		},
		{
			name:    "invalid base64 characters",
			input:   "Invalid!@#$",
			wantErr: true,
		},
		{
			name:    "invalid padding",
			input:   "SGVsbG8===",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode64(tt.input)
			
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBtoiItobRoundTrip(t *testing.T) {
	tests := []bool{true, false}
	
	for _, input := range tests {
		t.Run(string(rune(btoi(input))+'0'), func(t *testing.T) {
			// Convert bool to int and back to bool
			result := itob(btoi(input))
			assert.Equal(t, input, result)
		})
	}
}

func TestItobBtoiEdgeCases(t *testing.T) {
	tests := []int{0, 1, -1, 42, -42, 1000000, -1000000}
	
	for _, input := range tests {
		t.Run(string(rune(input)), func(t *testing.T) {
			// Convert int to bool and back to int
			result := btoi(itob(input))
			
			// Only 0 should round-trip as 0, everything else becomes 1
			if input == 0 {
				assert.Equal(t, 0, result)
			} else {
				assert.Equal(t, 1, result)
			}
		})
	}
}

func FuzzDecode64(f *testing.F) {
	// Add seed corpus with valid base64 strings
	f.Add("SGVsbG8gV29ybGQ=")
	f.Add("SGVsbG8gV29ybGQ")
	f.Add("")
	f.Add("QQ==")
	f.Add("QQ")
	f.Add("QUJD")
	
	f.Fuzz(func(t *testing.T, input string) {
		result, err := decode64(input)
		
		if err == nil {
			// If decode succeeds, verify we can encode it back
			var encoded string
			if len(input)%4 == 0 {
				encoded = base64.StdEncoding.EncodeToString(result)
			} else {
				encoded = base64.RawStdEncoding.EncodeToString(result)
			}
			
			// Decode the encoded result and verify it matches
			result2, err2 := decode64(encoded)
			require.NoError(t, err2)
			assert.Equal(t, result, result2)
		}
		// If decode fails, that's acceptable for invalid input
	})
}