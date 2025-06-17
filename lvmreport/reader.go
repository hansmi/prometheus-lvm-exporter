package lvmreport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type root struct {
	Reports []ReportData `json:"report"`
	Log     []any        `json:"log"`
}

type reader struct {
	inner io.Reader
	data  *root
	err   error
}

func newReader(r io.Reader) *reader {
	return &reader{
		inner: r,
	}
}

func (r *reader) Decode() {
	if !(r.err == nil && r.data == nil) {
		return
	}

	var data root

	rawData, err := io.ReadAll(r.inner)
	if err != nil {
		r.err = fmt.Errorf("reading src data failed: %w", err)
		return
	}

	dec := json.NewDecoder(bytes.NewReader(rawData))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&data); err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			// LVM has a bug which can add escaped null literals (`\0`) in strings
			// This is invalid JSON, so we attempt to work around it by just removing all escaped nulls
			fixedRawData := bytes.ReplaceAll(rawData, []byte{'\\', '0'}, []byte{})

			dec = json.NewDecoder(bytes.NewReader(fixedRawData))
			dec.DisallowUnknownFields()

			if err := dec.Decode(&data); err != nil {
				r.err = fmt.Errorf("decoding JSON failed: %w", err)
				return
			}
		} else {
			r.err = fmt.Errorf("decoding JSON failed: %w", err)
			return
		}
	}

	var placeholder struct{}

	if err := dec.Decode(&placeholder); err == nil || !errors.Is(err, io.EOF) {
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
