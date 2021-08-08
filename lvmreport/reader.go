package lvmreport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type root struct {
	Reports []ReportData `json:"report"`
}

type reader struct {
	dec  *json.Decoder
	data *root
	err  error
}

func newReader(r io.Reader) *reader {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	return &reader{
		dec: dec,
	}
}

func (r *reader) Decode() {
	if !(r.err == nil && r.data == nil) {
		return
	}

	var data root

	if err := r.dec.Decode(&data); err != nil {
		r.err = fmt.Errorf("decoding JSON failed: %w", err)
		return
	}

	var placeholder struct{}

	if err := r.dec.Decode(&placeholder); err == nil || !errors.Is(err, io.EOF) {
		r.err = fmt.Errorf("extra data after JSON fragment: %w", err)
		return
	}

	r.data = &data
}

func (r *reader) Data() (*ReportData, error) {
	if r.data == nil && r.err == nil {
		r.Decode()
	}

	if r.err != nil {
		return nil, r.err
	}

	if len(r.data.Reports) == 0 {
		return nil, errors.New("missing report")
	}

	result := &ReportData{}

	for _, i := range r.data.Reports {
		result.merge(i)
	}

	return result, nil
}
