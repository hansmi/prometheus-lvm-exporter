package main

type descriptor struct {
	fieldName string

	desc string

	metricName  string
	metricValue metricValueFunc
}
