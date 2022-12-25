package main

import "github.com/hansmi/prometheus-lvm-exporter/lvmreport"

type group struct {
	_ noCopy

	name lvmreport.GroupName

	infoMetricName string

	keyFields    []*descriptor
	infoFields   []*descriptor
	metricFields []*descriptor
}

func (r *group) allDescriptors() []*descriptor {
	d := append([]*descriptor(nil), r.keyFields...)
	d = append(d, r.infoFields...)
	d = append(d, r.metricFields...)
	return d
}

func (r *group) fieldNames() []string {
	var names []string

	for _, d := range r.allDescriptors() {
		names = append(names, d.fieldName)
	}

	return names
}

var allGroups = []*group{pvGroup, vgGroup, lvGroup}
