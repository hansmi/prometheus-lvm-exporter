package main

import "github.com/hansmi/prometheus-lvm-exporter/lvmreport"

var pvGroup = &group{
	name:           lvmreport.PV,
	infoMetricName: "pv_info",

	keyFields: []*textField{
		{fieldName: "pv_uuid", metricName: "pv_uuid", desc: "Unique identifier"},
	},
	textFields: []*textField{
		{fieldName: "pv_attr", metricName: "pv_attr"},
		{fieldName: "pv_fmt", metricName: "pv_fmt", desc: "Type of metadata"},
		{
			fieldName:  "pv_name",
			metricName: "pv_name",
			flags:      asInfoLabel,
			desc:       "Name",
		},
		{fieldName: "pv_tags", metricName: "pv_tags"},
	},
	numericFields: []*numericField{
		{
			fieldName:  "dev_size",
			metricName: "pv_dev_size_bytes",
			desc:       "Size of underlying device",
		},
		{
			fieldName:  "pe_start",
			metricName: "pv_pe_start",
			desc:       "Offset to the start of data on the underlying device",
		},
		{
			fieldName:  "pv_allocatable",
			metricName: "pv_allocatable",
			desc:       "Set if this device can be used for allocation",
		},
		{
			fieldName:  "pv_ba_size",
			metricName: "pv_ba_size_bytes",
			desc:       "Size of PV Bootloader Area",
		},
		{
			fieldName:  "pv_ba_start",
			metricName: "pv_ba_start",
			desc:       "Offset to the start of PV Bootloader Area on the underlying device",
		},
		{
			fieldName:  "pv_duplicate",
			metricName: "pv_duplicate",
			desc:       "Set if PV is an unchosen duplicate",
		},
		{
			fieldName:  "pv_exported",
			metricName: "pv_exported",
			desc:       "Set if this device is exported",
		},
		{
			fieldName:  "pv_ext_vsn",
			metricName: "pv_ext_vsn",
			desc:       "PV header extension version",
		},
		{
			fieldName:  "pv_free",
			metricName: "pv_free_bytes",
			desc:       "Total amount of unallocated space",
		},
		{
			fieldName:  "pv_in_use",
			metricName: "pv_in_use",
			desc:       "Set if PV is used",
		},
		{
			fieldName:  "pv_major",
			metricName: "pv_major",
			desc:       "Device major number",
		},
		{
			fieldName:  "pv_mda_count",
			metricName: "pv_mda_count",
			desc:       "Number of metadata areas",
		},
		{
			fieldName:  "pv_mda_free",
			metricName: "pv_mda_free_bytes",
			desc:       "Free metadata area space",
		},
		{
			fieldName:  "pv_mda_size",
			metricName: "pv_mda_size_bytes",
			desc:       "Size of smallest metadata area",
		},
		{
			fieldName:  "pv_mda_used_count",
			metricName: "pv_mda_used_count",
			desc:       "Number of metadata areas in use",
		},
		{
			fieldName:  "pv_minor",
			metricName: "pv_minor",
			desc:       "Device minor number",
		},
		{
			fieldName:  "pv_missing",
			metricName: "pv_missing",
			desc:       "Set if this device is missing in system",
		},
		{
			fieldName:  "pv_pe_alloc_count",
			metricName: "pv_pe_alloc_count",
			desc:       "Total number of allocated Physical Extents",
		},
		{
			fieldName:  "pv_pe_count",
			metricName: "pv_pe_count",
			desc:       "Total number of Physical Extents",
		},
		{
			fieldName:  "pv_size",
			metricName: "pv_size_bytes",
			desc:       "Size of PV",
		},
		{
			fieldName:  "pv_used",
			metricName: "pv_used",
			desc:       "Total amount of allocated space",
		},
	},
}
