package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"testing"

	"github.com/hansmi/prometheus-lvm-exporter/lvmreport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/common/expfmt"
	"github.com/sebdah/goldie/v2"
)

func disableLogOutput(t *testing.T) {
	t.Helper()

	previous := log.Writer()

	log.SetOutput(io.Discard)

	t.Cleanup(func() {
		log.SetOutput(previous)
	})
}

func gatherAndFormat(t *testing.T, c prometheus.Collector) []byte {
	t.Helper()

	reg := prometheus.NewPedanticRegistry()

	if err := prometheus.WrapRegistererWithPrefix(metricPrefix, reg).Register(c); err != nil {
		t.Fatalf("registering collector failed: %v", err)
	}

	if problems, err := testutil.GatherAndLint(reg); !(err == nil || len(problems) > 0) {
		t.Errorf("GatherAndLint() failed: %v\n%v", err, problems)
	}

	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gathering failed: %v", err)
	}

	var buf bytes.Buffer

	for _, mf := range families {
		if _, err := expfmt.MetricFamilyToText(&buf, mf); err != nil {
			t.Fatalf("MetricFamilyToText(%v) failed: %v", mf, err)
		}
	}

	return buf.Bytes()
}

func TestCollector(t *testing.T) {
	disableLogOutput(t)

	for _, tc := range []struct {
		name string
	}{
		{name: "vgdata-loop"},
		{name: "vgdata-cached"},
		{name: "multivg"},
		{name: "mirrored"},
		{name: "mirrored-unhealthy"},
		{name: "snapshot"},
		{name: "issue29-vdo-online"},
		{name: "issue30-lockargs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, enableLegacyInfoLabels := range []bool{false, true} {
				expectedName := tc.name

				if enableLegacyInfoLabels {
					expectedName += "-legacy"
				}

				t.Run(strconv.FormatBool(enableLegacyInfoLabels), func(t *testing.T) {
					c := newCollector(enableLegacyInfoLabels)
					c.load = func(ctx context.Context) (*lvmreport.ReportData, error) {
						return lvmreport.FromFile(fmt.Sprintf("testdata/%s.json", tc.name))
					}

					g := goldie.New(t)
					g.Assert(t, expectedName, gatherAndFormat(t, c))
				})
			}
		})
	}
}
