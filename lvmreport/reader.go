package lvmreport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
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

func decode(raw []byte, v any) error {
	var dec *json.Decoder
	var recovery struct {
		original error

		started  bool
		offset   int
		nullDone bool
	}

	for {
		dec = json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()

		err := dec.Decode(v)
		if err == nil {
			break
		}

		var errSyntax *json.SyntaxError

		// Workaround for incorrect JSON escaping in LVM. Backslashes in
		// strings are emitted without escaping. A "\0" in a device ID triggers
		// "invalid character '0' in string escape code". The troublesome
		// characters are prefixed with another backslash before parsing is
		// attempted once more.
		//
		// https://gitlab.com/lvmteam/lvm2/-/issues/35
		// https://github.com/hansmi/prometheus-lvm-exporter/issues/92
		if errors.As(err, &errSyntax) && strings.HasPrefix(errSyntax.Error(), "invalid character") {
			if !recovery.started {
				// Keep error from before recovery attempts.
				recovery.original = err

				recovery.started = true

				// Make modifications on a local slice.
				raw = slices.Clone(raw)

				if !recovery.nullDone {
					// "\0" is common enough that a full replace can be done at
					// once.
					if idx := bytes.Index(raw, []byte("\\0")); idx >= 0 {
						modified := bytes.ReplaceAll(raw[idx:], []byte("\\0"), []byte("\\\\0"))

						raw = slices.Replace(raw, idx, len(raw), modified...)

						recovery.nullDone = true
						continue
					}
				}
			}

			if offset := int(errSyntax.Offset); offset > recovery.offset {
				const maxEscapeLength = 10

				start := max(recovery.offset, offset-maxEscapeLength)

				// Never look further back than the previous error.
				recovery.offset = offset

				if idx := bytes.LastIndexByte(raw[start:offset], '\\'); idx >= 0 {
					// Add a backslash before the invalid character.
					raw = slices.Insert(raw, start+idx, '\\')

					continue
				}
			}
		}

		if recovery.original != nil {
			err = recovery.original
		}

		return fmt.Errorf("decoding JSON: %w", err)
	}

	var placeholder struct{}

	if err := dec.Decode(&placeholder); err == nil || !errors.Is(err, io.EOF) {
		return fmt.Errorf("extra data after JSON fragment: %w", err)
	}

	return nil
}

func (r *reader) Decode() {
	if !(r.err == nil && r.data == nil) {
		return
	}

	rawData, err := io.ReadAll(r.inner)
	if err != nil {
		r.err = fmt.Errorf("reading src data failed: %w", err)
		return
	}

	var data root

	if r.err = decode(rawData, &data); r.err == nil {
		r.data = &data
	}
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
