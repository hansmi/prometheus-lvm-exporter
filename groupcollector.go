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
	rawLabel   bool
	convert    metricValueFunc
	metricDesc *prometheus.Desc
}

func (f *groupField) collect(ch chan<- prometheus.Metric, rawValue string, keyValues []string) error {
	value, err := f.convert(rawValue)
	if err != nil {
		return err
	}

	labels := keyValues

	if f.rawLabel {
		labels = defaultStringSlicePool.get()
		defer defaultStringSlicePool.put(labels)

		labels = append(labels, keyValues...)
		labels = append(labels, rawValue)
	}

	ch <- prometheus.MustNewConstMetric(f.metricDesc, prometheus.GaugeValue, value, labels...)

	return nil
}

type groupTextField struct {
	metricDesc *prometheus.Desc
}

func (f *groupTextField) collect(ch chan<- prometheus.Metric, rawValue string, keyValues []string) error {
	labels := defaultStringSlicePool.get()
	defer defaultStringSlicePool.put(labels)

	labels = append(labels, keyValues...)
	labels = append(labels, rawValue)

	ch <- prometheus.MustNewConstMetric(f.metricDesc, prometheus.GaugeValue, 1, labels...)

	return nil
}

type groupCollector struct {
	name lvmreport.GroupName

	infoDesc    *prometheus.Desc
	unknownDesc *prometheus.Desc

	keyFields       []string
	textFields      map[string]*groupTextField
	numericFields   map[string]*groupField
	infoLabelFields []string
	knownFields     map[string]struct{}
}

func newGroupCollector(enableLegacyInfoLabels bool, g *group) *groupCollector {
	c := &groupCollector{
		name:          g.name,
		unknownDesc:   prometheus.NewDesc("unknown_field_count", "Fields reported by LVM not recognized by exporter", []string{"group", "details"}, nil),
		textFields:    map[string]*groupTextField{},
		numericFields: map[string]*groupField{},
		knownFields:   map[string]struct{}{},
	}

	var keyLabelNames []string

	for _, f := range g.keyFields {
		c.keyFields = append(c.keyFields, f.fieldName)
		c.knownFields[f.fieldName] = struct{}{}
		keyLabelNames = append(keyLabelNames, f.metricName)
	}

	for _, f := range g.textFields {
		c.knownFields[f.fieldName] = struct{}{}

		if f.flags&asInfoLabel == 0 {
			textLabelNames := slices.Clone(keyLabelNames)
			textLabelNames = append(textLabelNames, f.metricName)

			c.textFields[f.fieldName] = &groupTextField{
				metricDesc: prometheus.NewDesc(f.metricName, f.desc, textLabelNames, nil),
			}
		}
	}

	for _, f := range g.numericFields {
		labelNames := keyLabelNames

		if f.flags&asRawLabel != 0 {
			labelNames = slices.Clone(labelNames)
			labelNames = append(labelNames, f.metricName)
		}

		info := &groupField{
			rawLabel:   f.flags&asRawLabel != 0,
			convert:    f.metricValue,
			metricDesc: prometheus.NewDesc(f.metricName, f.desc, labelNames, nil),
		}
		if info.convert == nil {
			info.convert = fromNumeric
		}
		c.numericFields[f.fieldName] = info
		c.knownFields[f.fieldName] = struct{}{}
	}

	infoLabelNames := slices.Clone(keyLabelNames)

	for _, f := range g.textFields {
		if enableLegacyInfoLabels || f.flags&asInfoLabel != 0 {
			c.infoLabelFields = append(c.infoLabelFields, f.fieldName)
			infoLabelNames = append(infoLabelNames, f.metricName)
		}
	}

	c.infoDesc = prometheus.NewDesc(g.infoMetricName, "", infoLabelNames, nil)

	return c
}

func (c *groupCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.infoDesc
	ch <- c.unknownDesc

	for _, info := range c.textFields {
		ch <- info.metricDesc
	}

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

		for _, name := range c.infoLabelFields {
			infoValues = append(infoValues, row[name])
		}

		ch <- prometheus.MustNewConstMetric(c.infoDesc, prometheus.GaugeValue, 1, infoValues...)

		for fieldName, rawValue := range row {
			var collector interface {
				collect(chan<- prometheus.Metric, string, []string) error
			}

			if info, ok := c.textFields[fieldName]; ok {
				collector = info
			} else if rawValue == "" {
				continue
			} else if info, ok := c.numericFields[fieldName]; ok {
				collector = info
			} else {
				if _, ok := c.knownFields[fieldName]; !ok {
					unknown[fieldName] = struct{}{}
				}
				continue
			}

			if err := collector.collect(ch, rawValue, keyValues); err != nil {
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
