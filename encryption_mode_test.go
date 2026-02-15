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

func TestEncryptionMode_String(t *testing.T) {
	tests := []struct {
		name string
		em   EncryptionMode
		want string
	}{
		{
			name: "GCM mode",
			em:   EncryptionModeGCM,
			want: "gcm",
		},
		{
			name: "Unknown mode",
			em:   EncryptionModeUnknow,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.em.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncryptionMode_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		em   EncryptionMode
		want string
	}{
		{
			name: "Marshal GCM",
			em:   EncryptionModeGCM,
			want: `"gcm"`,
		},
		{
			name: "Marshal unknown",
			em:   EncryptionModeUnknow,
			want: `"unknown"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.em.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestEncryptionMode_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    EncryptionMode
		wantErr bool
	}{
		{
			name:  "Unmarshal gcm",
			input: `"gcm"`,
			want:  EncryptionModeGCM,
		},
		{
			name:  "Unmarshal unknown string",
			input: `"cbc"`,
			want:  EncryptionModeUnknow,
		},
		{
			name:  "Unmarshal empty string",
			input: `""`,
			want:  EncryptionModeUnknow,
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "Non-string JSON",
			input:   `456`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var em EncryptionMode
			err := em.UnmarshalJSON([]byte(tt.input))
			
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, em)
			}
		})
	}
}

func FuzzEncryptionMode_UnmarshalJSON(f *testing.F) {
	// Add seed corpus with known valid and edge cases
	f.Add(`"gcm"`)
	f.Add(`"unknown"`)
	f.Add(`""`)
	f.Add(`"GCM"`)
	f.Add(`"aes-gcm"`)
	f.Add(`"cbc"`)
	f.Add(`"ctr"`)
	f.Add(`"ecb"`)
	f.Add(`null`)
	f.Add(`false`)
	f.Add(`789`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`"!@#$%^&*()"`)
	f.Add(`"very long encryption mode string that should default to unknown"`)
	
	f.Fuzz(func(t *testing.T, input string) {
		var em EncryptionMode
		err := em.UnmarshalJSON([]byte(input))
		
		// If unmarshaling succeeds, verify the result is valid
		if err == nil {
			// The result should be either GCM or Unknown
			assert.True(t, em == EncryptionModeGCM || em == EncryptionModeUnknow,
				"UnmarshalJSON produced invalid mode: %v", em)
			
			// Verify that marshaling the result works
			data, marshalErr := em.MarshalJSON()
			require.NoError(t, marshalErr, "Failed to marshal valid mode %v", em)
			
			// Verify round-trip consistency for valid inputs
			var em2 EncryptionMode
			err2 := json.Unmarshal(data, &em2)
			require.NoError(t, err2, "Round-trip unmarshal failed")
			assert.Equal(t, em, em2, "Round-trip produced different value")
		}
		// If unmarshaling fails, that's fine - we just want to ensure no panic
	})
}