package main

import "github.com/hansmi/prometheus-lvm-exporter/lvmreport"

var lvGroup = &group{
	name:           lvmreport.LV,
	infoMetricName: "lv_info",

	keyFields: []*textField{
		{
			fieldName:  "lv_uuid",
			metricName: "lv_uuid",
			desc:       "Unique identifier",
		},
	},
	textFields: []*textField{
		{
			fieldName:  "convert_lv",
			metricName: "lv_convert_lv",
			desc:       "For lvconvert, Name of temporary LV created by lvconvert",
		},
		{
			fieldName:  "convert_lv_uuid",
			metricName: "lv_convert_lv_uuid",
			desc:       "For lvconvert, UUID of temporary LV created by lvconvert",
		},
		{
			fieldName:  "data_lv",
			metricName: "lv_data_lv",
			desc:       "For cache/thin/vdo pools, the LV holding the associated data",
		},
		{
			fieldName:  "data_lv_uuid",
			metricName: "lv_data_lv_uuid",
			desc:       "For cache/thin/vdo pools, the UUID of the LV holding the associated data",
		},
		{
			fieldName:  "kernel_cache_policy",
			metricName: "lv_kernel_cache_policy",
			desc:       "Cache policy used in kernel",
		},
		{
			fieldName:  "kernel_cache_settings",
			metricName: "lv_kernel_cache_settings",
			desc:       "Cache settings/parameters as set in kernel, including default values (cached segments only)",
		},
		{
			fieldName:  "kernel_discards",
			metricName: "lv_kernel_discards",
			desc:       "For thin pools, how discards are handled in kernel",
		},
		{
			fieldName:  "kernel_metadata_format",
			metricName: "lv_kernel_metadata_format",
			desc:       "Cache metadata format used in kernel",
		},
		{
			fieldName:  "lv_active",
			metricName: "lv_active",
			desc:       "Active state of the LV",
		},
		{
			fieldName:  "lv_allocation_policy",
			metricName: "lv_allocation_policy",
			desc:       "LV allocation policy",
		},
		{
			fieldName:  "lv_ancestors",
			metricName: "lv_ancestors",
			desc:       "LV ancestors ignoring any stored history of the ancestry chain",
		},
		{
			fieldName:  "lv_attr",
			metricName: "lv_attr",
			desc:       "Various attributes",
		},
		{
			fieldName:  "lv_descendants",
			metricName: "lv_descendants",
			desc:       "LV descendants ignoring any stored history of the ancestry chain",
		},
		{
			fieldName:  "lv_dm_path",
			metricName: "lv_dm_path",
			desc:       "Internal device-mapper pathname for LV (in /dev/mapper directory)",
		},
		{
			fieldName:  "lv_full_ancestors",
			metricName: "lv_full_ancestors",
			desc:       "LV ancestors including stored history of the ancestry chain",
		},
		{
			fieldName:  "lv_full_descendants",
			metricName: "lv_full_descendants",
			desc:       "LV descendants including stored history of the ancestry chain",
		},
		{
			fieldName:  "lv_full_name",
			metricName: "lv_full_name",
			flags:      asInfoLabel,
			desc:       "Full name of LV including its VG, namely VG/LV",
		},
		{
			fieldName:  "lv_health_status",
			metricName: "lv_health_status",
			desc:       "LV health status",
		},
		{
			fieldName:  "lv_host",
			metricName: "lv_host",
			desc:       "Creation host of the LV, if known",
		},
		{
			fieldName:  "lv_kernel_read_ahead",
			metricName: "lv_kernel_read_ahead_bytes",
			desc:       "Currently-in-use read ahead setting",
		},
		{
			fieldName:  "lv_layout",
			metricName: "lv_layout",
			desc:       "LV layout",
		},
		{
			fieldName:  "lv_lockargs",
			metricName: "lv_lockargs",
			desc:       "Lock args of the LV used by lvmlockd",
		},
		{
			fieldName:  "lv_modules",
			metricName: "lv_modules",
			desc:       "Kernel device-mapper modules required for this LV",
		},
		{
			fieldName:  "lv_name",
			metricName: "lv_name",
			flags:      asInfoLabel,
			desc:       "Name; LVs created for internal use are enclosed in brackets",
		},
		{
			fieldName:  "lv_parent",
			metricName: "lv_parent",
			desc:       "For LVs that are components of another LV, the parent LV",
		},
		{
			fieldName:  "lv_path",
			metricName: "lv_path",
			desc:       "Full pathname for LV. Blank for internal LVs",
		},
		{
			fieldName:  "lv_permissions",
			metricName: "lv_permissions",
			desc:       "LV permissions",
		},
		{
			fieldName:  "lv_role",
			metricName: "lv_role",
			desc:       "LV role",
		},
		{
			fieldName:  "lv_tags",
			metricName: "lv_tags",
			desc:       "Tags, if any",
		},
		{
			fieldName:  "lv_when_full",
			metricName: "lv_when_full",
			desc:       "For thin pools, behavior when full",
		},
		{
			fieldName:  "metadata_lv",
			metricName: "lv_metadata_lv",
			desc:       "For cache/thin pools, the LV holding the associated metadata",
		},
		{
			fieldName:  "metadata_lv_uuid",
			metricName: "lv_metadata_lv_uuid",
			desc:       "For cache/thin pools, the UUID of the LV holding the associated metadata",
		},
		{
			fieldName:  "mirror_log",
			metricName: "lv_mirror_log",
			desc:       "For mirrors, the LV holding the synchronisation log",
		},
		{
			fieldName:  "mirror_log_uuid",
			metricName: "lv_mirror_log_uuid",
			desc:       "For mirrors, the UUID of the LV holding the synchronisation log",
		},
		{
			fieldName:  "move_pv",
			metricName: "lv_move_pv",
			desc:       "For pvmove, Source PV of temporary LV created by pvmove",
		},
		{
			fieldName:  "move_pv_uuid",
			metricName: "lv_move_pv_uuid",
			desc:       "For pvmove, the UUID of Source PV of temporary LV created by pvmove",
		},
		{
			fieldName:  "origin",
			metricName: "lv_origin",
			desc:       "For snapshots and thins, the origin device of this LV",
		},
		{
			fieldName:  "origin_uuid",
			metricName: "lv_origin_uuid",
			desc:       "For snapshots and thins, the UUID of origin device of this LV",
		},
		{
			fieldName:  "pool_lv",
			metricName: "lv_pool_lv",
			desc:       "For cache/thin/vdo volumes, the cache/thin/vdo pool LV for this volume",
		},
		{
			fieldName:  "pool_lv_uuid",
			metricName: "lv_pool_lv_uuid",
			desc:       "For cache/thin/vdo volumes, the UUID of the cache/thin/vdo pool LV for this volume",
		},
		{
			fieldName:  "raid_sync_action",
			metricName: "lv_raid_sync_action",
			desc:       "For RAID, the current synchronization action being performed",
		},
		{
			fieldName:  "raidintegritymode",
			metricName: "lv_raidintegritymode",
			desc:       "The integrity mode",
		},
		{
			fieldName:  "vdo_compression_state",
			metricName: "lv_vdo_compression_state",
			desc:       "For vdo pools, whether compression is running",
		},
		{
			fieldName:  "vdo_index_state",
			metricName: "lv_vdo_index_state",
			desc:       "For vdo pools, state of index for deduplication",
		},
		{
			fieldName:  "vdo_operating_mode",
			metricName: "lv_vdo_operating_mode",
			desc:       "For vdo pools, its current operating mode",
		},
	},
	numericFields: []*numericField{
		{
			fieldName:  "cache_dirty_blocks",
			metricName: "lv_cache_dirty_blocks",
			desc:       "Dirty cache blocks",
		},
		{
			fieldName:  "cache_read_hits",
			metricName: "lv_cache_read_hits",
			desc:       "Cache read hits",
		},
		{
			fieldName:  "cache_read_misses",
			metricName: "lv_cache_read_misses",
			desc:       "Cache read misses",
		},
		{
			fieldName:  "cache_total_blocks",
			metricName: "lv_cache_total_blocks",
			desc:       "Total cache blocks",
		},
		{
			fieldName:  "cache_used_blocks",
			metricName: "lv_cache_used_blocks",
			desc:       "Used cache blocks",
		},
		{
			fieldName:  "cache_write_hits",
			metricName: "lv_cache_write_hits",
			desc:       "Cache write hits",
		},
		{
			fieldName:  "cache_write_misses",
			metricName: "lv_cache_write_misses",
			desc:       "Cache write misses",
		},
		{
			fieldName:  "copy_percent",
			metricName: "lv_copy_percent",
			desc:       "For Cache, RAID, mirrors and pvmove, current percentage in-sync",
		},
		{
			fieldName:  "data_percent",
			metricName: "lv_data_percent",
			desc:       "For snapshot, cache and thin pools and volumes, the percentage full if LV is active",
		},
		{
			fieldName:  "integritymismatches",
			metricName: "lv_integritymismatches",
			desc:       "The number of integrity mismatches",
		},
		{
			fieldName:  "lv_active_exclusively",
			metricName: "lv_active_exclusively",
			desc:       "Set if the LV is active exclusively",
		},
		{
			fieldName:  "lv_active_locally",
			metricName: "lv_active_locally",
			desc:       "Set if the LV is active locally",
		},
		{
			fieldName:  "lv_active_remotely",
			metricName: "lv_active_remotely",
			desc:       "Set if the LV is active remotely",
		},
		{
			fieldName:  "lv_allocation_locked",
			metricName: "lv_allocation_locked",
			desc:       "Set if LV is locked against allocation changes",
		},
		{
			fieldName:  "lv_check_needed",
			metricName: "lv_check_needed",
			desc:       "For thin pools and cache volumes, whether metadata check is needed",
		},
		{
			fieldName:  "lv_converting",
			metricName: "lv_converting",
			desc:       "Set if LV is being converted",
		},
		{
			fieldName:  "lv_device_open",
			metricName: "lv_device_open",
			desc:       "Set if LV device is open",
		},
		{
			fieldName:  "lv_fixed_minor",
			metricName: "lv_fixed_minor",
			desc:       "Set if LV has fixed minor number assigned",
		},
		{
			fieldName:  "lv_historical",
			metricName: "lv_historical",
			desc:       "Set if the LV is historical",
		},
		{
			fieldName:  "lv_image_synced",
			metricName: "lv_image_synced",
			desc:       "Set if mirror/RAID image is synchronized",
		},
		{
			fieldName:  "lv_inactive_table",
			metricName: "lv_inactive_table",
			desc:       "Set if LV has inactive table present",
		},
		{
			fieldName:  "lv_initial_image_sync",
			metricName: "lv_initial_image_sync",
			desc:       "Set if mirror/RAID images underwent initial resynchronization",
		},
		{
			fieldName:  "lv_kernel_major",
			metricName: "lv_kernel_major",
			desc:       "Currently assigned major number or -1 if LV is not active",
		},
		{
			fieldName:  "lv_kernel_minor",
			metricName: "lv_kernel_minor",
			desc:       "Currently assigned minor number or -1 if LV is not active",
		},
		{
			fieldName:  "lv_live_table",
			metricName: "lv_live_table",
			desc:       "Set if LV has live table present",
		},
		{
			fieldName:  "lv_major",
			metricName: "lv_major",
			desc:       "Persistent major number or -1 if not persistent",
		},
		{
			fieldName:  "lv_merge_failed",
			metricName: "lv_merge_failed",
			desc:       "Set if snapshot merge failed",
		},
		{
			fieldName:  "lv_merging",
			metricName: "lv_merging",
			desc:       "Set if snapshot LV is being merged to origin",
		},
		{
			fieldName:  "lv_metadata_size",
			metricName: "lv_metadata_size_bytes",
			desc:       "For thin and cache pools, the size of the LV that holds the metadata",
		},
		{
			fieldName:  "lv_minor",
			metricName: "lv_minor",
			desc:       "Persistent minor number or -1 if not persistent",
		},
		{
			fieldName:  "lv_profile",
			metricName: "lv_profile",
			desc:       "Configuration profile attached to this LV",
		},
		{
			fieldName:  "lv_read_ahead",
			metricName: "lv_read_ahead_bytes",
			desc:       "Read ahead setting",
			metricValue: func(raw string) (float64, error) {
				if raw == "auto" {
					return -1, nil
				}

				return fromNumeric(raw)
			},
		},
		{
			fieldName:  "lv_size",
			metricName: "lv_size_bytes",
			desc:       "Size of LV",
		},
		{
			fieldName:  "lv_skip_activation",
			metricName: "lv_skip_activation",
			desc:       "Set if LV is skipped on activation",
		},
		{
			fieldName:  "lv_snapshot_invalid",
			metricName: "lv_snapshot_invalid",
			desc:       "Set if snapshot LV is invalid",
		},
		{
			fieldName:  "lv_suspended",
			metricName: "lv_suspended",
			desc:       "Set if LV is suspended",
		},
		{
			fieldName:  "lv_time",
			metricName: "lv_time",
			desc:       "Creation time of the LV, if known",
		},
		{
			fieldName:  "lv_time_removed",
			metricName: "lv_time_removed",
			desc:       "Removal time of the LV, if known",
		},
		{
			fieldName:  "metadata_percent",
			metricName: "lv_metadata_percent",
			desc:       "For cache and thin pools, the percentage of metadata full if LV is active",
		},
		{
			fieldName:  "origin_size",
			metricName: "lv_origin_size_bytes",
			desc:       "For snapshots, the size of the origin device of this LV",
		},
		{
			fieldName:  "raid_max_recovery_rate",
			metricName: "lv_raid_max_recovery_rate",
			desc:       "For RAID1, the maximum recovery I/O load in kiB/sec/disk",
		},
		{
			fieldName:  "raid_min_recovery_rate",
			metricName: "lv_raid_min_recovery_rate",
			desc:       "For RAID1, the minimum recovery I/O load in kiB/sec/disk",
		},
		{
			fieldName:  "raid_mismatch_count",
			metricName: "lv_raid_mismatch_count",
			desc:       "For RAID, number of mismatches found or repaired",
		},
		{
			fieldName:  "raid_write_behind",
			metricName: "lv_raid_write_behind",
			desc:       "For RAID1, the number of outstanding writes allowed to writemostly devices",
		},
		{
			fieldName:  "raidintegrityblocksize",
			metricName: "lv_raidintegrityblocksize",
			desc:       "The integrity block size",
		},
		{
			fieldName:  "seg_count",
			metricName: "lv_seg_count",
			desc:       "Number of segments in LV",
		},
		{
			fieldName:  "snap_percent",
			metricName: "lv_snap_percent",
			desc:       "For snapshots, the percentage full if LV is active",
		},
		{
			fieldName:  "sync_percent",
			metricName: "lv_sync_percent",
			desc:       "For Cache, RAID, mirrors and pvmove, current percentage in-sync",
		},
		{
			fieldName:  "vdo_saving_percent",
			metricName: "lv_vdo_saving_percent",
			desc:       "For vdo pools, percentage of saved space",
		},
		{
			fieldName:  "vdo_used_size",
			metricName: "lv_vdo_used_size_bytes",
			desc:       "For vdo pools, currently used space",
		},
		{
			fieldName:  "writecache_error",
			metricName: "lv_writecache_error",
			desc:       "Total writecache errors",
		},
		{
			fieldName:  "writecache_free_blocks",
			metricName: "lv_writecache_free_blocks",
			desc:       "Total writecache free blocks",
		},
		{
			fieldName:  "writecache_total_blocks",
			metricName: "lv_writecache_total_blocks",
			desc:       "Total writecache blocks",
		},
		{
			fieldName:  "writecache_writeback_blocks",
			metricName: "lv_writecache_writeback_blocks",
			desc:       "Total writecache writeback blocks",
		},
	},
}
