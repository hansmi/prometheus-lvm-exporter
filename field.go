package main

import "strconv"

type metricValueFunc func(string) (float64, error)

func fromNumeric(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

type fieldFlag uint

const (
	// Whether the field should be included as a label on an info metric.
	asInfoLabel fieldFlag = 1 << iota
)

type field interface {
	Name() string
	MetricName() string
	Help() string
}

// textField is an LVM report field whose value can not be made numeric, e.g.
// a device name or path.
type textField struct {
	fieldName string
	desc      string

	flags fieldFlag

	metricName string
}

var _ field = (*textField)(nil)

func (f *textField) Name() string {
	return f.fieldName
}

func (f *textField) MetricName() string {
	return f.metricName
}

func (f *textField) Help() string {
	return f.desc
}

// numericField is an LVM report field whose value is numeric or can be
// converted to a number.
type numericField struct {
	fieldName string
	desc      string

	metricName  string
	metricValue metricValueFunc
}

var _ field = (*numericField)(nil)

func (f *numericField) Name() string {
	return f.fieldName
}

func (f *numericField) MetricName() string {
	return f.metricName
}

func (f *numericField) Help() string {
	return f.desc
}
