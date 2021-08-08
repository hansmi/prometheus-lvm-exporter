package main

type descriptor struct {
	noCopy noCopy

	fieldName string

	desc string

	metricName  string
	metricValue metricValueFunc
}
