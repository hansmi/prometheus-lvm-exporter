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
