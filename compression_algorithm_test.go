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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressionAlgorithm_String(t *testing.T) {
	tests := []struct {
		name string
		ca   CompressionAlgorithm
		want string
	}{
		{
			name: "None compression",
			ca:   CompressionAlgorithmNone,
			want: "none",
		},
		{
			name: "GZip compression",
			ca:   CompressionAlgorithmGZip,
			want: "zlib",
		},
		{
			name: "Unknown compression",
			ca:   CompressionAlgorithmUnknow,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ca.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCompressionAlgorithm_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		ca   CompressionAlgorithm
		want string
	}{
		{
			name: "Marshal none",
			ca:   CompressionAlgorithmNone,
			want: `"none"`,
		},
		{
			name: "Marshal zlib",
			ca:   CompressionAlgorithmGZip,
			want: `"zlib"`,
		},
		{
			name: "Marshal unknown",
			ca:   CompressionAlgorithmUnknow,
			want: `"unknown"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ca.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestCompressionAlgorithm_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    CompressionAlgorithm
		wantErr bool
	}{
		{
			name:  "Unmarshal none",
			input: `"none"`,
			want:  CompressionAlgorithmNone,
		},
		{
			name:  "Unmarshal zlib",
			input: `"zlib"`,
			want:  CompressionAlgorithmGZip,
		},
		{
			name:  "Unmarshal unknown string",
			input: `"gzip"`,
			want:  CompressionAlgorithmUnknow,
		},
		{
			name:  "Unmarshal empty string",
			input: `""`,
			want:  CompressionAlgorithmUnknow,
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "Non-string JSON",
			input:   `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ca CompressionAlgorithm
			err := ca.UnmarshalJSON([]byte(tt.input))
			
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, ca)
			}
		})
	}
}

func FuzzCompressionAlgorithm_UnmarshalJSON(f *testing.F) {
	// Add seed corpus with known valid and edge cases
	f.Add(`"none"`)
	f.Add(`"zlib"`)
	f.Add(`"gzip"`)
	f.Add(`"unknown"`)
	f.Add(`""`)
	f.Add(`"NONE"`)
	f.Add(`"ZLIB"`)
	f.Add(`"deflate"`)
	f.Add(`"bzip2"`)
	f.Add(`"lz4"`)
	f.Add(`"snappy"`)
	f.Add(`null`)
	f.Add(`true`)
	f.Add(`999`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`"!@#$%^&*()"`)
	f.Add(`"very long compression algorithm string that should default to unknown"`)
	
	f.Fuzz(func(t *testing.T, input string) {
		var ca CompressionAlgorithm
		err := ca.UnmarshalJSON([]byte(input))
		
		// If unmarshaling succeeds, verify the result is valid
		if err == nil {
			// The result should be None, GZip, or Unknown
			validValues := []CompressionAlgorithm{
				CompressionAlgorithmNone,
				CompressionAlgorithmGZip,
				CompressionAlgorithmUnknow,
			}
			assert.Contains(t, validValues, ca,
				"UnmarshalJSON produced invalid algorithm: %v", ca)
			
			// Verify that marshaling the result works
			data, marshalErr := ca.MarshalJSON()
			require.NoError(t, marshalErr, "Failed to marshal valid algorithm %v", ca)
			
			// Verify round-trip consistency for valid inputs
			var ca2 CompressionAlgorithm
			err2 := json.Unmarshal(data, &ca2)
			require.NoError(t, err2, "Round-trip unmarshal failed")
			assert.Equal(t, ca, ca2, "Round-trip produced different value")
		}
		// If unmarshaling fails, that's fine - we just want to ensure no panic
	})
}