package lvmreport

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDecode(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    any
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: io.EOF,
		},
		{
			name:  "valid JSON",
			input: `{"hello": "world"}`,
			want:  map[string]any{"hello": "world"},
		},
		{
			name:  "bad escape",
			input: `{"name": "zero\0null"}`,
			want:  map[string]any{"name": `zero\0null`},
		},
		{
			name:  "multiple bad escapes",
			input: `{"v": "zero=\u0000, zero.bad=\x00, bad=\z, newline=\\\n, end"}`,
			want: map[string]any{
				"v": "zero=\x00, zero.bad=\\x00, bad=\\z, newline=\\\n, end",
			},
		},
		{
			name:  "single null",
			input: `"\0"`,
			want:  `\0`,
		},
		{
			name:  "null only",
			input: `"\0\0\0\0"`,
			want:  `\0\0\0\0`,
		},
		{
			name:  "mixed",
			input: `"\01234\t\1\u2603\2\uAAA_\3\x\y\z"`,
			want:  "\\01234\t\\1\u2603\\2\\uAAA_\\3\\x\\y\\z",
		},
		{
			name:    "single backslash in string",
			input:   `"\"`,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "backslash only",
			input:   `\\\\\\`,
			wantErr: cmpopts.AnyError,
		},
		{
			name:    "extra content",
			input:   `{}{}`,
			want:    map[string]any{},
			wantErr: cmpopts.AnyError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var value any

			err := decode([]byte(tc.input), &value)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("decode() error diff (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, value, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("decode() result diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReader(t *testing.T) {
	for _, tc := range []struct {
		name       string
		input      string
		want       *ReportData
		wantErr    error
		wantErrMsg *regexp.Regexp
	}{
		{
			name:    "empty input",
			wantErr: io.EOF,
		},
		{
			name:       "missing report",
			input:      `{"report": []}`,
			wantErrMsg: regexp.MustCompile(`(?mi)^missing report$`),
		},
		{
			name:       "extra content",
			input:      `{}{}{}`,
			wantErrMsg: regexp.MustCompile(`(?mi)^extra data after`),
		},
		{
			name:  "multiple reports",
			input: `{"report": [{"pv": [{"a": "1"}]}, {"pv": [{"a": "2"}]}]}`,
			want: &ReportData{
				PV: []Row{{"a": "1"}, {"a": "2"}},
			},
		},
		{
			name:  "single report",
			input: `{"report": [{"vg": [], "lv": []}]}`,
			want: &ReportData{
				LV: []Row{},
				VG: []Row{},
			},
		},
		{
			name:  "escaped nulls in report",
			input: `{"report": [{"vg": [{"a": "1\0\0\0"}], "lv": [{"a": "\02\0"}]}]}`,
			want: &ReportData{
				LV: []Row{{"a": `\02\0`}},
				VG: []Row{{"a": `1\0\0\0`}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, decodeCalls := range []int{0, 1, 7} {
				t.Run(fmt.Sprint(decodeCalls), func(t *testing.T) {
					r := newReader(strings.NewReader(tc.input))

					for i := 0; i < decodeCalls; i++ {
						r.Decode()
					}

					got, err := r.Data()

					if tc.wantErr != nil {
						if err == nil || !errors.Is(err, tc.wantErr) {
							t.Errorf("reader failed with %v, want %v", err, tc.wantErr)
						}
					} else if tc.wantErrMsg != nil {
						if err == nil || !tc.wantErrMsg.MatchString(err.Error()) {
							t.Errorf("reader failed with %q, want match for %q", err, tc.wantErrMsg.String())
						}
					} else if err != nil {
						t.Errorf("reader failed with %v", err)
					}

					if err == nil {
						if diff := cmp.Diff(got, tc.want, cmpopts.EquateEmpty()); diff != "" {
							t.Errorf("Report difference (-got +want):\n%s", diff)
						}
					}
				})
			}
		})
	}
}
