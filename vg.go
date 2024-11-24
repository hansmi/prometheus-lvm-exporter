package main

import "github.com/hansmi/prometheus-lvm-exporter/lvmreport"

var vgGroup = &group{
	name:           lvmreport.VG,
	infoMetricName: "vg_info",

	keyFields: []*textField{
		{fieldName: "vg_uuid", metricName: "vg_uuid", desc: "Unique identifier"},
	},
	textFields: []*textField{
		{fieldName: "vg_allocation_policy", metricName: "vg_allocation_policy"},
		{fieldName: "vg_attr", metricName: "vg_attr"},
		{fieldName: "vg_fmt", metricName: "vg_fmt", desc: "Type of metadata"},
		{fieldName: "vg_lock_args", metricName: "vg_lock_args"},
		{fieldName: "vg_lock_type", metricName: "vg_lock_type"},
		{fieldName: "vg_name", metricName: "vg_name", desc: "Name"},
		{fieldName: "vg_permissions", metricName: "vg_permissions"},
		{fieldName: "vg_systemid", metricName: "vg_systemid"},
		{fieldName: "vg_tags", metricName: "vg_tags"},
	},
	numericFields: []*numericField{
		{
			fieldName:  "lv_count",
			metricName: "vg_lv_count",
			desc:       "Number of LVs",
		},
		{
			fieldName:  "max_lv",
			metricName: "vg_max_lv",
			desc:       "Maximum number of LVs allowed in VG or 0 if unlimited",
		},
		{
			fieldName:  "max_pv",
			metricName: "vg_max_pv",
			desc:       "Maximum number of PVs allowed in VG or 0 if unlimited",
		},
		{
			fieldName:  "pv_count",
			metricName: "vg_pv_count",
			desc:       "Number of PVs in VG",
		},
		{
			fieldName:  "snap_count",
			metricName: "vg_snap_count",
			desc:       "Number of snapshots",
		},
		{
			fieldName:  "vg_clustered",
			metricName: "vg_clustered",
			desc:       "Set if VG is clustered",
		},
		{
			fieldName:  "vg_exported",
			metricName: "vg_exported",
			desc:       "Set if VG is exported",
		},
		{
			fieldName:  "vg_extendable",
			metricName: "vg_extendable",
			desc:       "Set if VG is extendable",
		},
		{
			fieldName:  "vg_extent_count",
			metricName: "vg_extent_count",
			desc:       "Total number of Physical Extents",
		},
		{
			fieldName:  "vg_extent_size",
			metricName: "vg_extent_size_bytes",
			desc:       "Size of Physical Extents",
		},
		{
			fieldName:  "vg_free",
			metricName: "vg_free_bytes",
			desc:       "Total amount of free space in bytes",
		},
		{
			fieldName:  "vg_free_count",
			metricName: "vg_free_count",
			desc:       "Total number of unallocated Physical Extents",
		},
		{
			fieldName:  "vg_mda_copies",
			metricName: "vg_mda_copies",
			desc:       "Target number of in use metadata areas in the VG",
			metricValue: func(raw string) (float64, error) {
				if raw == "unmanaged" {
					return 0, nil
				}

				return fromNumeric(raw)
			},
		},
		{
			fieldName:  "vg_mda_count",
			metricName: "vg_mda_count",
			desc:       "Number of metadata areas",
		},
		{
			fieldName:  "vg_mda_free",
			metricName: "vg_mda_free_bytes",
			desc:       "Free metadata area space for this VG",
		},
		{
			fieldName:  "vg_mda_size",
			metricName: "vg_mda_size_bytes",
			desc:       "Size of smallest metadata area for this VG",
		},
		{
			fieldName:  "vg_mda_used_count",
			metricName: "vg_mda_used_count",
			desc:       "Number of metadata areas in use on this VG",
		},
		{
			fieldName:  "vg_missing_pv_count",
			metricName: "vg_missing_pv_count",
			desc:       "Number of PVs in VG which are missing",
		},
		{
			fieldName:  "vg_partial",
			metricName: "vg_partial",
			desc:       "Set if VG is partial",
		},
		{
			fieldName:  "vg_seqno",
			metricName: "vg_seqno",
			desc:       "Revision number of internal metadata",
		},
		{
			fieldName:  "vg_shared",
			metricName: "vg_shared",
			desc:       "Set if VG is shared",
		},
		{
			fieldName:  "vg_size",
			metricName: "vg_size_bytes",
			desc:       "Total size of VG in bytes",
		},
	},
}
