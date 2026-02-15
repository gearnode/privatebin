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

func TestShowPasteRequestMeta_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    showPasteRequestMeta
		wantErr bool
	}{
		{
			name:  "Valid object with fields",
			input: `{"created":1707900000,"time_to_live":86400}`,
			want: showPasteRequestMeta{
				Created:    1707900000,
				TimeToLive: 86400,
			},
		},
		{
			name:  "Empty object",
			input: `{}`,
			want:  showPasteRequestMeta{},
		},
		{
			name:  "Empty array (PHP empty associative array)",
			input: `[]`,
			want:  showPasteRequestMeta{},
		},
		{
			name:  "Empty array with whitespace",
			input: `  []  `,
			want:  showPasteRequestMeta{},
		},
		{
			name:    "Non-empty array",
			input:   `[1, 2, 3]`,
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var meta showPasteRequestMeta
			err := json.Unmarshal([]byte(tt.input), &meta)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, meta)
			}
		})
	}
}
