package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/hansmi/prometheus-lvm-exporter/lvmreport"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/maps"
)

type groupField struct {
	metricValue metricValueFunc
	metricDesc  *prometheus.Desc
}

func (f *groupField) collect(ch chan<- prometheus.Metric, rawValue string, keyValues []string) error {
	fn := f.metricValue

	if fn == nil {
		fn = fromNumeric
	}

	value, err := fn(rawValue)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(f.metricDesc, prometheus.GaugeValue, value, keyValues...)

	return nil
}

type groupCollector struct {
	name lvmreport.GroupName

	infoDesc    *prometheus.Desc
	unknownDesc *prometheus.Desc

	keyFields     []string
	textFields    []string
	numericFields map[string]*groupField
	knownFields   map[string]struct{}
}

func newGroupCollector(g *group) *groupCollector {
	c := &groupCollector{
		name:          g.name,
		unknownDesc:   prometheus.NewDesc("unknown_field_count", "Fields reported by LVM not recognized by exporter", []string{"group", "details"}, nil),
		numericFields: map[string]*groupField{},
		knownFields:   map[string]struct{}{},
	}

	var keyLabelNames []string

	for _, f := range g.keyFields {
		c.keyFields = append(c.keyFields, f.fieldName)
		c.knownFields[f.fieldName] = struct{}{}
		keyLabelNames = append(keyLabelNames, f.metricName)
	}

	infoLabelNames := slices.Clone(keyLabelNames)

	for _, f := range g.textFields {
		c.textFields = append(c.textFields, f.fieldName)
		c.knownFields[f.fieldName] = struct{}{}
		infoLabelNames = append(infoLabelNames, f.metricName)
	}

	c.infoDesc = prometheus.NewDesc(g.infoMetricName, "", infoLabelNames, nil)

	for _, f := range g.numericFields {
		c.numericFields[f.fieldName] = &groupField{
			metricValue: f.metricValue,
			metricDesc:  prometheus.NewDesc(f.metricName, f.desc, keyLabelNames, nil),
		}
		c.knownFields[f.fieldName] = struct{}{}
	}

	return c
}

func (c *groupCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.unknownDesc

	for _, info := range c.numericFields {
		ch <- info.metricDesc
	}
}

func (c *groupCollector) collect(ch chan<- prometheus.Metric, data *lvmreport.ReportData) error {
	var allErrors prometheus.MultiError

	unknown := map[string]struct{}{}

	for _, row := range data.GroupByName(c.name) {
		var keyValues []string

		for _, name := range c.keyFields {
			keyValues = append(keyValues, row[name])
		}

		infoValues := slices.Clone(keyValues)

		for _, name := range c.textFields {
			infoValues = append(infoValues, row[name])
		}

		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1, infoValues...)

		for fieldName, rawValue := range row {
			if rawValue == "" {
				continue
			}

			info, ok := c.numericFields[fieldName]
			if !ok {
				if _, ok := c.knownFields[fieldName]; !ok {
					unknown[fieldName] = struct{}{}
				}
				continue
			}

			if err := info.collect(ch, rawValue, keyValues); err != nil {
				allErrors.Append(fmt.Errorf("field %s: %w", fieldName, err))
				continue
			}
		}
	}

	details := maps.Keys(unknown)

	sort.Strings(details)

	ch <- prometheus.MustNewConstMetric(c.unknownDesc, prometheus.GaugeValue, float64(len(unknown)), string(c.name), strings.Join(details, ", "))

	if len(allErrors) == 0 {
		return nil
	}

	slices.SortFunc(allErrors, func(a, b error) int {
		return strings.Compare(a.Error(), b.Error())
	})

	return allErrors
}
