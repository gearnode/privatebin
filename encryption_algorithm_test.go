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

func TestEncryptionAlgorithm_String(t *testing.T) {
	tests := []struct {
		name string
		ea   EncryptionAlgorithm
		want string
	}{
		{
			name: "AES algorithm",
			ea:   EncryptionAlgorithmAES,
			want: "aes",
		},
		{
			name: "Unknown algorithm",
			ea:   EncryptionAlgorithmUnknow,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ea.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncryptionAlgorithm_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		ea   EncryptionAlgorithm
		want string
	}{
		{
			name: "Marshal AES",
			ea:   EncryptionAlgorithmAES,
			want: `"aes"`,
		},
		{
			name: "Marshal unknown",
			ea:   EncryptionAlgorithmUnknow,
			want: `"unknown"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ea.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestEncryptionAlgorithm_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    EncryptionAlgorithm
		wantErr bool
	}{
		{
			name:  "Unmarshal aes",
			input: `"aes"`,
			want:  EncryptionAlgorithmAES,
		},
		{
			name:  "Unmarshal unknown string",
			input: `"rsa"`,
			want:  EncryptionAlgorithmUnknow,
		},
		{
			name:  "Unmarshal empty string",
			input: `""`,
			want:  EncryptionAlgorithmUnknow,
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
			var ea EncryptionAlgorithm
			err := ea.UnmarshalJSON([]byte(tt.input))
			
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, ea)
			}
		})
	}
}

func FuzzEncryptionAlgorithm_UnmarshalJSON(f *testing.F) {
	// Add seed corpus with known valid and edge cases
	f.Add(`"aes"`)
	f.Add(`"unknown"`)
	f.Add(`""`)
	f.Add(`"AES"`)
	f.Add(`"aes256"`)
	f.Add(`"rsa"`)
	f.Add(`null`)
	f.Add(`123`)
	f.Add(`true`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`"very long string with spaces and special chars !@#$%^&*()"`)
	
	f.Fuzz(func(t *testing.T, input string) {
		var ea EncryptionAlgorithm
		err := ea.UnmarshalJSON([]byte(input))
		
		// If unmarshaling succeeds, verify the result is valid
		if err == nil {
			// The result should be either AES or Unknown
			assert.True(t, ea == EncryptionAlgorithmAES || ea == EncryptionAlgorithmUnknow,
				"UnmarshalJSON produced invalid algorithm: %v", ea)
			
			// Verify that marshaling the result works
			data, marshalErr := ea.MarshalJSON()
			require.NoError(t, marshalErr, "Failed to marshal valid algorithm %v", ea)
			
			// Verify round-trip consistency for valid inputs
			var ea2 EncryptionAlgorithm
			err2 := json.Unmarshal(data, &ea2)
			require.NoError(t, err2, "Round-trip unmarshal failed")
			assert.Equal(t, ea, ea2, "Round-trip produced different value")
		}
		// If unmarshaling fails, that's fine - we just want to ensure no panic
	})
}