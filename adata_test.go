// Copyright (c) 2020-2024 Bryan Frimin <bryan@frimin.fr>.
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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func specEqual(a, b Spec) bool {
	return bytes.Equal(a.IV, b.IV) &&
		bytes.Equal(a.Salt, b.Salt) &&
		a.Iterations == b.Iterations &&
		a.KeySize == b.KeySize &&
		a.TagSize == b.TagSize &&
		a.Algorithm == b.Algorithm &&
		a.Mode == b.Mode &&
		a.Compression == b.Compression
}

func adataEqual(a, b AData) bool {
	return specEqual(a.Spec, b.Spec) &&
		a.Formatter == b.Formatter &&
		a.OpenDiscussion == b.OpenDiscussion &&
		a.BurnAfterReading == b.BurnAfterReading
}

func TestSpec_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		spec    Spec
		wantErr bool
	}{
		{
			name: "Valid spec with all fields",
			spec: Spec{
				IV:          []byte("initialization"),
				Salt:        []byte("saltsalt"),
				Iterations:  100000,
				KeySize:     256,
				TagSize:     128,
				Algorithm:   EncryptionAlgorithmAES,
				Mode:        EncryptionModeGCM,
				Compression: CompressionAlgorithmNone,
			},
		},
		{
			name: "Empty spec",
			spec: Spec{},
		},
		{
			name: "Spec with unknown values",
			spec: Spec{
				IV:          []byte("test"),
				Salt:        []byte("test"),
				Iterations:  1,
				KeySize:     128,
				TagSize:     96,
				Algorithm:   EncryptionAlgorithmUnknow,
				Mode:        EncryptionModeUnknow,
				Compression: CompressionAlgorithmUnknow,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.spec.MarshalJSON()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if !tt.wantErr {
				// Verify it's a valid JSON array with 8 elements
				var array [8]json.RawMessage
				err = json.Unmarshal(data, &array)
				require.NoError(t, err, "Spec.MarshalJSON() produced invalid JSON array")

				// Verify round-trip
				var spec2 Spec
				err = json.Unmarshal(data, &spec2)
				require.NoError(t, err, "Failed to unmarshal marshaled Spec")
				assert.True(t, specEqual(tt.spec, spec2), "Round-trip failed: got %+v, want %+v", spec2, tt.spec)
			}
		})
	}
}

func TestSpec_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Spec
		wantErr bool
	}{
		{
			name: "Valid spec array",
			input: `["aW5pdGlhbGl6YXRpb24=","c2FsdHNhbHQ=",100000,256,128,"aes","gcm","none"]`,
			want: Spec{
				IV:          []byte("initialization"),
				Salt:        []byte("saltsalt"),
				Iterations:  100000,
				KeySize:     256,
				TagSize:     128,
				Algorithm:   EncryptionAlgorithmAES,
				Mode:        EncryptionModeGCM,
				Compression: CompressionAlgorithmNone,
			},
		},
		{
			name: "Empty base64 strings",
			input: `["","",0,0,0,"unknow","unknow","unknow"]`,
			want: Spec{
				IV:          []byte{},
				Salt:        []byte{},
				Iterations:  0,
				KeySize:     0,
				TagSize:     0,
				Algorithm:   EncryptionAlgorithmUnknow,
				Mode:        EncryptionModeUnknow,
				Compression: CompressionAlgorithmUnknow,
			},
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "Not an array",
			input:   `{"iv": "test"}`,
			wantErr: true,
		},
		{
			name:    "Wrong array length",
			input:   `["test","test",1]`,
			wantErr: true,
		},
		{
			name:    "Invalid base64 in IV",
			input:   `["!!!invalid","c2FsdA==",1,2,3,"aes","gcm","none"]`,
			wantErr: true,
		},
		{
			name:    "Invalid base64 in Salt",
			input:   `["aXY=","!!!invalid",1,2,3,"aes","gcm","none"]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var spec Spec
			err := spec.UnmarshalJSON([]byte(tt.input))
			
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.True(t, specEqual(spec, tt.want), "Spec.UnmarshalJSON() = %+v, want %+v", spec, tt.want)
			}
		})
	}
}

func TestAData_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		adata   AData
		wantErr bool
	}{
		{
			name: "Full AData",
			adata: AData{
				Spec: Spec{
					IV:          []byte("iv"),
					Salt:        []byte("salt"),
					Iterations:  100000,
					KeySize:     256,
					TagSize:     128,
					Algorithm:   EncryptionAlgorithmAES,
					Mode:        EncryptionModeGCM,
					Compression: CompressionAlgorithmGZip,
				},
				Formatter:        "plaintext",
				OpenDiscussion:   true,
				BurnAfterReading: false,
			},
		},
		{
			name: "Empty AData",
			adata: AData{
				Spec:             Spec{},
				Formatter:        "",
				OpenDiscussion:   false,
				BurnAfterReading: false,
			},
		},
		{
			name: "All flags true",
			adata: AData{
				Spec: Spec{
					IV:          []byte("test"),
					Salt:        []byte("test"),
					Iterations:  1,
					KeySize:     128,
					TagSize:     96,
					Algorithm:   EncryptionAlgorithmAES,
					Mode:        EncryptionModeGCM,
					Compression: CompressionAlgorithmNone,
				},
				Formatter:        "markdown",
				OpenDiscussion:   true,
				BurnAfterReading: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.adata.MarshalJSON()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if !tt.wantErr {
				// Verify it's a valid JSON array with 4 elements
				var array [4]json.RawMessage
				err = json.Unmarshal(data, &array)
				require.NoError(t, err, "AData.MarshalJSON() produced invalid JSON array")

				// Verify round-trip
				var adata2 AData
				err = json.Unmarshal(data, &adata2)
				require.NoError(t, err, "Failed to unmarshal marshaled AData")
				assert.True(t, adataEqual(tt.adata, adata2), "Round-trip failed: got %+v, want %+v", adata2, tt.adata)
			}
		})
	}
}

func TestAData_UnmarshalJSON(t *testing.T) {
	specJSON := `["aXY=","c2FsdA==",100000,256,128,"aes","gcm","zlib"]`
	
	tests := []struct {
		name    string
		input   string
		want    AData
		wantErr bool
	}{
		{
			name:  "Valid AData with false flags",
			input: `[` + specJSON + `,"plaintext",0,0]`,
			want: AData{
				Spec: Spec{
					IV:          []byte("iv"),
					Salt:        []byte("salt"),
					Iterations:  100000,
					KeySize:     256,
					TagSize:     128,
					Algorithm:   EncryptionAlgorithmAES,
					Mode:        EncryptionModeGCM,
					Compression: CompressionAlgorithmGZip,
				},
				Formatter:        "plaintext",
				OpenDiscussion:   false,
				BurnAfterReading: false,
			},
		},
		{
			name:  "Valid AData with true flags",
			input: `[` + specJSON + `,"markdown",1,1]`,
			want: AData{
				Spec: Spec{
					IV:          []byte("iv"),
					Salt:        []byte("salt"),
					Iterations:  100000,
					KeySize:     256,
					TagSize:     128,
					Algorithm:   EncryptionAlgorithmAES,
					Mode:        EncryptionModeGCM,
					Compression: CompressionAlgorithmGZip,
				},
				Formatter:        "markdown",
				OpenDiscussion:   true,
				BurnAfterReading: true,
			},
		},
		{
			name:  "Non-zero integers treated as true",
			input: `[` + specJSON + `,"syntaxhighlighting",42,-1]`,
			want: AData{
				Spec: Spec{
					IV:          []byte("iv"),
					Salt:        []byte("salt"),
					Iterations:  100000,
					KeySize:     256,
					TagSize:     128,
					Algorithm:   EncryptionAlgorithmAES,
					Mode:        EncryptionModeGCM,
					Compression: CompressionAlgorithmGZip,
				},
				Formatter:        "syntaxhighlighting",
				OpenDiscussion:   true,
				BurnAfterReading: true,
			},
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "Not an array",
			input:   `{"spec": {}}`,
			wantErr: true,
		},
		{
			name:    "Wrong array length",
			input:   `[` + specJSON + `,"plaintext"]`,
			wantErr: true,
		},
		{
			name:    "Invalid spec",
			input:   `[{"invalid":"spec"},"plaintext",0,0]`,
			wantErr: true,
		},
		{
			name:    "Invalid formatter type",
			input:   `[` + specJSON + `,123,0,0]`,
			wantErr: true,
		},
		{
			name:    "Invalid open discussion type",
			input:   `[` + specJSON + `,"plaintext","true",0]`,
			wantErr: true,
		},
		{
			name:    "Invalid burn after reading type",
			input:   `[` + specJSON + `,"plaintext",0,"false"]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var adata AData
			err := adata.UnmarshalJSON([]byte(tt.input))
			
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.True(t, adataEqual(adata, tt.want), "AData.UnmarshalJSON() = %+v, want %+v", adata, tt.want)
			}
		})
	}
}

func FuzzSpec_UnmarshalJSON(f *testing.F) {
	// Add seed corpus
	f.Add(`["aW5pdA==","c2FsdA==",100000,256,128,"aes","gcm","none"]`)
	f.Add(`["","",0,0,0,"unknow","unknow","unknow"]`)
	f.Add(`[]`)
	f.Add(`{}`)
	f.Add(`null`)
	f.Add(`[null,null,null,null,null,null,null,null]`)
	f.Add(`["!!!","-",1,2,3,"x","y","z"]`)
	f.Add(`[1,2,3,4,5,6,7,8]`)
	f.Add(`["dGVzdA==","dGVzdA==",1,128,96,"aes","gcm","zlib"]`)
	
	f.Fuzz(func(t *testing.T, input string) {
		var spec Spec
		err := spec.UnmarshalJSON([]byte(input))
		
		if err == nil {
			// If unmarshaling succeeds, verify marshaling works
			data, marshalErr := spec.MarshalJSON()
			require.NoError(t, marshalErr, "Failed to marshal valid Spec")
			
			// Verify round-trip consistency
			var spec2 Spec
			err2 := json.Unmarshal(data, &spec2)
			require.NoError(t, err2, "Round-trip failed")
			assert.True(t, specEqual(spec, spec2), "Round-trip produced different value")
			
			// Verify algorithm, mode, and compression are valid
			assert.True(t, spec.Algorithm == EncryptionAlgorithmAES || spec.Algorithm == EncryptionAlgorithmUnknow,
				"Invalid algorithm: %v", spec.Algorithm)
			assert.True(t, spec.Mode == EncryptionModeGCM || spec.Mode == EncryptionModeUnknow,
				"Invalid mode: %v", spec.Mode)
			validCompression := spec.Compression == CompressionAlgorithmNone ||
				spec.Compression == CompressionAlgorithmGZip ||
				spec.Compression == CompressionAlgorithmUnknow
			assert.True(t, validCompression, "Invalid compression: %v", spec.Compression)
		}
	})
}

func FuzzAData_UnmarshalJSON(f *testing.F) {
	// Add seed corpus
	validSpec := `["aXY=","c2FsdA==",100000,256,128,"aes","gcm","none"]`
	f.Add(`[` + validSpec + `,"plaintext",0,0]`)
	f.Add(`[` + validSpec + `,"markdown",1,1]`)
	f.Add(`[]`)
	f.Add(`{}`)
	f.Add(`null`)
	f.Add(`[null,null,null,null]`)
	f.Add(`[` + validSpec + `,null,null,null]`)
	f.Add(`[{},"",-1,999]`)
	f.Add(`[[],[],0,0]`)
	
	f.Fuzz(func(t *testing.T, input string) {
		var adata AData
		err := adata.UnmarshalJSON([]byte(input))
		
		if err == nil {
			// If unmarshaling succeeds, verify marshaling works
			data, marshalErr := adata.MarshalJSON()
			require.NoError(t, marshalErr, "Failed to marshal valid AData")
			
			// Verify round-trip consistency
			var adata2 AData
			err2 := json.Unmarshal(data, &adata2)
			require.NoError(t, err2, "Round-trip failed")
			assert.True(t, adataEqual(adata, adata2), "Round-trip produced different value")
			
			// Verify boolean fields are valid (booleans are always true or false in Go)
			// These checks are redundant but kept for clarity
			assert.True(t, adata.OpenDiscussion == true || adata.OpenDiscussion == false,
				"Invalid OpenDiscussion value")
			assert.True(t, adata.BurnAfterReading == true || adata.BurnAfterReading == false,
				"Invalid BurnAfterReading value")
		}
	})
}

func TestSpec_Base64Encoding(t *testing.T) {
	tests := []struct {
		name  string
		iv    []byte
		salt  []byte
	}{
		{
			name: "Standard padding",
			iv:   []byte("1234567890123456"), // 16 bytes -> needs padding
			salt: []byte("salt1234"),         // 8 bytes -> needs padding
		},
		{
			name: "No padding needed",
			iv:   []byte("123456789012"),     // 12 bytes -> no padding
			salt: []byte("123456789012"),     // 12 bytes -> no padding
		},
		{
			name: "Empty bytes",
			iv:   []byte{},
			salt: []byte{},
		},
		{
			name: "Single byte",
			iv:   []byte{0xFF},
			salt: []byte{0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := Spec{
				IV:          tt.iv,
				Salt:        tt.salt,
				Iterations:  1,
				KeySize:     128,
				TagSize:     96,
				Algorithm:   EncryptionAlgorithmAES,
				Mode:        EncryptionModeGCM,
				Compression: CompressionAlgorithmNone,
			}

			data, err := spec.MarshalJSON()
			require.NoError(t, err, "MarshalJSON failed")

			var spec2 Spec
			err = spec2.UnmarshalJSON(data)
			require.NoError(t, err, "UnmarshalJSON failed")

			assert.Equal(t, spec.IV, spec2.IV, "IV mismatch")
			assert.Equal(t, spec.Salt, spec2.Salt, "Salt mismatch")
		})
	}
}

func TestDecode64_RawEncoding(t *testing.T) {
	// Test that decode64 handles both standard and raw base64 encoding
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:  "Standard encoding with padding",
			input: base64.StdEncoding.EncodeToString([]byte("test")),
			want:  []byte("test"),
		},
		{
			name:  "Raw encoding without padding",
			input: base64.RawStdEncoding.EncodeToString([]byte("test")),
			want:  []byte("test"),
		},
		{
			name:    "Invalid base64",
			input:   "!!!invalid!!!",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode64(tt.input)
			
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got, "decode64() mismatch")
			}
		})
	}
}