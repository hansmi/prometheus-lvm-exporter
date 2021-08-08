package lvmreport

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGroupNameFields(t *testing.T) {
	got := PV.fields()
	want := []string{"pv_uuid", "pv_name", "pv_all"}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("fields difference (-got +want):\n%s", diff)
	}
}
