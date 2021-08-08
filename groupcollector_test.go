package main

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/hansmi/prometheus-lvm-exporter/lvmreport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sebdah/goldie/v2"
	"golang.org/x/sync/errgroup"
)

func TestGroupCollectorCollect(t *testing.T) {
	for _, tc := range []struct {
		name       string
		data       *lvmreport.ReportData
		collector  *groupCollector
		wantErrMsg []*regexp.Regexp
	}{
		{
			name: "group-collector-with-error",
			data: &lvmreport.ReportData{
				SEG: []lvmreport.Row{
					{"key1": "1", "col1": "123", "count": "true"},
					{"key1": "2", "text": "Hello World", "count": "22", "info1": "info:2"},
					{"key1": "3", "unknown key": "unknown", "count": "11.5"},
				},
			},
			collector: newGroupCollector(&group{
				name:           lvmreport.SEG,
				infoMetricName: "test_info",
				keyFields: []*descriptor{
					{
						fieldName:  "key1",
						metricName: "m_key1",
					},
				},
				infoFields: []*descriptor{
					{
						fieldName:  "info1",
						metricName: "m_info1",
					},
				},
				metricFields: []*descriptor{
					{
						fieldName:  "col1",
						metricName: "m_col1",
						metricValue: func(raw string) (float64, error) {
							return 0, errors.New("col1 error")
						},
					},
					{
						fieldName:  "count",
						metricName: "m_count",
					},
				},
			}),
			wantErrMsg: []*regexp.Regexp{
				regexp.MustCompile(`(?m)field col1: col1 error$`),
				regexp.MustCompile(`(?m)field count:.*: invalid syntax$`),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			disableLogOutput(t)

			c := newEmptyCollector()
			c.gc = append(c.gc, tc.collector)
			c.load = func(ctx context.Context) (*lvmreport.ReportData, error) {
				return tc.data, nil
			}

			g := goldie.New(t)
			g.Assert(t, tc.name, gatherAndFormat(t, c))

			{
				var g errgroup.Group

				ch := make(chan prometheus.Metric)

				g.Go(func() error {
					defer close(ch)

					for _, gc := range c.gc {
						if err := gc.collect(ch, tc.data); err != nil {
							return err
						}
					}

					return nil
				})

				for range ch {
				}

				if err := g.Wait(); len(tc.wantErrMsg) > 0 {
					match := (err != nil)

					if err != nil {
						for _, i := range tc.wantErrMsg {
							match = match && i.MatchString(err.Error())
						}
					}

					if !match {
						t.Errorf("collector failed with %q, want match for %q", err, tc.wantErrMsg)
					}
				} else if err != nil {
					t.Errorf("collector failed with %v", err)
				}
			}
		})
	}
}
