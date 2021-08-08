package lvmreport

import "os"

func FromFile(path string) (*ReportData, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	return newReader(fh).Data()
}
