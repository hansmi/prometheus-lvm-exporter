package lvmreport

import (
	"context"
	"errors"
	"io"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPrepareCommand(t *testing.T) {
	c := NewCommand([]string{"/test/binary", "arg1"})

	cmd := prepareCommand(context.Background(), c.completeArgs())

	if diff := cmp.Diff(cmd.Path, "/test/binary"); diff != "" {
		t.Errorf("Path difference (-got +want):\n%s", diff)
	}

	if diff := cmp.Diff(cmd.Args[:3], []string{"/test/binary", "arg1", "fullreport"}); diff != "" {
		t.Errorf("Path difference (-got +want):\n%s", diff)
	}
}

func TestFromCommand(t *testing.T) {
	for _, tc := range []struct {
		name       string
		args       []string
		want       *ReportData
		wantErr    error
		wantErrMsg *regexp.Regexp
	}{
		{
			name:    "empty",
			args:    []string{"/bin/true"},
			wantErr: io.EOF,
		},
		{
			name:       "failure",
			args:       []string{"/bin/false"},
			wantErrMsg: regexp.MustCompile(`(?mi)\bexit status \d+\b`),
		},
		{
			name:       "no report",
			args:       []string{"/bin/echo", "{}"},
			wantErrMsg: regexp.MustCompile(`\bmissing report$`),
		},
		{
			name: "multiple reports",
			args: []string{"/bin/echo", `{
"report": [
	{"vg": [], "pv": []},
	{"vg": []}
]
}`},
			want: &ReportData{},
		},
		{
			name: "single report",
			args: []string{"/bin/echo", `{
"report": [
	{"vg": [{"vg_name": "vgAAA"}]}
]
}`},
			want: &ReportData{
				VG: []Row{
					{
						"vg_name": "vgAAA",
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			got, err := runCommand(ctx, tc.args)

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
}
