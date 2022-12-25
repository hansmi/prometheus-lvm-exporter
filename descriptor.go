package main

type descriptor struct {
	_ noCopy

	fieldName string

	desc string

	metricName  string
	metricValue metricValueFunc
}
