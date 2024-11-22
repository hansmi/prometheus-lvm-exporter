package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/hansmi/prometheus-lvm-exporter/lvmreport"
	"github.com/prometheus/client_golang/prometheus"
)

type groupField struct {
	fieldDesc  *descriptor
	metricDesc *prometheus.Desc
}

func (f *groupField) collect(ch chan<- prometheus.Metric, rawValue string, keyValues []string) error {
	fn := f.fieldDesc.metricValue

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

	keyFields    []string
	infoFields   []string
	metricFields map[string]*groupField
	knownFields  map[string]struct{}
}

func newGroupCollector(g *group) *groupCollector {
	c := &groupCollector{
		name:         g.name,
		unknownDesc:  prometheus.NewDesc("unknown_field_count", "Fields reported by LVM not recognized by exporter", []string{"group", "details"}, nil),
		metricFields: map[string]*groupField{},
		knownFields:  map[string]struct{}{},
	}

	var keyLabelNames []string

	for _, d := range g.keyFields {
		c.keyFields = append(c.keyFields, d.fieldName)
		keyLabelNames = append(keyLabelNames, d.metricName)
	}

	infoLabelNames := slices.Clone(keyLabelNames)

	for _, d := range g.infoFields {
		c.infoFields = append(c.infoFields, d.fieldName)
		infoLabelNames = append(infoLabelNames, d.metricName)
	}

	c.infoDesc = prometheus.NewDesc(g.infoMetricName, "", infoLabelNames, nil)

	for _, d := range g.metricFields {
		c.metricFields[d.fieldName] = &groupField{
			fieldDesc:  d,
			metricDesc: prometheus.NewDesc(d.metricName, d.desc, keyLabelNames, nil),
		}
	}

	for _, i := range g.fieldNames() {
		c.knownFields[i] = struct{}{}
	}

	return c
}

func (c *groupCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.unknownDesc

	for _, info := range c.metricFields {
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

		for _, name := range c.infoFields {
			infoValues = append(infoValues, row[name])
		}

		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1, infoValues...)

		for fieldName, rawValue := range row {
			if rawValue == "" {
				continue
			}

			info, ok := c.metricFields[fieldName]
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

	var details []string

	for i := range unknown {
		details = append(details, i)
	}

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
