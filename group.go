package main

import (
	"github.com/hansmi/prometheus-lvm-exporter/lvmreport"
)

type group struct {
	name lvmreport.GroupName

	infoMetricName string

	// Fields applied to all metrics from the group.
	keyFields []*textField

	// Non-numeric fields, e.g. a device path.
	textFields []*textField

	// Numeric fields, either directly or after conversion.
	numericFields []*numericField
}

var allGroups = []*group{pvGroup, vgGroup, lvGroup}
