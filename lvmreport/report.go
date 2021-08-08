package lvmreport

import "fmt"

type GroupName string

func (n GroupName) String() string {
	return string(n)
}

func (n GroupName) fields() []string {
	var selected []string

	switch n {
	case PV, LV, VG:
		selected = append(selected,
			fmt.Sprintf("%s_uuid", string(n)),
			fmt.Sprintf("%s_name", string(n)),
			fmt.Sprintf("%s_all", string(n)),
		)

	default:
		selected = append(selected, "-all")
	}

	return selected
}

const (
	PV    = GroupName("pv")
	LV    = GroupName("lv")
	VG    = GroupName("vg")
	SEG   = GroupName("seg")
	PVSEG = GroupName("pvseg")
)

var AllGroupNames = []GroupName{PV, LV, VG, SEG, PVSEG}

type ReportData struct {
	PV    []Row `json:"pv"`
	LV    []Row `json:"lv"`
	VG    []Row `json:"vg"`
	SEG   []Row `json:"seg"`
	PVSEG []Row `json:"pvseg"`
}

func (d *ReportData) merge(other ReportData) {
	d.PV = append(d.PV, other.PV...)
	d.LV = append(d.LV, other.LV...)
	d.VG = append(d.VG, other.VG...)
	d.SEG = append(d.SEG, other.SEG...)
	d.PVSEG = append(d.PVSEG, other.PVSEG...)
}

func (d *ReportData) GroupByName(name GroupName) []Row {
	switch name {
	case PV:
		return d.PV
	case LV:
		return d.LV
	case VG:
		return d.VG
	case SEG:
		return d.SEG
	case PVSEG:
		return d.PVSEG
	}

	return nil
}

type Row map[string]string
